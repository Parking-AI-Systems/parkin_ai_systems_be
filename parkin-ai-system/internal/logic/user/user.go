package user

import (
	"context"
	"fmt"
	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model"
	"parkin-ai-system/internal/service"
	"parkin-ai-system/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"

	"github.com/gogf/gf/v2/errors/gerror"
	"golang.org/x/crypto/bcrypt"
)

type sUser struct {
}

func Init() {
	service.RegisterUser(&sUser{})
}

func (s *sUser) SignUp(ctx context.Context, req *user.RegisterReq) (res *user.RegisterRes, err error) {
	if req.Email == "" || req.Password == "" || req.Username == "" || req.Phone == "" {
		return nil, gerror.New("Invalid input")
	}

	count, err := dao.Users.Ctx(ctx).Where("email", req.Email).Count()
	if err != nil {
		return nil, gerror.New("Database error")
	}
	if count > 0 {
		return nil, gerror.New("Email already exists")
	}

	hashedPwd, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, gerror.New("Failed to hash password")
	}

	userId, err := dao.Users.Ctx(ctx).Data(g.Map{
		"email":    req.Email,
		"password": hashedPwd,
		"username": req.Username,
		"phone":    req.Phone,
	}).InsertAndGetId()
	if err != nil {
		return nil, gerror.New("Failed to create user")
	}
	res = &user.RegisterRes{
		UserID:    fmt.Sprintf("%d", userId),
		Username:  req.Username,
		Email:     req.Email,
		Phone:     req.Phone,
		FullName:  req.FullName,
		Gender:    req.Gender,
		BirthDate: req.BirthDate,
	}
	return
}

func (s *sUser) Login(ctx context.Context, req *user.UserLoginReq) (res *user.UserLoginRes, err error) {
	if req.Account == "" || req.Password == "" {
		return nil, gerror.New("Invalid input")
	}

	userRecord, err := dao.Users.Ctx(ctx).Where("email", req.Account).One()
	if err != nil {
		return nil, gerror.New("Database error")
	}
	if userRecord.IsEmpty() {
		return nil, gerror.New("User not found")
	}

	password := userRecord["password_hash"].String()
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)); err != nil {
		return nil, gerror.New("Invalid password")
	}

	userId := userRecord["id"].Int64()

	refreshTokenStr := guid.S()

	_, err = dao.ApiTokens.Ctx(ctx).Data(g.Map{
		"user_id":     userId,
		"token":       refreshTokenStr,
		"description": "Login refresh token",
		"is_active":   true,
	}).Insert()
	if err != nil {
		return nil, gerror.New("Failed to save refresh token")
	}

	accessToken := &service.AccessToken{
		Iss: "parkin-ai-system",
		Sub: fmt.Sprintf("%d", userId),
		Exp: 0,
	}

	accessTokenStr, err := accessToken.Gen()
	if err != nil {
		return nil, gerror.New("Failed to generate access token")
	}

	res = &user.UserLoginRes{
		SignInOutput: model.SignInOutput{
			AccessTokenOutput: model.AccessTokenOutput{
				Uid:         uint64(userId),
				AccessToken: accessTokenStr,
				ExpTime:     accessToken.Exp,
			},
			RefreshToken: refreshTokenStr,
		},
	}
	return
}

func (s *sUser) RefreshToken(ctx context.Context, req *user.RefreshTokenReq) (res *user.RefreshTokenRes, err error) {
	if req.RefreshToken == "" {
		return nil, gerror.New("Refresh token is required")
	}

	tokenRecord, err := dao.ApiTokens.Ctx(ctx).Where("token", req.RefreshToken).Where("is_active", true).One()
	if err != nil {
		return nil, gerror.New("Database error")
	}
	if tokenRecord.IsEmpty() {
		return nil, gerror.New("Invalid or expired refresh token")
	}

	userId := tokenRecord["user_id"].Int64()

	accessToken := &service.AccessToken{
		Iss: "parkin-ai-system",
		Sub: fmt.Sprintf("%d", userId),
		Exp: 0,
	}

	accessTokenStr, err := accessToken.Gen()
	if err != nil {
		return nil, gerror.New("Failed to generate access token")
	}

	newRefreshTokenStr := guid.S()

	_, err = dao.ApiTokens.Ctx(ctx).Where("token", req.RefreshToken).Data(g.Map{
		"is_active": false,
	}).Update()
	if err != nil {
		return nil, gerror.New("Failed to update old token")
	}

	_, err = dao.ApiTokens.Ctx(ctx).Data(g.Map{
		"user_id":     userId,
		"token":       newRefreshTokenStr,
		"description": "Refreshed token",
		"is_active":   true,
	}).Insert()
	if err != nil {
		return nil, gerror.New("Failed to save new refresh token")
	}

	res = &user.RefreshTokenRes{
		AccessToken:  accessTokenStr,
		RefreshToken: newRefreshTokenStr,
	}
	return
}

func (s *sUser) Logout(ctx context.Context, req *user.UserLogoutReq) (res *user.UserLogoutRes, err error) {

	if refreshToken := g.RequestFromCtx(ctx).Header.Get("Refresh-Token"); refreshToken != "" {
		_, err = dao.ApiTokens.Ctx(ctx).Where("token", refreshToken).Data(g.Map{
			"is_active": false,
		}).Update()
		if err != nil {
			g.Log().Error(ctx, "Failed to deactivate refresh token:", err)
		}
	}

	if authHeader := g.RequestFromCtx(ctx).Header.Get("Authorization"); authHeader != "" {
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			accessToken := authHeader[7:]
			claims, err := utility.ParseJWT(accessToken)
			if err == nil {
				if sub, ok := claims["sub"].(string); ok {
					_, err = dao.ApiTokens.Ctx(ctx).Where("user_id", sub).Data(g.Map{
						"is_active": false,
					}).Update()
					if err != nil {
						g.Log().Error(ctx, "Failed to deactivate user tokens:", err)
					}
				}
			}
		}
	}

	res = &user.UserLogoutRes{
		Message: "Logout successful",
	}
	return
}

func (s *sUser) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (s *sUser) UserProfile(ctx context.Context, req *user.UserProfileReq) (res *user.UserProfileRes, err error) {
	// Extract user ID from context (middleware sets this)
	userIDStr := ""
	if v := g.RequestFromCtx(ctx).GetCtxVar("user_id"); v != nil {
		userIDStr = v.String()
	}
	if userIDStr == "" {
		return nil, gerror.New("Unauthorized: user_id missing in context")
	}

	// Query user info from DB
	userRecord, err := dao.Users.Ctx(ctx).Where("id", userIDStr).One()
	if err != nil {
		return nil, gerror.New("Database error")
	}
	if userRecord.IsEmpty() {
		return nil, gerror.New("User not found")
	}

	res = &user.UserProfileRes{
		UserID:    userRecord["id"].Int64(),
		Username:  userRecord["username"].String(),
		Email:     userRecord["email"].String(),
		Phone:     userRecord["phone"].String(),
		FullName:  userRecord["full_name"].String(),
		Gender:    userRecord["gender"].String(),
		BirthDate: userRecord["birth_date"].String(),
	}
	return
}
