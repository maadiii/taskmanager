package task

import (
	"time"
	"uuid"
)

type Entity struct {
	ID          string
	Title       string
	Description string
	Status      Status
	Priority    Priority
	DueDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func Create(title, description string) *Entity {
	return &Entity{
		ID:          uuid.NewV7().String(),
		Title:       title,
		Description: description,
		Status:      StatusTodo,
		Priority:    PriorityMedium,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
