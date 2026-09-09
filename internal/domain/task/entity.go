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

func (e *Entity) Update(ctx *appcontext.Context, title, description, status, priority string) error {
	if !slices.Contains(ctx.Identity.Permissions, "update") {
		return fmt.Errorf("forbidden")
	}

	if title != "" {
		e.Title = title
	}
	if description != "" {
		e.Description = description
	}
	if status != "" {
		parsedStatus, ok := parseStatus(status)
		if !ok {
			return fmt.Errorf("invalid status: %s", status)
		}
		e.Status = parsedStatus
	}
	if priority != "" {
		parsedPriority, ok := parsePriority(priority)
		if !ok {
			return fmt.Errorf("invalid priority: %s", priority)
		}
		e.Priority = parsedPriority
	}

	e.UpdatedAt = time.Now().UTC()

	return nil
}

func (e *Entity) Delete(ctx *appcontext.Context) error {
	if !slices.Contains(ctx.Identity.Permissions, "delete") {
		return fmt.Errorf("forbidden")
	}

	return nil
}

func parseStatus(value string) (Status, bool) {
	switch Status(value) {
	case StatusTodo, StatusInProgress, StatusDone:
		return Status(value), true
	default:
		return "", false
	}
}

func parsePriority(value string) (Priority, bool) {
	switch Priority(value) {
	case PriorityLow, PriorityMedium, PriorityHigh:
		return Priority(value), true
	default:
		return "", false
	}
}
