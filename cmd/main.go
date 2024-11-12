package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/storage/redis/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pageton/todo-list/config"
	db "github.com/pageton/todo-list/db/model"
	"github.com/pageton/todo-list/handlers"
	"github.com/pageton/todo-list/utils"
)

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder:             sonic.Marshal,
		JSONDecoder:             sonic.Unmarshal,
		Prefork:                 true,
		ServerHeader:            "Application",
		StrictRouting:           true,
		CaseSensitive:           true,
		BodyLimit:               4 * 1024 * 1024, // 4MB
		Concurrency:             256 * 1024,      // 256k
		ReadTimeout:             10 * time.Second,
		WriteTimeout:            15 * time.Second,  // 15 seconds
		IdleTimeout:             120 * time.Second, // 120 seconds =  2 minutes
		CompressedFileSuffix:    ".task.gz",
		EnableTrustedProxyCheck: true,
		TrustedProxies:          []string{"127.0.0.1"},
		ProxyHeader:             "X-Forwarded-For",
		DisableKeepalive:        false,
		EnablePrintRoutes:       false,
		GETOnly:                 false,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}})

	redisStore := redis.New(redis.Config{
		Host:      "localhost",
		Port:      6379,
		Username:  "",
		Password:  "",
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10,
	})

	if err := redisStore.Conn(); err != nil {
		log.Fatal("Failed to connect to redis store", err)
	}

	defer redisStore.Close()

	app.Use(cache.New(cache.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		Expiration:           5 * time.Minute, // 5 minutes
		CacheControl:         true,
		StoreResponseHeaders: true,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Path() + c.Query("limit") + c.Get("Authorization")
		},
		Storage:  nil,
		MaxBytes: 5 * 1024 * 1024, // 5MB
		ExpirationGenerator: func(c *fiber.Ctx, cfg *cache.Config) time.Duration {
			if c.Path() == "/api/tasks" || c.Path() == "/api/task" {
				return 2 * time.Minute // 2 minutes
			}
			return cfg.Expiration
		},
		CacheHeader: "X-Cache",
	}))

	app.Use(compress.New())
	app.Use(logger.New())
	app.Use(recover.New())

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Type", "application/json")
		c.Set("User-Agent", c.Get("User-Agent"))
		c.Set("Server", c.Get("Server"))
		c.Set("Cache-Control", "public, max-age=300")
		return c.Next()
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE",
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        100,             // 100 requests per user
		Expiration: 1 * time.Minute, // 1 minute
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() + c.Get("User-Agent") + c.Get("Authorization")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{
				"error": "Too many requests",
			})
		},
		Storage: redisStore,
	}))

	app.Use(func(c *fiber.Ctx) error {
		if deadline, ok := c.Context().Deadline(); ok {
			if time.Until(deadline) < 100*time.Millisecond {
				return c.Status(fiber.StatusRequestTimeout).JSON(
					&fiber.Map{
						"ok":      false,
						"message": "request timeout",
					})
			}
		}
		return c.Next()
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Frame-Options", "DENY")
		c.Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Set("Content-Security-Policy", "default-src 'self'")
		return c.Next()
	})

	cfg, err := config.LoadConfig()

	if err != nil {
		utils.LogError("Failed to load config: " + err.Error())
	}

	database, err := sql.Open("sqlite3", cfg.DatabasePath)

	if err != nil {
		utils.LogError("Failed to connect to database: " + err.Error())
	}

	defer database.Close()

	queries := db.New(database)

	app.Post("/api/task/create", func(c *fiber.Ctx) error {
		return handlers.CreateTaskHandler(c, queries)
	})

	app.Get("/api/tasks", func(c *fiber.Ctx) error {
		return handlers.GetTasksHandler(c, queries)
	})

	app.Get("/api/tasks/:limit", func(c *fiber.Ctx) error {
		return handlers.TaskLimiterHandler(c, queries)
	})

	app.Get("/api/task/:task_id", func(c *fiber.Ctx) error {
		return handlers.TaskByIdHandler(c, queries)
	})

	app.Put("/api/task/:task_id", func(c *fiber.Ctx) error {
		return handlers.UpdateTaskHandler(c, queries)
	})

	app.Delete("/api/task/:task_id", func(c *fiber.Ctx) error {
		return handlers.DeleteTaskHandler(c, queries)
	})

	app.Post("/api/auth/register", func(c *fiber.Ctx) error {
		return handlers.RegisterHandler(c, queries)
	})

	app.Post("/api/auth/login", func(c *fiber.Ctx) error {
		return handlers.LoginHandler(c, queries)
	})

	log.Fatal(app.Listen(cfg.Port))

}
