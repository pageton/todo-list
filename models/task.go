package models

import (
	"context"
	"errors"
	"fmt"

	db "github.com/pageton/todo-list/db/model"
	"github.com/pageton/todo-list/services"
)

func (t *TaskModel) Validate() error {
	if t == nil {
		return errors.New("task model is nil")
	}

	if t.Title == "" {
		return errors.New("title is required")
	}

	if t.Description == "" {
		return errors.New("description is required")
	}

	if t.Status == "" {
		return errors.New("status is required")
	}

	return nil
}

func (t *TaskModel) Encrypt(u *UserModel) error {
	if t == nil {
		return errors.New("task model is nil")
	}

	if u == nil {
		return errors.New("user model is nil")
	}

	if t.Description == "" {
		return errors.New("description is required")
	}

	privateKey, err := services.RetrievePrivateKey(u.ID)

	if err != nil {
		return fmt.Errorf("failed to retrieve private key: %w", err)
	}

	encrypted, err := services.EncryptTask(t.Description, privateKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt task: %w", err)
	}
	t.Description = encrypted
	return nil
}

func (t *TaskModel) Decrypt(u *UserModel) (string, error) {
	if t == nil {
		return "", errors.New("task model is nil")
	}

	if u == nil {
		return "", errors.New("user model is nil")
	}

	if t.Description == "" {
		return "", errors.New("description is required")
	}

	privateKey, err := services.RetrievePrivateKey(u.ID)

	if err != nil {
		return "", fmt.Errorf("failed to retrieve private key: %w", err)
	}

	decrypted, err := services.DecryptTask(t.Description, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt task: %w", err)
	}

	return decrypted, nil
}

func (t *TaskModel) UpdateTask(queries *db.Queries) error {
	if t == nil {
		return errors.New("task model is nil")
	}

	task, err := queries.GetTask(context.Background(), db.GetTaskParams{ID: t.ID, UserID: t.UserID})

	if err != nil {
		return fmt.Errorf("failed to get task: %w", err)
	}

	if t.Title == "" {
		t.Title = task.Title
	}

	if t.Description == "" {
		t.Description = task.Description
	}

	if t.Status == "" {
		t.Status = task.Status
	}

	if t.PriorityID.String == "" {
		t.PriorityID = task.PriorityID
	}

	if t.DueDate.Time.String() != "" {
		t.DueDate = task.DueDate
	}

	return nil
}
