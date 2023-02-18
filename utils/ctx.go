package utils

import "context"

type contextKey struct {
	Name string
}

var requestCtxKey = &contextKey{"context"}

func GetRequestContext(ctx context.Context, key string) string {
	if v := ctx.Value(requestCtxKey); v != nil {
		return v.(map[string]string)[key]
	} else {
		return ""
	}
}

func ResetCtxKey(ctx context.Context, key interface{}) context.Context {
	return context.WithValue(ctx, key, nil)
}
