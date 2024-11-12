package handlers

import (
	"github.com/gofiber/fiber/v2"
	db "github.com/pageton/todo-list/db/model"
	"github.com/pageton/todo-list/models"
	"github.com/pageton/todo-list/services"
	"github.com/pageton/todo-list/utils"
)

func UpdateTaskHandler(c *fiber.Ctx, queries *db.Queries) error {
	var task *models.TaskModel = new(models.TaskModel)

	if err := c.BodyParser(&task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "Invalid request body",
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
	taskId := c.Params("task_id")

	if taskId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "missing task id",
			})
	}

	tokenAuth, err := services.ValidateToken(token)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "invalid token",
			})
	}

	task.ID = taskId
	task.UserID = tokenAuth.UserID

	err = task.UpdateTask(queries)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "internal server error",
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

	err = task.Encrypt(user)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "internal server error",
			})
	}

	_, err = queries.UpdateTask(c.Context(), db.UpdateTaskParams{
		ID:          taskId,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		UserID:      tokenAuth.UserID,
		DueDate:     task.DueDate,
		PriorityID:  task.PriorityID,
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
			"message": "task updated successfully",
			"data": utils.ResponseTask{
				ID:          task.ID,
				UserId:      task.UserID,
				Title:       task.Title,
				Description: task.Description,
				PriorityID:  task.PriorityID.String,
				Status:      task.Status,
				CreatedAt:   task.CreatedAt.String(),
				UpdatedAt:   task.UpdatedAt.String(),
				DueDate:     task.DueDate.Time.String(),
			},
		})

}
