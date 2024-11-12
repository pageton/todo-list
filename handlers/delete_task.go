package handlers

import (
	"github.com/gofiber/fiber/v2"
	db "github.com/pageton/todo-list/db/model"
	"github.com/pageton/todo-list/services"
)

func DeleteTaskHandler(c *fiber.Ctx, queries *db.Queries) error {
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

	taskId := c.Params("task_id")

	if taskId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "missing task id",
			})
	}

	_, err = queries.GetTask(c.Context(), db.GetTaskParams{
		ID:     taskId,
		UserID: tokenAuth.UserID,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "failed to get task",
			})
	}

	err = queries.DeleteTask(c.Context(), db.DeleteTaskParams{
		ID:     taskId,
		UserID: tokenAuth.UserID,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "failed to delete task",
			})
	}

	return c.Status(fiber.StatusOK).JSON(
		&fiber.Map{
			"ok":      true,
			"message": "task deleted successfully",
			"task_id": taskId,
		})
}
