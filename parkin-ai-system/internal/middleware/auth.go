package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/service"
	"parkin-ai-system/utility"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/guid"
)

func Auth(r *ghttp.Request) {
	skipPaths := []string{
		"/user/login",
		"/user/register",
		"/user/refresh",
		"/health",
		"/swagger",
	}

	for _, path := range skipPaths {
		if strings.HasPrefix(r.URL.Path, path) {
			r.Middleware.Next()
			return
		}
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		r.Response.WriteStatusExit(http.StatusUnauthorized, g.Map{
			"code":    401,
			"message": "Authorization header is required",
		})
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		r.Response.WriteStatusExit(http.StatusUnauthorized, g.Map{
			"code":    401,
			"message": "Invalid authorization header format",
		})
		return
	}

	accessToken := authHeader[7:]
	claims, err := utility.ParseJWT(accessToken)
	if err != nil {
		if refreshErr := s.handleTokenRefresh(r); refreshErr != nil {
			r.Response.WriteStatusExit(http.StatusUnauthorized, g.Map{
				"code":    401,
				"message": "Invalid or expired token",
				"error":   refreshErr.Error(),
			})
			return
		}
		r.Middleware.Next()
		return
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		r.Response.WriteStatusExit(http.StatusUnauthorized, g.Map{
			"code":    401,
			"message": "Invalid token claims",
		})
		return
	}

	r.SetCtxVar("user_id", userID)
	r.SetCtxVar("claims", claims)

	r.Middleware.Next()
}

type authService struct{}

var s = &authService{}

func (s *authService) handleTokenRefresh(r *ghttp.Request) error {
	refreshToken := r.Header.Get("Refresh-Token")
	if refreshToken == "" {
		refreshTokenVar := r.Cookie.Get("refresh_token")
		if refreshTokenVar != nil {
			refreshToken = refreshTokenVar.String()
		}
	}
	if refreshToken == "" {
		return gerror.NewCode(gcode.CodeValidationFailed, "Refresh token not provided")
	}

	tokenRecord, err := dao.ApiTokens.Ctx(r.Context()).
		Where("token", refreshToken).
		Where("is_active", true).
		One()
	if err != nil {
		g.Log().Error(r.Context(), "Database error while checking refresh token:", err)
		return gerror.NewCode(gcode.CodeInternalError, "Database error")
	}
	if tokenRecord.IsEmpty() {
		return gerror.NewCode(gcode.CodeValidationFailed, "Invalid or expired refresh token")
	}

	userId := tokenRecord["user_id"].Int64()

	accessToken := &service.AccessToken{
		Iss: "parkin-ai-system",
		Sub: fmt.Sprintf("%d", userId),
		Exp: time.Now().Add(15 * time.Minute).Unix(),
	}

	newAccessToken, err := accessToken.Gen()
	if err != nil {
		g.Log().Error(r.Context(), "Failed to generate new access token:", err)
		return gerror.NewCode(gcode.CodeInternalError, "Failed to generate new access token")
	}

	newRefreshToken := guid.S()

	_, err = dao.ApiTokens.Ctx(r.Context()).
		Where("token", refreshToken).
		Data(g.Map{"is_active": false}).
		Update()
	if err != nil {
		g.Log().Error(r.Context(), "Failed to deactivate old refresh token:", err)
	}

	_, err = dao.ApiTokens.Ctx(r.Context()).Data(g.Map{
		"user_id":     userId,
		"token":       newRefreshToken,
		"description": "Auto-refreshed token",
		"is_active":   true,
	}).Insert()
	if err != nil {
		g.Log().Error(r.Context(), "Failed to save new refresh token:", err)
		newRefreshToken = refreshToken
	}

	r.Response.Header().Set("New-Access-Token", newAccessToken)
	r.Response.Header().Set("New-Refresh-Token", newRefreshToken)

	r.SetCtxVar("user_id", userId)
	r.SetCtxVar("token_refreshed", true)

	return nil
}

func RequireAuth(r *ghttp.Request) {
	Auth(r)
}

func OptionalAuth(r *ghttp.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		accessToken := authHeader[7:]
		claims, err := utility.ParseJWT(accessToken)
		if err == nil {
			if userID, ok := claims["sub"].(string); ok {
				r.SetCtxVar("user_id", userID)
				r.SetCtxVar("claims", claims)
			}
		}
	}
	r.Middleware.Next()
}
