// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: queries.sql

package db

import (
	"context"
	"database/sql"
	"time"
)

const addCategoriesToTask = `-- name: AddCategoriesToTask :exec
INSERT INTO task_categories (task_id, category_id)
VALUES (?, ?)
`

type AddCategoriesToTaskParams struct {
	TaskID     string
	CategoryID string
}

func (q *Queries) AddCategoriesToTask(ctx context.Context, arg AddCategoriesToTaskParams) error {
	_, err := q.db.ExecContext(ctx, addCategoriesToTask, arg.TaskID, arg.CategoryID)
	return err
}

const createAuthToken = `-- name: CreateAuthToken :one
INSERT INTO auth (id, user_id, token, user_agent, expires_at)
VALUES (?, ?, ?, ?, ?)
RETURNING id, user_id, token, user_agent, expires_at, created_at
`

type CreateAuthTokenParams struct {
	ID        string
	UserID    string
	Token     string
	UserAgent sql.NullString
	ExpiresAt time.Time
}

func (q *Queries) CreateAuthToken(ctx context.Context, arg CreateAuthTokenParams) (Auth, error) {
	row := q.db.QueryRowContext(ctx, createAuthToken,
		arg.ID,
		arg.UserID,
		arg.Token,
		arg.UserAgent,
		arg.ExpiresAt,
	)
	var i Auth
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.UserAgent,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const createTask = `-- name: CreateTask :one
INSERT INTO tasks (id, user_id, title, description, status, priority_id, due_date, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING id, title, description, status, priority_id, due_date, created_at, updated_at, user_id
`

type CreateTaskParams struct {
	ID          string
	UserID      string
	Title       string
	Description string
	Status      string
	PriorityID  sql.NullString
	DueDate     sql.NullTime
}

func (q *Queries) CreateTask(ctx context.Context, arg CreateTaskParams) (Task, error) {
	row := q.db.QueryRowContext(ctx, createTask,
		arg.ID,
		arg.UserID,
		arg.Title,
		arg.Description,
		arg.Status,
		arg.PriorityID,
		arg.DueDate,
	)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Status,
		&i.PriorityID,
		&i.DueDate,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, username, password, email, private_key)
VALUES (?, ?, ?, ?, ?)
RETURNING id, username, password, email, private_key, created_at, updated_at
`

type CreateUserParams struct {
	ID         string
	Username   string
	Password   string
	Email      string
	PrivateKey string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.Username,
		arg.Password,
		arg.Email,
		arg.PrivateKey,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.PrivateKey,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAuthToken = `-- name: DeleteAuthToken :exec
DELETE FROM auth
WHERE id = ?
    OR user_id = ?
    OR token = ?
`

type DeleteAuthTokenParams struct {
	ID     string
	UserID string
	Token  string
}

func (q *Queries) DeleteAuthToken(ctx context.Context, arg DeleteAuthTokenParams) error {
	_, err := q.db.ExecContext(ctx, deleteAuthToken, arg.ID, arg.UserID, arg.Token)
	return err
}

const deleteCategory = `-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = ?
`

func (q *Queries) DeleteCategory(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, deleteCategory, id)
	return err
}

const deleteTask = `-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = ? AND user_id = ?
`

type DeleteTaskParams struct {
	ID     string
	UserID string
}

func (q *Queries) DeleteTask(ctx context.Context, arg DeleteTaskParams) error {
	_, err := q.db.ExecContext(ctx, deleteTask, arg.ID, arg.UserID)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?
    OR username = ?
    OR email = ?
`

type DeleteUserParams struct {
	ID       string
	Username string
	Email    string
}

func (q *Queries) DeleteUser(ctx context.Context, arg DeleteUserParams) error {
	_, err := q.db.ExecContext(ctx, deleteUser, arg.ID, arg.Username, arg.Email)
	return err
}

const getAuthToken = `-- name: GetAuthToken :one
SELECT id, user_id, token, user_agent, expires_at, created_at
FROM auth
WHERE token = ?
    AND expires_at > CURRENT_TIMESTAMP
LIMIT 1
`

func (q *Queries) GetAuthToken(ctx context.Context, token string) (Auth, error) {
	row := q.db.QueryRowContext(ctx, getAuthToken, token)
	var i Auth
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.UserAgent,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}

const getCategories = `-- name: GetCategories :many
SELECT id, name, color
FROM categories
`

func (q *Queries) GetCategories(ctx context.Context) ([]Category, error) {
	rows, err := q.db.QueryContext(ctx, getCategories)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Category
	for rows.Next() {
		var i Category
		if err := rows.Scan(&i.ID, &i.Name, &i.Color); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPriorityForTask = `-- name: GetPriorityForTask :one
SELECT id, name, level, created_at, updated_at
FROM task_priorities
WHERE id = ?
`

func (q *Queries) GetPriorityForTask(ctx context.Context, id string) (TaskPriority, error) {
	row := q.db.QueryRowContext(ctx, getPriorityForTask, id)
	var i TaskPriority
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Level,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getTask = `-- name: GetTask :one
SELECT id, title, description, status, priority_id, due_date, created_at, updated_at, user_id
FROM tasks
WHERE id = ? AND user_id = ?
LIMIT 1
`

type GetTaskParams struct {
	ID     string
	UserID string
}

func (q *Queries) GetTask(ctx context.Context, arg GetTaskParams) (Task, error) {
	row := q.db.QueryRowContext(ctx, getTask, arg.ID, arg.UserID)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Status,
		&i.PriorityID,
		&i.DueDate,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
	)
	return i, err
}

const getTaskCategories = `-- name: GetTaskCategories :many
SELECT c.name
FROM task_categories tc
JOIN categories c ON tc.category_id = c.id
WHERE tc.task_id = ?
`

func (q *Queries) GetTaskCategories(ctx context.Context, taskID string) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getTaskCategories, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTaskPriorities = `-- name: GetTaskPriorities :many
SELECT id, name, level, created_at, updated_at
FROM task_priorities
`

func (q *Queries) GetTaskPriorities(ctx context.Context) ([]TaskPriority, error) {
	rows, err := q.db.QueryContext(ctx, getTaskPriorities)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []TaskPriority
	for rows.Next() {
		var i TaskPriority
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Level,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTasks = `-- name: GetTasks :many
SELECT id, title, description, status, priority_id, due_date, created_at, updated_at, user_id
FROM tasks
WHERE user_id = ?
ORDER BY created_at DESC
`

func (q *Queries) GetTasks(ctx context.Context, userID string) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, getTasks, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Status,
			&i.PriorityID,
			&i.DueDate,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTasksByDueDate = `-- name: GetTasksByDueDate :many
SELECT id, title, description, status, priority_id, due_date, created_at, updated_at, user_id
FROM tasks
WHERE user_id = ? AND due_date <= ?
ORDER BY due_date DESC
`

type GetTasksByDueDateParams struct {
	UserID  string
	DueDate sql.NullTime
}

func (q *Queries) GetTasksByDueDate(ctx context.Context, arg GetTasksByDueDateParams) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, getTasksByDueDate, arg.UserID, arg.DueDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Status,
			&i.PriorityID,
			&i.DueDate,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTasksByPriority = `-- name: GetTasksByPriority :many
SELECT id, title, description, status, priority_id, due_date, created_at, updated_at, user_id
FROM tasks
WHERE user_id = ? AND priority_id = ?
ORDER BY created_at DESC
`

type GetTasksByPriorityParams struct {
	UserID     string
	PriorityID sql.NullString
}

func (q *Queries) GetTasksByPriority(ctx context.Context, arg GetTasksByPriorityParams) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, getTasksByPriority, arg.UserID, arg.PriorityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Status,
			&i.PriorityID,
			&i.DueDate,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTasksByStatus = `-- name: GetTasksByStatus :many
SELECT id, title, description, status, priority_id, due_date, created_at, updated_at, user_id
FROM tasks
WHERE user_id = ? AND status = ?
ORDER BY created_at DESC
`

type GetTasksByStatusParams struct {
	UserID string
	Status string
}

func (q *Queries) GetTasksByStatus(ctx context.Context, arg GetTasksByStatusParams) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, getTasksByStatus, arg.UserID, arg.Status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Status,
			&i.PriorityID,
			&i.DueDate,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getTasksLimit = `-- name: GetTasksLimit :many
SELECT id, title, description, status, priority_id, due_date, created_at, updated_at, user_id
FROM tasks
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ?
`

type GetTasksLimitParams struct {
	UserID string
	Limit  int64
}

func (q *Queries) GetTasksLimit(ctx context.Context, arg GetTasksLimitParams) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, getTasksLimit, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Status,
			&i.PriorityID,
			&i.DueDate,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUser = `-- name: GetUser :one
SELECT id, username, password, email, private_key, created_at, updated_at
FROM users
WHERE id = ?
    OR username = ?
    OR email = ?
LIMIT 1
`

type GetUserParams struct {
	ID       string
	Username string
	Email    string
}

func (q *Queries) GetUser(ctx context.Context, arg GetUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, arg.ID, arg.Username, arg.Email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.PrivateKey,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, username, password, email, private_key, created_at, updated_at
FROM users
ORDER BY username
`

func (q *Queries) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Username,
			&i.Password,
			&i.Email,
			&i.PrivateKey,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeCategoriesFromTask = `-- name: RemoveCategoriesFromTask :exec
DELETE FROM task_categories
WHERE task_id = ? AND category_id = ?
`

type RemoveCategoriesFromTaskParams struct {
	TaskID     string
	CategoryID string
}

func (q *Queries) RemoveCategoriesFromTask(ctx context.Context, arg RemoveCategoriesFromTaskParams) error {
	_, err := q.db.ExecContext(ctx, removeCategoriesFromTask, arg.TaskID, arg.CategoryID)
	return err
}

const searchTasks = `-- name: SearchTasks :many
SELECT id, title, description, status, priority_id, due_date, created_at, updated_at, user_id
FROM tasks
WHERE user_id = ? AND (title LIKE ? OR description LIKE ?)
ORDER BY created_at DESC
`

type SearchTasksParams struct {
	UserID      string
	Title       string
	Description string
}

func (q *Queries) SearchTasks(ctx context.Context, arg SearchTasksParams) ([]Task, error) {
	rows, err := q.db.QueryContext(ctx, searchTasks, arg.UserID, arg.Title, arg.Description)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Task
	for rows.Next() {
		var i Task
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Description,
			&i.Status,
			&i.PriorityID,
			&i.DueDate,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateTask = `-- name: UpdateTask :one
UPDATE tasks
SET title = COALESCE(?, title),
    description = COALESCE(?, description),
    status = COALESCE(?, status),
    priority_id = COALESCE(?, priority_id),
    due_date = COALESCE(?, due_date),
    updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND user_id = ?
RETURNING id, title, description, status, priority_id, due_date, created_at, updated_at, user_id
`

type UpdateTaskParams struct {
	Title       string
	Description string
	Status      string
	PriorityID  sql.NullString
	DueDate     sql.NullTime
	ID          string
	UserID      string
}

func (q *Queries) UpdateTask(ctx context.Context, arg UpdateTaskParams) (Task, error) {
	row := q.db.QueryRowContext(ctx, updateTask,
		arg.Title,
		arg.Description,
		arg.Status,
		arg.PriorityID,
		arg.DueDate,
		arg.ID,
		arg.UserID,
	)
	var i Task
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Description,
		&i.Status,
		&i.PriorityID,
		&i.DueDate,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
	)
	return i, err
}

const updateTaskTimestamp = `-- name: UpdateTaskTimestamp :exec
UPDATE tasks
SET updated_at = CURRENT_TIMESTAMP
WHERE id = ?
`

func (q *Queries) UpdateTaskTimestamp(ctx context.Context, id string) error {
	_, err := q.db.ExecContext(ctx, updateTaskTimestamp, id)
	return err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET username = ?,
    password = ?,
    email = ?,
    updated_at = CURRENT_TIMESTAMP,
    private_key = ?
WHERE id = ?
    OR username = ?
    OR email = ?
RETURNING id, username, password, email, private_key, created_at, updated_at
`

type UpdateUserParams struct {
	Username   string
	Password   string
	Email      string
	PrivateKey string
	ID         string
	Username_2 string
	Email_2    string
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.Username,
		arg.Password,
		arg.Email,
		arg.PrivateKey,
		arg.ID,
		arg.Username_2,
		arg.Email_2,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.Password,
		&i.Email,
		&i.PrivateKey,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const validateAuthToken = `-- name: ValidateAuthToken :one
SELECT id, user_id, token, user_agent, expires_at, created_at
FROM auth
WHERE token = ?
  AND expires_at > CURRENT_TIMESTAMP
  AND user_id = ?
LIMIT 1
`

type ValidateAuthTokenParams struct {
	Token  string
	UserID string
}

func (q *Queries) ValidateAuthToken(ctx context.Context, arg ValidateAuthTokenParams) (Auth, error) {
	row := q.db.QueryRowContext(ctx, validateAuthToken, arg.Token, arg.UserID)
	var i Auth
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Token,
		&i.UserAgent,
		&i.ExpiresAt,
		&i.CreatedAt,
	)
	return i, err
}
