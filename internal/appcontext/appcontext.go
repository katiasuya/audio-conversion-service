// Package appcontext provides functions to work with the application context.
package appcontext

import "context"

type key int

const (
	userIDKey key = iota
)

// AddUserID adds user id to context.
func AddUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// GetUserID gets user id from context.
func GetUserID(ctx context.Context) (string, bool) {
	userIDctx, ok := ctx.Value(userIDKey).(string)
	return userIDctx, ok
}
