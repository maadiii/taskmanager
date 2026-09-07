package appcontext

import "context"

type Context struct {
	context.Context

	UserID    string
	RequestID string
}
