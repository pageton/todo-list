package models

import (
	"context"
	"errors"
	"fmt"

	db "github.com/pageton/todo-list/db/model"
	"github.com/pageton/todo-list/services"
	"golang.org/x/crypto/bcrypt"
)

func (u *UserModel) Validate() error {
	if u.Username == "" {
		return errors.New("username is required")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

func (u *UserModel) GenerateKey() error {
	key, err := services.GenerateKey()
	if err != nil {
		return err
	}
	if err := services.StorePrivateKey(u.ID, key); err != nil {
		return fmt.Errorf("failed to store private key in keyring: %v", err)
	}
	return nil
}

func (u *UserModel) CheckUserExists(queries *db.Queries) error {
	user, err := queries.GetUser(context.Background(), db.GetUserParams{
		ID: u.ID,
	})
	if err != nil {
		return err
	}
	if user.ID == "" {
		return errors.New("user does not exist")
	}

	u.PrivateKey = user.PrivateKey

	return nil
}

func (u *UserModel) HashPassword() error {
	hashedPasword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPasword)
	return nil
}

func (u *UserModel) HashValidation(password string) error {
	if password == "" {
		return errors.New("password is required")
	}
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(u.Password))
	if err != nil {
		return err
	}
	return nil
}
