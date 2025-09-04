package user

import (
	"context"
	"fmt"
	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/consts"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/middleware"
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

func init() {
	Init()
}

func (s *sUser) SignUp(ctx context.Context, req *user.RegisterReq) (res *user.RegisterRes, err error) {
	if req.Email == "" || req.Password == "" || req.Username == "" || req.Phone == "" {
		return nil, gerror.NewCode(consts.CodeInvalidInput)
	}

	count, err := dao.Users.Ctx(ctx).Where("email", req.Email).Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError)
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeEmailExists)
	}

	hashedPwd, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeHashPasswordFailed)
	}

	userId, err := dao.Users.Ctx(ctx).Data(g.Map{
		"email":         req.Email,
		"password_hash": hashedPwd,
		"username":      req.Username,
		"phone":         req.Phone,
		"full_name":     req.FullName,
		"gender":        req.Gender,
		"birth_date":    req.BirthDate,
		"role":          consts.RoleUser,
	}).InsertAndGetId()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToCreate)
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
		return nil, gerror.NewCode(consts.CodeInvalidInput)
	}

	userRecord, err := dao.Users.Ctx(ctx).Where("email", req.Account).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError)
	}
	fmt.Println(userRecord)
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound)
	}

	password := userRecord["password_hash"].String()
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)); err != nil {
		return nil, gerror.NewCode(consts.CodeIncorrectPassword)
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
		return nil, gerror.NewCode(consts.CodeFailedToCreate)
	}

	accessToken := &service.AccessToken{
		Iss: "parkin-ai-system",
		Sub: fmt.Sprintf("%d", userId),
		Exp: 0,
	}

	accessTokenStr, err := accessToken.Gen()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToCreate)
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
		return nil, gerror.NewCode(consts.CodeInvalidInput)
	}

	tokenRecord, err := dao.ApiTokens.Ctx(ctx).Where("token", req.RefreshToken).Where("is_active", true).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError)
	}
	if tokenRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeInvalidToken)
	}

	userId := tokenRecord["user_id"].Int64()

	accessToken := &service.AccessToken{
		Iss: "parkin-ai-system",
		Sub: fmt.Sprintf("%d", userId),
		Exp: 0,
	}

	accessTokenStr, err := accessToken.Gen()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToCreate)
	}

	newRefreshTokenStr := guid.S()

	_, err = dao.ApiTokens.Ctx(ctx).Where("token", req.RefreshToken).Data(g.Map{
		"is_active": false,
	}).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToUpdate)
	}

	_, err = dao.ApiTokens.Ctx(ctx).Data(g.Map{
		"user_id":     userId,
		"token":       newRefreshTokenStr,
		"description": "Refreshed token",
		"is_active":   true,
	}).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToCreate)
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
	userIDStr := ""
	if v := g.RequestFromCtx(ctx).GetCtxVar("user_id"); v != nil {
		userIDStr = v.String()
	}
	if userIDStr == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized)
	}

	userRecord, err := dao.Users.Ctx(ctx).Where("id", userIDStr).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError)
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound)
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

func (s *sUser) UserById(ctx context.Context, req *user.UserByIdReq) (res *user.UserByIdRes, err error) {
	userRecord, err := dao.Users.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError)
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound)
	}
	res = &user.UserByIdRes{
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

func (s *sUser) UserUpdateProfile(ctx context.Context, req *user.UserUpdateProfileReq) (res *user.UserUpdateProfileRes, err error) {
	userIDStr := ""
	if v := g.RequestFromCtx(ctx).GetCtxVar("user_id"); v != nil {
		userIDStr = v.String()
	}
	if userIDStr == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized)
	}

	_, err = dao.Users.Ctx(ctx).Where("id", userIDStr).Data(g.Map{
		"full_name":  req.FullName,
		"phone":      req.Phone,
		"gender":     req.Gender,
		"birth_date": req.BirthDate,
	}).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToUpdate)
	}

	res = &user.UserUpdateProfileRes{
		Success: true,
		Message: "Profile updated successfully",
	}
	return
}

func (s *sUser) GetAllUsers(ctx context.Context, req *user.GetAllUsersReq) (res *user.GetAllUsersRes, err error) {
	offset := (req.Page - 1) * req.Size
	total, err := dao.Users.Ctx(ctx).Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError)
	}
	users, err := dao.Users.Ctx(ctx).
		Fields("id", "username", "email", "phone", "full_name", "role", "created_at").
		Offset(offset).
		Limit(req.Size).
		Order("created_at DESC").
		All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError)
	}
	userList := make([]user.UserInfo, 0, len(users))
	for _, u := range users {
		userList = append(userList, user.UserInfo{
			UserID:    u["id"].Int64(),
			Username:  u["username"].String(),
			Email:     u["email"].String(),
			Phone:     u["phone"].String(),
			FullName:  u["full_name"].String(),
			Role:      u["role"].String(),
			CreatedAt: u["created_at"].String(),
		})
	}
	res = &user.GetAllUsersRes{
		Users: userList,
		Total: int(total),
		Page:  req.Page,
		Size:  req.Size,
	}
	return
}

func (s *sUser) DeleteUser(ctx context.Context, req *user.DeleteUserReq) (res *user.DeleteUserRes, err error) {
	userRecord, err := dao.Users.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, err
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound)
	}
	currentUserID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if currentUserID == userRecord["id"].String() {
		return nil, gerror.NewCode(consts.CodeCannotDeleteSelf)
	}
	_, err = dao.Users.Ctx(ctx).Where("id", req.Id).Delete()
	if err != nil {
		return nil, err
	}
	_, err = dao.ApiTokens.Ctx(ctx).Where("user_id", req.Id).Delete()
	if err != nil {
		g.Log().Warning(ctx, "Failed to delete user tokens:", err)
	}
	res = &user.DeleteUserRes{
		Message: "User deleted successfully",
	}
	return
}

func (s *sUser) UpdateUserRole(ctx context.Context, req *user.UpdateUserRoleReq) (res *user.UpdateUserRoleRes, err error) {
	userRecord, err := dao.Users.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, err
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound)
	}
	currentUserID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if currentUserID == userRecord["id"].String() {
		return nil, gerror.NewCode(consts.CodeNotAdmin)
	}
	_, err = dao.Users.Ctx(ctx).
		Where("id", req.Id).
		Data(g.Map{"role": req.Role}).
		Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToUpdate)
	}
	_, err = dao.ApiTokens.Ctx(ctx).
		Where("user_id", req.Id).
		Data(g.Map{"is_active": false}).
		Update()
	if err != nil {
		g.Log().Warning(ctx, "Failed to invalidate user tokens:", err)
	}
	res = &user.UpdateUserRoleRes{
		Message: "User role updated successfully",
	}
	return
}

func (s *sUser) UserUpdateProfileWithRBAC(ctx context.Context, req *user.UserUpdateProfileReq) (res *user.UserUpdateProfileRes, err error) {
	currentUserID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if !middleware.CheckResourceOwnership(g.RequestFromCtx(ctx), currentUserID) {
		return nil, gerror.NewCode(consts.CodeNotOwner)
	}
	// ... existing update logic here ...
	res = &user.UserUpdateProfileRes{
		Message: "Profile updated successfully",
	}
	return
}
