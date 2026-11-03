package user

import "context"

type contextKey struct{}

func WithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, contextKey{}, user)
}

func FromContext(ctx context.Context) *User {
	user, ok := ctx.Value(contextKey{}).(*User)
	if !ok {
		return nil
	}
	return user
}
