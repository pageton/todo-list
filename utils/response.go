package utils

import (
	db "github.com/pageton/todo-list/db/model"
)

type ResponseTask struct {
	ID          string `json:"id"`
	UserId      string `json:"user_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PriorityID  string `json:"priority_id"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DueDate     string `json:"due_date"`
}

func ResponseTasks(t *db.Task) *ResponseTask {
	return &ResponseTask{
		ID:          t.ID,
		UserId:      t.UserID,
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		PriorityID:  t.PriorityID.String,
		CreatedAt:   t.CreatedAt.Format("2006-01-02 03:04:05 PM"),
		UpdatedAt:   t.UpdatedAt.Format("2006-01-02 03:04:05 PM"),
		DueDate:     t.DueDate.Time.Format("2006-01-02 03:04 PM"),
	}
}
