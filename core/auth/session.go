package auth

import (
	"context"

	"github.com/nedpranson/the-christmas-breach/core/db"
)

func DeleteSession(ctx context.Context, token string) error {
	return db.GetQueries().DeleteSession(ctx, token)
}

func CleanupSessions(ctx context.Context) error {
	return db.GetQueries().DeleteExpiredSessions(ctx)
}
