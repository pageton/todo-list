package handlers

import (
	"github.com/gofiber/fiber/v2"
	db "github.com/pageton/todo-list/db/model"
	"github.com/pageton/todo-list/services"
	"github.com/pageton/todo-list/utils"
)

func TaskByIdHandler(c *fiber.Ctx, queries *db.Queries) error {

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

	task, err := queries.GetTask(c.Context(), db.GetTaskParams{
		ID:     taskId,
		UserID: tokenAuth.UserID,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "internal server error",
			})
	}

	return c.Status(fiber.StatusOK).JSON(
		&fiber.Map{
			"ok":      true,
			"message": "task fetched successfully",
			"task": &utils.ResponseTask{
				ID:          task.ID,
				Title:       task.Title,
				Description: task.Description,
				Status:      task.Status,
				PriorityID:  task.PriorityID.String,
				UserId:      task.UserID,
				DueDate:     task.DueDate.Time.String(),
				CreatedAt:   task.CreatedAt.String(),
				UpdatedAt:   task.UpdatedAt.String()},
		})
}
