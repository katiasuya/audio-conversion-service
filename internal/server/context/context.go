package context

import "context"

type key int

const userID key = 0

func ContextWithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userID, id)
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	login, ok := ctx.Value(userID).(string)
	return login, ok
}
