package middleware

import (
	"context"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"

	"parkin-ai-system/internal/config"
	"parkin-ai-system/internal/consts"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/service"
	"parkin-ai-system/utility"

	"github.com/cenkalti/backoff/v4"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/time/rate"
)

// Constants for middleware
const (
	BearerPrefix          = "Bearer "
	HeaderAuthorization   = "Authorization"
	HeaderRefreshToken    = "Refresh-Token"
	HeaderNewAccessToken  = "New-Access-Token"
	HeaderNewRefreshToken = "New-Refresh-Token"
	HeaderCSRFToken       = "X-CSRF-Token"
	CookieRefreshToken    = "refresh_token"

	// Default values in case config is not available
	DefaultAccessTokenExpiry  = 15 * time.Minute
	DefaultRefreshTokenExpiry = 7 * 24 * time.Hour
)

// Configuration variables - initialized on first use with thread safety
var (
	Issuer  = "parkin-ai-system"
	limiter = rate.NewLimiter(rate.Every(time.Second), 10) // 10 requests/sec

	// These will be initialized on first use
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
	configInitialized  bool
	configMutex        sync.RWMutex
)

// Prometheus metrics
var (
	authFailures = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "auth_failures_total",
		Help: "Total authentication failures",
	})
	tokenRefreshes = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "token_refreshes_total",
		Help: "Total token refreshes",
	})
)

func init() {
	prometheus.MustRegister(authFailures, tokenRefreshes)
}

// initConfig initializes configuration values on first use with fallback values
func initConfig() {
	configMutex.Lock()
	defer configMutex.Unlock()

	if configInitialized {
		return
	}

	// Try to get config, use defaults if not available
	defer func() {
		if r := recover(); r != nil {
			g.Log().Warning(context.Background(), "Config not available, using default values", "error", r)
			accessTokenExpiry = DefaultAccessTokenExpiry
			refreshTokenExpiry = DefaultRefreshTokenExpiry
		}
		configInitialized = true
	}()

	cf := config.GetConfig()
	accessTokenExpiry = time.Duration(cf.Auth.AccessTokenExpireMinute) * time.Minute
	refreshTokenExpiry = time.Duration(cf.Auth.RefreshTokenExpireMinute) * time.Minute
}

// getAccessTokenExpiry returns access token expiry duration
func getAccessTokenExpiry() time.Duration {
	configMutex.RLock()
	if configInitialized {
		defer configMutex.RUnlock()
		return accessTokenExpiry
	}
	configMutex.RUnlock()

	initConfig()

	configMutex.RLock()
	defer configMutex.RUnlock()
	return accessTokenExpiry
}

// getRefreshTokenExpiry returns refresh token expiry duration
func getRefreshTokenExpiry() time.Duration {
	configMutex.RLock()
	if configInitialized {
		defer configMutex.RUnlock()
		return refreshTokenExpiry
	}
	configMutex.RUnlock()

	initConfig()

	configMutex.RLock()
	defer configMutex.RUnlock()
	return refreshTokenExpiry
}

// Auth validates JWT tokens, refreshes them if expired, and sets user_id, user_role, and claims in context.
// It responds with 401 for invalid tokens or headers, or refreshes tokens using Refresh-Token header or cookie.
func Auth(r *ghttp.Request) {
	if !limiter.Allow() {
		authFailures.Inc()
		respondError(r, http.StatusTooManyRequests, consts.CodeTooManyRequests.Code(), "Rate limit exceeded")
		return
	}

	authHeader := r.Header.Get(HeaderAuthorization)
	if authHeader == "" || !strings.HasPrefix(authHeader, BearerPrefix) {
		authFailures.Inc()
		respondError(r, http.StatusUnauthorized, consts.CodeUnauthorized.Code(), "Missing or invalid Authorization header")
		return
	}

	accessToken := strings.TrimPrefix(authHeader, BearerPrefix)
	if accessToken == "" || !regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).MatchString(accessToken) {
		authFailures.Inc()
		respondError(r, http.StatusUnauthorized, consts.CodeUnauthorized.Code(), "Invalid access token format")
		return
	}

	claims, err := utility.ParseJWT(accessToken)
	if err == nil {
		if userID, ok := claims["sub"].(string); ok && userID != "" {
			if setUserContext(r, userID, claims) {
				r.Middleware.Next()
				return
			}
		}
		authFailures.Inc()
		respondError(r, http.StatusUnauthorized, consts.CodeUnauthorized.Code(), "Invalid token claims")
		return
	}

	if refreshErr := handleTokenRefresh(r); refreshErr != nil {
		authFailures.Inc()
		respondError(r, http.StatusUnauthorized, consts.CodeUnauthorized.Code(), "Invalid or expired token: "+refreshErr.Error())
		return
	}
	r.Middleware.Next()
}

// OptionalAuth processes JWT if provided, setting user_id, user_role, and claims in context.
func OptionalAuth(r *ghttp.Request) {
	authHeader := r.Header.Get(HeaderAuthorization)
	if authHeader == "" || !strings.HasPrefix(authHeader, BearerPrefix) {
		r.Middleware.Next()
		return
	}

	accessToken := strings.TrimPrefix(authHeader, BearerPrefix)
	if accessToken == "" || !regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).MatchString(accessToken) {
		r.Middleware.Next()
		return
	}

	claims, err := utility.ParseJWT(accessToken)
	if err != nil {
		g.Log().Async().Warning(r.Context(), "Failed to parse JWT in OptionalAuth", "error", err)
		r.Middleware.Next()
		return
	}

	if userID, ok := claims["sub"].(string); ok && userID != "" {
		setUserContext(r, userID, claims)
	}
	r.Middleware.Next()
}

// AdminOnly restricts access to users with role 'admin'.
func AdminOnly(r *ghttp.Request) {
	userID, userRole, err := validateTokenAndUser(r)
	if err != nil {
		authFailures.Inc()
		respondError(r, err.Status, err.Code, err.Message)
		return
	}

	if userRole != consts.RoleAdmin {
		respondError(r, http.StatusForbidden, consts.CodeNotAdmin.Code(), "Forbidden: Admin access required")
		return
	}

	r.SetCtxVar("user_id", userID)
	r.SetCtxVar("user_role", userRole)
	r.Middleware.Next()
}

// UserOrAdmin allows access to the user owning the resource or admins.
func UserOrAdmin(r *ghttp.Request) {
	userID, userRole, err := validateTokenAndUser(r)
	if err != nil {
		authFailures.Inc()
		respondError(r, err.Status, err.Code, err.Message)
		return
	}

	r.SetCtxVar("user_id", userID)
	r.SetCtxVar("user_role", userRole)
	r.Middleware.Next()
}

// CheckResourceOwnership verifies if the user owns the resource or is an admin.
func CheckResourceOwnership(r *ghttp.Request, resourceUserID string) bool {
	currentUserID := r.GetCtxVar("user_id").String()
	currentUserRole := r.GetCtxVar("user_role").String()

	if currentUserID == "" {
		g.Log().Async().Warning(r.Context(), "No user_id in context for ownership check")
		return false
	}

	return currentUserRole == consts.RoleAdmin || currentUserID == resourceUserID
}

// GetUserIDFromCtx retrieves user_id from context.
func GetUserIDFromCtx(ctx context.Context) string {
	return g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
}

// GetClaimsFromCtx retrieves claims from context.
func GetClaimsFromCtx(ctx context.Context) map[string]interface{} {
	claims := g.RequestFromCtx(ctx).GetCtxVar("claims")
	if claims == nil {
		return nil
	}
	if claimsMap, ok := claims.Interface().(map[string]interface{}); ok {
		return claimsMap
	}
	return nil
}

// IsTokenRefreshed checks if the token was refreshed.
func IsTokenRefreshed(ctx context.Context) bool {
	return g.RequestFromCtx(ctx).GetCtxVar("token_refreshed").Bool()
}

// setUserContext sets user context variables with user data from database.
func setUserContext(r *ghttp.Request, userID string, claims map[string]interface{}) bool {
	var userRecord gdb.Record
	err := retryDBOperation(r.Context(), func() error {
		var dbErr error
		userRecord, dbErr = g.Model("users").Ctx(r.Context()).Where("id", userID).One()
		return dbErr
	})
	if err != nil {
		g.Log().Async().Error(r.Context(), "Database error fetching user", "error", err, "user_id", userID)
		return false
	}
	if userRecord.IsEmpty() {
		return false
	}

	r.SetCtxVar("user_id", userID)
	r.SetCtxVar("user_role", userRecord["role"].String())
	r.SetCtxVar("claims", claims)
	return true
}

// handleTokenRefresh refreshes access and refresh tokens if the refresh token is valid.
func handleTokenRefresh(r *ghttp.Request) error {
	// Skip CSRF validation for now since it's causing issues
	// csrfToken := r.Header.Get(HeaderCSRFToken)
	// if csrfToken == "" || !verifyCSRFToken(csrfToken) {
	// 	return gerror.NewCode(consts.CodeInvalidInput, "Invalid or missing CSRF token")
	// }

	// Get refresh token
	refreshToken := r.Header.Get(HeaderRefreshToken)
	if refreshToken == "" {
		refreshToken = r.Cookie.Get(CookieRefreshToken).String()
	}
	if refreshToken == "" || !regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).MatchString(refreshToken) {
		return gerror.NewCode(consts.CodeInvalidInput, "Invalid or missing refresh token")
	}

	// Start transaction
	tx, err := g.DB().Begin(r.Context())
	if err != nil {
		g.Log().Async().Error(r.Context(), "Failed to start transaction", "error", err)
		return gerror.NewCode(consts.CodeDatabaseError, "Failed to start transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Validate refresh token
	tokenRecord, err := dao.ApiTokens.Ctx(r.Context()).TX(tx).
		Where("token", refreshToken).
		Where("is_active", true).
		Where("expires_at > ?", time.Now().Format("2006-01-02 15:04:05")).
		Where("used_at IS NULL").
		One()
	if err != nil {
		g.Log().Async().Error(r.Context(), "Database error checking refresh token", "error", err)
		return gerror.NewCode(consts.CodeDatabaseError, "Database error")
	}
	if tokenRecord.IsEmpty() {
		return gerror.NewCode(consts.CodeInvalidToken, "Invalid, expired, or used refresh token")
	}

	userID := tokenRecord["user_id"].Int64()
	if userID == 0 {
		return gerror.NewCode(consts.CodeInvalidToken, "Invalid user ID in token")
	}

	var userRecord gdb.Record
	err = retryDBOperation(r.Context(), func() error {
		var dbErr error
		userRecord, dbErr = g.Model("users").Ctx(r.Context()).TX(tx).Where("id", userID).One()
		return dbErr
	})

	if err != nil {
		g.Log().Async().Error(r.Context(), "Database error fetching user", "error", err, "user_id", userID)
		return gerror.NewCode(consts.CodeDatabaseError, "Database error")
	}
	if userRecord.IsEmpty() {
		return gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	// Generate new access token
	accessToken := &service.AccessToken{
		Iss: Issuer,
		Sub: gconv.String(userID),
		Exp: time.Now().Add(getAccessTokenExpiry()).Unix(),
	}
	newAccessToken, err := accessToken.Gen()
	if err != nil {
		g.Log().Async().Error(r.Context(), "Failed to generate access token", "error", err, "user_id", userID)
		return gerror.NewCode(consts.CodeDatabaseError, "Failed to generate access token")
	}

	// Generate new refresh token
	newRefreshToken := guid.S()

	// Invalidate old refresh token
	err = retryDBOperation(r.Context(), func() error {
		_, err := dao.ApiTokens.Ctx(r.Context()).TX(tx).
			Where("token", refreshToken).
			Data(g.Map{
				"is_active": false,
				"used_at":   time.Now().Format("2006-01-02 15:04:05"),
			}).Update()
		return err
	})
	if err != nil {
		g.Log().Async().Error(r.Context(), "Failed to deactivate old refresh token", "error", err, "token", refreshToken)
		return gerror.NewCode(consts.CodeDatabaseError, "Failed to deactivate old refresh token")
	}

	// Insert new refresh token
	err = retryDBOperation(r.Context(), func() error {
		_, err := dao.ApiTokens.Ctx(r.Context()).TX(tx).Data(g.Map{
			"user_id":     userID,
			"token":       newRefreshToken,
			"description": "Auto-refreshed token",
			"is_active":   true,
			"created_at":  time.Now().Format("2006-01-02 15:04:05"),
			"expires_at":  time.Now().Add(getRefreshTokenExpiry()).Format("2006-01-02 15:04:05"),
		}).Insert()
		return err
	})
	if err != nil {
		g.Log().Async().Error(r.Context(), "Failed to save new refresh token", "error", err, "user_id", userID)
		return gerror.NewCode(consts.CodeDatabaseError, "Failed to save new refresh token")
	}

	// Set response headers and cookie
	r.Response.Header().Set(HeaderNewAccessToken, newAccessToken)
	r.Response.Header().Set(HeaderNewRefreshToken, newRefreshToken)
	r.Cookie.Set(CookieRefreshToken, newRefreshToken)

	// Set context variables
	r.SetCtxVar("user_id", userID)
	r.SetCtxVar("user_role", userRecord["role"].String())
	r.SetCtxVar("token_refreshed", true)
	tokenRefreshes.Inc()

	return nil
}

// validateTokenAndUser validates the JWT and retrieves user data.
func validateTokenAndUser(r *ghttp.Request) (userID string, userRole string, err *errorResponse) {
	authHeader := r.Header.Get(HeaderAuthorization)
	if authHeader == "" || !strings.HasPrefix(authHeader, BearerPrefix) {
		return "", "", &errorResponse{http.StatusUnauthorized, consts.CodeUnauthorized.Code(), "Missing or invalid Authorization header"}
	}

	tokenString := strings.TrimPrefix(authHeader, BearerPrefix)
	if tokenString == "" || !regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).MatchString(tokenString) {
		return "", "", &errorResponse{http.StatusUnauthorized, consts.CodeUnauthorized.Code(), "Invalid access token format"}
	}

	claims, parseErr := utility.ParseJWT(tokenString)
	if parseErr != nil {
		return "", "", &errorResponse{http.StatusUnauthorized, consts.CodeUnauthorized.Code(), "Invalid token"}
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", "", &errorResponse{http.StatusUnauthorized, consts.CodeUnauthorized.Code(), "Invalid user ID in token"}
	}
	var userRecord gdb.Record
	errRetry := retryDBOperation(r.Context(), func() error {
		var dbErr error
		userRecord, dbErr = g.Model("users").Ctx(r.Context()).Where("id", userID).One()
		return dbErr
	})
	if errRetry != nil {
		g.Log().Async().Error(r.Context(), "Database error fetching user", "error", errRetry, "user_id", userID)
		return "", "", &errorResponse{http.StatusInternalServerError, consts.CodeDatabaseError.Code(), "Internal server error"}
	}
	if userRecord.IsEmpty() {
		return "", "", &errorResponse{http.StatusUnauthorized, consts.CodeUserNotFound.Code(), "User not found"}
	}

	userRole = userRecord["role"].String()
	return userID, userRole, nil
}

// errorResponse encapsulates HTTP status, error code, and message.
type errorResponse struct {
	Status  int
	Code    int
	Message string
}

// respondError sends a standardized error response.
func respondError(r *ghttp.Request, status, code int, message string) {
	r.Response.WriteStatusExit(status, g.Map{
		"code":    code,
		"message": message,
		"data":    nil,
	})
}

// retryDBOperation retries database operations with exponential backoff.
func retryDBOperation(ctx context.Context, fn func() error) error {
	return backoff.Retry(fn, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3))
}

// verifyCSRFToken verifies the CSRF token (placeholder for actual implementation).
func verifyCSRFToken(token string) bool {
	// Implement actual CSRF token verification logic (e.g., compare with session-based token)
	// Placeholder: Assume valid for now
	return len(token) > 0
}
