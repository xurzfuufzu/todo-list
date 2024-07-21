package models

import (
	"github.com/google/uuid"
)

type Task struct {
	ID       uuid.UUID `json:"id"`
	Title    string    `json:"title"`
	ActiveAt string    `json:"activeAt"`
	Status   string    `json:"status"`
}

func NewTask(title, activeAt string) *Task {
	return &Task{
		ID:       uuid.New(),
		Title:    title,
		ActiveAt: activeAt,
		Status:   "active",
	}
}
