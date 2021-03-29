package context

import "context"

type key int

const userID key = 0

func SetWithUserID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, userID, id)
}

func GetUserID(ctx context.Context) (string, bool) {
	login, ok := ctx.Value(userID).(string)
	return login, ok
}
