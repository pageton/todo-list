package handlers

import (
	"database/sql"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	db "github.com/pageton/todo-list/db/model"
	"github.com/pageton/todo-list/models"
	"github.com/pageton/todo-list/services"
)

func LoginHandler(c *fiber.Ctx, queries *db.Queries) error {

	var user *models.UserModel = new(models.UserModel)

	if err := c.BodyParser(user); err != nil {

		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{
				"ok":    false,
				"error": "Invalid request body",
			})
	}

	if user.Username == "" || user.Password == "" {

		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{
				"ok":    false,
				"error": "Username and password are required",
			})
	}

	user.Username = strings.TrimSpace(user.Username)
	user.Password = strings.TrimSpace(user.Password)

	userDB, err := queries.GetUser(c.Context(),
		db.GetUserParams{
			Username: strings.ToLower(user.Username),
		})

	if err != nil {

		if err == sql.ErrNoRows {

			return c.Status(fiber.StatusUnauthorized).JSON(
				&fiber.Map{
					"ok":      false,
					"message": "Invalid username or password",
				})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "Seems we misplaced our records, please try again later",
			})
	}

	err = user.HashValidation(userDB.Password)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "Invalid username or password",
			})
	}

	token, err := services.CreateToken(userDB.ID, strings.ToLower(user.Username))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "Could not create token",
			})
	}

	userAgent := c.Get("User-Agent")

	if userAgent == "" {
		userAgent = "Unknown"
	}

	authID := uuid.New().String()

	_, err = queries.CreateAuthToken(c.Context(),
		db.CreateAuthTokenParams{
			ID:        authID,
			UserID:    userDB.ID,
			Token:     token,
			ExpiresAt: time.Now().Add(24 * 30 * time.Hour), // 30 days
			UserAgent: sql.NullString{String: userAgent, Valid: userAgent != ""},
		},
	)
	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "Could not create auth token",
			})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    "Bearer " + token,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteStrictMode,
	})

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"ok":      true,
		"message": "Login successful",
		"token":   token,
	})
}
