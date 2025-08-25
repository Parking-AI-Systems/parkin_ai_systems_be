package middleware

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
)

func GetUserIDFromCtx(ctx context.Context) string {
	if userID := g.RequestFromCtx(ctx).GetCtxVar("user_id"); userID != nil {
		return userID.String()
	}
	return ""
}

func GetClaimsFromCtx(ctx context.Context) map[string]interface{} {
	if claims := g.RequestFromCtx(ctx).GetCtxVar("claims"); claims != nil {
		if claimsMap, ok := claims.Interface().(map[string]interface{}); ok {
			return claimsMap
		}
	}
	return nil
}

func IsTokenRefreshed(ctx context.Context) bool {
	if refreshed := g.RequestFromCtx(ctx).GetCtxVar("token_refreshed"); refreshed != nil {
		return refreshed.Bool()
	}
	return false
}
