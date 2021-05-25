// Package mycontext provides functions to work with the application's context.
package mycontext

import "context"

type key int

const userIDKey key = 0

// ContextWithUserID adds user id to the context.
func ContextWithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userIDKey, id)
}

// UserIDFromContext retrieves user id from context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	userIDctx, ok := ctx.Value(userIDKey).(string)
	return userIDctx, ok
}
