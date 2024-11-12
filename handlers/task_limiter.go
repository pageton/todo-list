package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	db "github.com/pageton/todo-list/db/model"
	"github.com/pageton/todo-list/services"
	"github.com/pageton/todo-list/utils"
)

func TaskLimiterHandler(c *fiber.Ctx, queries *db.Queries) error {
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

	limit := c.Params("limit")

	if limit == "" {

		return c.Status(400).JSON(fiber.Map{
			"ok":      false,
			"message": "Limit is required",
		})
	}

	limitInt, err := strconv.Atoi(limit)

	if err != nil {

		return c.Status(400).JSON(fiber.Map{
			"ok":      false,
			"message": "Limit must be a number",
		})
	}

	if limitInt < 1 || limitInt > 100 {

		return c.Status(400).JSON(fiber.Map{
			"ok":      false,
			"message": "Limit must be between 1 and 100",
		})
	}

	taskLimit, err := queries.GetTasksLimit(c.Context(), db.GetTasksLimitParams{
		UserID: tokenAuth.UserID,
		Limit:  int64(limitInt),
	})

	if err != nil {

		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "internal server error",
			})
	}

	if len(taskLimit) == 0 {

		return c.Status(fiber.StatusNoContent).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "no tasks found",
			})
	}

	var responseTasks []utils.ResponseTask

	for _, task := range taskLimit {
		responseTasks = append(responseTasks, *utils.ResponseTasks(&task))
	}

	return c.Status(fiber.StatusOK).JSON(
		&fiber.Map{
			"ok":      true,
			"limit":   limitInt,
			"count":   len(taskLimit),
			"message": "tasks have been fetched successfully",
			"data": fiber.Map{
				"tasks": responseTasks,
			},
		})

}
