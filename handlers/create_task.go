package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	db "github.com/pageton/todo-list/db/model"
	"github.com/pageton/todo-list/models"
	"github.com/pageton/todo-list/services"
)

func CreateTaskHandler(c *fiber.Ctx, queries *db.Queries) error {
	var task *models.TaskModel = new(models.TaskModel)

	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "invalid request body",
			})
	}

	if err := task.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "invalid task",
			})
	}

	authHeader := c.Get("Authorization")

	if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		return c.Status(fiber.StatusUnauthorized).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "missing or invalid authorization token",
			})
	}

	token := authHeader[7:]

	tokenAuth, err := services.ValidateToken(token)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "invalid token",
			})
	}

	user := &models.UserModel{User: &db.User{ID: tokenAuth.UserID}}

	if err := user.CheckUserExists(queries); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "user not found",
			})
	}

	taskID := uuid.New().String()

	task.ID = taskID

	err = task.Encrypt(user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "internal server error",
			})
	}

	params := &db.CreateTaskParams{
		ID:          task.ID,
		UserID:      tokenAuth.UserID,
		Title:       task.Title,
		Status:      task.Status,
		DueDate:     task.DueDate,
		Description: task.Description,
	}

	_, err = queries.CreateTask(c.Context(), *params)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "internal server error",
			})
	}

	return c.Status(fiber.StatusCreated).JSON(
		&fiber.Map{
			"ok":      true,
			"message": "task created successfully",
			"task_id": task.ID,
		})
}
