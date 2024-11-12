package handlers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	db "github.com/pageton/todo-list/db/model"
	"github.com/pageton/todo-list/models"
)

func RegisterHandler(c *fiber.Ctx, queries *db.Queries) error {
	var user *models.UserModel = new(models.UserModel)

	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "Invalid request body",
			})
	}

	if err := user.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "Invalid form data or validation error",
			})
	}

	existingUser, err := queries.GetUser(c.Context(),
		db.GetUserParams{
			Username: strings.ToLower(user.Username),
		})

	if err == nil && strings.ToLower(existingUser.Username) != "" {
		return c.Status(fiber.StatusBadRequest).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "Username already exists",
			})
	}

	if err := user.HashPassword(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "authentication failed",
			})
	}

	user.ID = uuid.New().String()

	err = user.GenerateKey()

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "Something went wrong",
			})
	}

	_, err = queries.CreateUser(c.Context(),
		db.CreateUserParams{
			ID:         user.ID,
			Username:   user.Username,
			Password:   user.Password,
			Email:      user.Email,
			PrivateKey: user.PrivateKey,
		})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(
			&fiber.Map{
				"ok":      false,
				"message": "could not register user",
			})
	}

	return c.Status(fiber.StatusCreated).JSON(
		&fiber.Map{
			"ok":       true,
			"message":  "user registered successfully",
			"id":       user.ID,
			"username": user.Username,
		})
}
