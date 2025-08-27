package middleware

import (
	"parkin-ai-system/internal/consts"
	"parkin-ai-system/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func AdminOnly(r *ghttp.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		r.Response.Status = 401
		r.Response.WriteJson(g.Map{
			"code":    401,
			"message": "Unauthorized: Missing or invalid token",
			"data":    nil,
		})
		return
	}

	tokenString := authHeader[7:]
	claims, err := utility.ParseJWT(tokenString)
	if err != nil {
		r.Response.Status = 401
		r.Response.WriteJson(g.Map{
			"code":    401,
			"message": "Unauthorized: Invalid token",
			"data":    nil,
		})
		return
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		r.Response.Status = 401
		r.Response.WriteJson(g.Map{
			"code":    401,
			"message": "Unauthorized: Invalid user ID in token",
			"data":    nil,
		})
		return
	}

	userRecord, err := g.Model("users").Where("id", userID).One()
	if err != nil {
		g.Log().Error(r.Context(), "Database error:", err)
		r.Response.Status = 500
		r.Response.WriteJson(g.Map{
			"code":    500,
			"message": "Internal server error",
			"data":    nil,
		})
		return
	}

	if userRecord.IsEmpty() {
		r.Response.Status = 401
		r.Response.WriteJson(g.Map{
			"code":    401,
			"message": "Unauthorized: User not found",
			"data":    nil,
		})
		return
	}

	userRole := userRecord["role"].String()
	if userRole != consts.RoleAdmin {
		r.Response.Status = 403
		r.Response.WriteJson(g.Map{
			"code":    403,
			"message": "Forbidden: Admin access required",
			"data":    nil,
		})
		return
	}

	r.SetCtxVar("user_id", userID)
	r.SetCtxVar("user_role", userRole)
	r.Middleware.Next()
}

func UserOrAdmin(r *ghttp.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		r.Response.Status = 401
		r.Response.WriteJson(g.Map{
			"code":    401,
			"message": "Unauthorized: Missing or invalid token",
			"data":    nil,
		})
		return
	}

	tokenString := authHeader[7:]
	claims, err := utility.ParseJWT(tokenString)
	if err != nil {
		r.Response.Status = 401
		r.Response.WriteJson(g.Map{
			"code":    401,
			"message": "Unauthorized: Invalid token",
			"data":    nil,
		})
		return
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		r.Response.Status = 401
		r.Response.WriteJson(g.Map{
			"code":    401,
			"message": "Unauthorized: Invalid user ID in token",
			"data":    nil,
		})
		return
	}

	userRecord, err := g.Model("users").Where("id", userID).One()
	if err != nil {
		g.Log().Error(r.Context(), "Database error:", err)
		r.Response.Status = 500
		r.Response.WriteJson(g.Map{
			"code":    500,
			"message": "Internal server error",
			"data":    nil,
		})
		return
	}

	if userRecord.IsEmpty() {
		r.Response.Status = 401
		r.Response.WriteJson(g.Map{
			"code":    401,
			"message": "Unauthorized: User not found",
			"data":    nil,
		})
		return
	}

	userRole := userRecord["role"].String()
	r.SetCtxVar("user_id", userID)
	r.SetCtxVar("user_role", userRole)
	r.Middleware.Next()
}

func CheckResourceOwnership(r *ghttp.Request, resourceUserID string) bool {
	currentUserID := r.GetCtxVar("user_id").String()
	currentUserRole := r.GetCtxVar("user_role").String()

	if currentUserRole == consts.RoleAdmin {
		return true
	}

	return currentUserID == resourceUserID
}
