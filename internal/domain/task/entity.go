package task

import (
	"fmt"
	"slices"
	"time"
	"uuid"

	"github.com/maadiii/taskmanager/pkg/appcontext"
)

type Entity struct {
	ID          string
	UserID      string
	Title       string
	Description string
	Status      Status
	Priority    Priority
	DueDate     time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func Create(ctx *appcontext.Context, title, description string) (*Entity, error) {
	// Creation business rule
	if !slices.Contains(ctx.Identity.Permissions, "create") {
		return nil, fmt.Errorf("forbidden")
	}

	return &Entity{
		ID:          uuid.NewV7().String(),
		UserID:      ctx.Identity.UserID,
		Title:       title,
		Description: description,
		Status:      StatusTodo,
		Priority:    PriorityMedium,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}
