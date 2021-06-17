// Package appcontext provides functions to work with the application context.
package appcontext

import "context"

type key int

const (
	userIDKey key = iota
	requestIDKey
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

// AddRequestID adds request id to context.
func AddRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

// GetRequestID get request id from context.
func GetRequestID(ctx context.Context) (string, bool) {
	requestIDctx, ok := ctx.Value(requestIDKey).(string)
	return requestIDctx, ok
}
