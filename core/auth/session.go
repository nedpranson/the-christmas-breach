package auth

import (
	"context"

	"github.com/nedpranson/the-christmas-breach/core/db"
)

func CleanupSessions(ctx context.Context) error {
	return db.GetQueries().DeleteExpiredSessions(ctx)
}
