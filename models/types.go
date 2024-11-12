package models

import db "github.com/pageton/todo-list/db/model"

type UserModel struct {
	*db.User
}

type TaskModel struct {
	*db.Task
}

type Validate interface {
	Validate() error
}

type Task interface {
	UpdateTask(queries *db.Queries) error
}

type Cipher interface {
	Encrypt(u *UserModel) error
	Decrypt(u *UserModel) error
	GenerateKey() error
	HashPassword() error
	HashValidation(password string) error
}

type Check interface {
	CheckUserExists(queries *db.Queries) error
}
