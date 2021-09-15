package handler

import "context"

type ContextKeyUID struct{}

func ReadContextString(ctx context.Context, key interface{}) string {
	v := ctx.Value(key)
	if v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}
