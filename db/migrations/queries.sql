-- name: AddCategoriesToTask :exec
INSERT INTO task_categories (task_id, category_id)
VALUES (?, ?);

-- name: CreateAuthToken :one
INSERT INTO auth (id, user_id, token, user_agent, expires_at)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: CreateTask :one
INSERT INTO tasks (id, user_id, title, description, status, priority_id, due_date, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
RETURNING *;

-- name: CreateUser :one
INSERT INTO users (id, username, password, email, private_key)
VALUES (?, ?, ?, ?, ?)
RETURNING *;

-- name: DeleteAuthToken :exec
DELETE FROM auth
WHERE id = ?
    OR user_id = ?
    OR token = ?;

-- name: DeleteCategory :exec
DELETE FROM categories
WHERE id = ?;

-- name: DeleteTask :exec
DELETE FROM tasks
WHERE id = ? AND user_id = ?;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = ?
    OR username = ?
    OR email = ?;

-- name: GetAuthToken :one
SELECT *
FROM auth
WHERE token = ?
    AND expires_at > CURRENT_TIMESTAMP
LIMIT 1;

-- name: GetCategories :many
SELECT *
FROM categories;

-- name: GetPriorityForTask :one
SELECT *
FROM task_priorities
WHERE id = ?;

-- name: GetTask :one
SELECT *
FROM tasks
WHERE id = ? AND user_id = ?
LIMIT 1;

-- name: GetTaskCategories :many
SELECT c.name
FROM task_categories tc
JOIN categories c ON tc.category_id = c.id
WHERE tc.task_id = ?;

-- name: GetTaskPriorities :many
SELECT *
FROM task_priorities;

-- name: GetTasks :many
SELECT *
FROM tasks
WHERE user_id = ?
ORDER BY created_at DESC;

-- name: GetTasksLimit :many
SELECT *
FROM tasks
WHERE user_id = ?
ORDER BY created_at DESC
LIMIT ?;

-- name: GetTasksByDueDate :many
SELECT *
FROM tasks
WHERE user_id = ? AND due_date <= ?
ORDER BY due_date DESC;

-- name: GetTasksByPriority :many
SELECT *
FROM tasks
WHERE user_id = ? AND priority_id = ?
ORDER BY created_at DESC;

-- name: GetTasksByStatus :many
SELECT *
FROM tasks
WHERE user_id = ? AND status = ?
ORDER BY created_at DESC;

-- name: GetUser :one
SELECT *
FROM users
WHERE id = ?
    OR username = ?
    OR email = ?
LIMIT 1;

-- name: GetUsers :many
SELECT *
FROM users
ORDER BY username;

-- name: RemoveCategoriesFromTask :exec
DELETE FROM task_categories
WHERE task_id = ? AND category_id = ?;

-- name: SearchTasks :many
SELECT *
FROM tasks
WHERE user_id = ? AND (title LIKE ? OR description LIKE ?)
ORDER BY created_at DESC;

-- name: UpdateTask :one
UPDATE tasks
SET title = COALESCE(?, title),
    description = COALESCE(?, description),
    status = COALESCE(?, status),
    priority_id = COALESCE(?, priority_id),
    due_date = COALESCE(?, due_date),
    updated_at = CURRENT_TIMESTAMP
WHERE id = ? AND user_id = ?
RETURNING *;

-- name: UpdateTaskTimestamp :exec
UPDATE tasks
SET updated_at = CURRENT_TIMESTAMP
WHERE id = ?;

-- name: UpdateUser :one
UPDATE users
SET username = ?,
    password = ?,
    email = ?,
    updated_at = CURRENT_TIMESTAMP,
    private_key = ?
WHERE id = ?
    OR username = ?
    OR email = ?
RETURNING *;

-- name: ValidateAuthToken :one
SELECT *
FROM auth
WHERE token = ?
  AND expires_at > CURRENT_TIMESTAMP
  AND user_id = ?
LIMIT 1;
