// Package mycontext provides functions to work with the application's context.
package mycontext

import "context"

type key int

const userID key = 0

// ContextWithUserID adds user id to the context.
func ContextWithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userID, id)
}

// UserIDFromContext retrieves user id from context.
func UserIDFromContext(ctx context.Context) (string, bool) {
	login, ok := ctx.Value(userID).(string)
	return login, ok
}
