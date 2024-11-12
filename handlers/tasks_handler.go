package handlers

import (
	"github.com/gofiber/fiber/v2"
	db "github.com/pageton/todo-list/db/model"
	"github.com/pageton/todo-list/services"
	"github.com/pageton/todo-list/utils"
)

func GetTasksHandler(c *fiber.Ctx, queries *db.Queries) error {

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

	tasks, err := queries.GetTasks(c.Context(), tokenAuth.UserID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "internal server error",
			})
	}

	if len(tasks) <= 0 {
		return c.Status(fiber.StatusNoContent).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "no tasks found",
			})
	}

	var responseTasks []utils.ResponseTask

	for _, task := range tasks {
		responseTasks = append(responseTasks, *utils.ResponseTasks(&task))
	}

	return c.Status(fiber.StatusOK).JSON(
		&fiber.Map{
			"ok":      true,
			"count":   len(tasks),
			"message": "tasks have been fetched successfully",
			"data": fiber.Map{
				"tasks": responseTasks,
			},
		})
}
