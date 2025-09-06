package user

import (
	"context"
	"fmt"
	"parkin-ai-system/internal/consts"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/middleware"
	"parkin-ai-system/internal/model/do"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
	"parkin-ai-system/utility"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"

	"github.com/gogf/gf/v2/errors/gerror"
	"golang.org/x/crypto/bcrypt"
)

type sUser struct{}

func Init() {
	service.RegisterUser(&sUser{})
}

func init() {
	Init()
}

func (s *sUser) SignUp(ctx context.Context, req *entity.UserRegisterReq) (res *entity.UserRegisterRes, err error) {
	if req.Email == "" || req.Password == "" || req.Username == "" || req.Phone == "" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Email, username, password, and phone are required")
	}

	// Validate unique fields
	count, err := dao.Users.Ctx(ctx).
		WhereOr("email", req.Email).
		WhereOr("username", req.Username).
		WhereOr("phone", req.Phone).
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeEmailExists, "Username, email, or phone already exists")
	}

	hashedPwd, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeHashPasswordFailed, "Error hashing password")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error starting transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	userId, err := dao.Users.Ctx(ctx).TX(tx).Data(g.Map{
		"email":          req.Email,
		"password_hash":  hashedPwd,
		"username":       req.Username,
		"phone":          req.Phone,
		"full_name":      req.FullName,
		"gender":         req.Gender,
		"birth_date":     req.BirthDate,
		"role":           consts.RoleUser,
		"avatar_url":     req.AvatarUrl,
		"wallet_balance": 0.0,
		"created_at":     gtime.Now(),
		"updated_at":     gtime.Now(),
	}).InsertAndGetId()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToCreate, "Error creating user")
	}

	// Notify admins
	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving admins")
	}
	for _, admin := range adminUsers {
		notiData := do.Notifications{
			UserId:         admin["id"].Int64(),
			Type:           "user_registered",
			Content:        fmt.Sprintf("New user #%d (%s) registered.", userId, req.Username),
			RelatedOrderId: userId,
			IsRead:         false,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	res = &entity.UserRegisterRes{
		UserId:    userId,
		Username:  req.Username,
		Email:     req.Email,
		Phone:     req.Phone,
		FullName:  req.FullName,
		Gender:    req.Gender,
		BirthDate: req.BirthDate,
	}
	return
}

func (s *sUser) Login(ctx context.Context, req *entity.UserLoginReq) (res *entity.UserLoginRes, err error) {
	if req.Account == "" || req.Password == "" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Account and password are required")
	}

	userRecord, err := dao.Users.Ctx(ctx).Where("email", req.Account).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	password := userRecord["password_hash"].String()
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)); err != nil {
		return nil, gerror.NewCode(consts.CodeIncorrectPassword, "Incorrect password")
	}

	userId := userRecord["id"].Int64()
	refreshTokenStr := guid.S()

	_, err = dao.ApiTokens.Ctx(ctx).Data(g.Map{
		"user_id":     userId,
		"token":       refreshTokenStr,
		"description": "Login refresh token",
		"is_active":   true,
		"created_at":  gtime.Now(),
	}).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToCreate, "Error creating refresh token")
	}

	accessToken := &service.AccessToken{
		Iss: "parkin-ai-system",
		Sub: fmt.Sprintf("%d", userId),
		Exp: 0, // Set appropriate expiry in production
	}

	accessTokenStr, err := accessToken.Gen()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToCreate, "Error generating access token")
	}

	res = &entity.UserLoginRes{
		AccessToken:   accessTokenStr,
		RefreshToken:  refreshTokenStr,
		UserId:        userId,
		Username:      userRecord["username"].String(),
		Role:          userRecord["role"].String(),
		WalletBalance: userRecord["wallet_balance"].Float64(),
	}
	return
}

func (s *sUser) RefreshToken(ctx context.Context, req *entity.UserRefreshTokenReq) (res *entity.UserRefreshTokenRes, err error) {
	if req.RefreshToken == "" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Refresh token is required")
	}

	tokenRecord, err := dao.ApiTokens.Ctx(ctx).Where("token", req.RefreshToken).Where("is_active", true).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking token")
	}
	if tokenRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeInvalidToken, "Invalid or inactive refresh token")
	}

	userId := tokenRecord["user_id"].Int64()
	userRecord, err := dao.Users.Ctx(ctx).Where("id", userId).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	accessToken := &service.AccessToken{
		Iss: "parkin-ai-system",
		Sub: fmt.Sprintf("%d", userId),
		Exp: 0, // Set appropriate expiry in production
	}

	accessTokenStr, err := accessToken.Gen()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToCreate, "Error generating access token")
	}

	newRefreshTokenStr := guid.S()

	_, err = dao.ApiTokens.Ctx(ctx).Where("token", req.RefreshToken).Data(g.Map{
		"is_active": false,
	}).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToUpdate, "Error deactivating old token")
	}

	_, err = dao.ApiTokens.Ctx(ctx).Data(g.Map{
		"user_id":     userId,
		"token":       newRefreshTokenStr,
		"description": "Refreshed token",
		"is_active":   true,
		"created_at":  gtime.Now(),
	}).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToCreate, "Error creating new refresh token")
	}

	res = &entity.UserRefreshTokenRes{
		AccessToken:  accessTokenStr,
		RefreshToken: newRefreshTokenStr,
	}
	return
}

func (s *sUser) Logout(ctx context.Context, req *entity.UserLogoutReq) (res *entity.UserLogoutRes, err error) {
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

	res = &entity.UserLogoutRes{
		Message: "Logout successful",
	}
	return
}

func (s *sUser) UserProfile(ctx context.Context, req *entity.UserProfileReq) (res *entity.UserProfileRes, err error) {
	userIDStr := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userIDStr == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	userRecord, err := dao.Users.Ctx(ctx).Where("id", userIDStr).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving user")
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	res = &entity.UserProfileRes{
		UserId:        userRecord["id"].Int64(),
		Username:      userRecord["username"].String(),
		Email:         userRecord["email"].String(),
		Phone:         userRecord["phone"].String(),
		FullName:      userRecord["full_name"].String(),
		Gender:        userRecord["gender"].String(),
		BirthDate:     userRecord["birth_date"].String(),
		Role:          userRecord["role"].String(),
		AvatarUrl:     userRecord["avatar_url"].String(),
		WalletBalance: userRecord["wallet_balance"].Float64(),
		CreatedAt:     userRecord["created_at"].Time().Format("2006-01-02 15:04:05"),
		UpdatedAt:     userRecord["updated_at"].Time().Format("2006-01-02 15:04:05"),
	}
	return
}

func (s *sUser) UserById(ctx context.Context, req *entity.UserByIdReq) (res *entity.UserByIdRes, err error) {
	userIDStr := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userIDStr == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	currentUser, err := dao.Users.Ctx(ctx).Where("id", userIDStr).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking current user")
	}
	if currentUser.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Current user not found")
	}

	isAdmin := currentUser["role"].String() == consts.RoleAdmin
	if !isAdmin && gconv.Int64(userIDStr) != req.Id {
		return nil, gerror.NewCode(consts.CodeNotOwner, "You can only view your own profile or must be an admin")
	}

	userRecord, err := dao.Users.Ctx(ctx).Where("id", req.Id).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving user")
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}

	res = &entity.UserByIdRes{
		UserId:        userRecord["id"].Int64(),
		Username:      userRecord["username"].String(),
		Email:         userRecord["email"].String(),
		Phone:         userRecord["phone"].String(),
		FullName:      userRecord["full_name"].String(),
		Gender:        userRecord["gender"].String(),
		BirthDate:     userRecord["birth_date"].String(),
		Role:          userRecord["role"].String(),
		AvatarUrl:     userRecord["avatar_url"].String(),
		WalletBalance: userRecord["wallet_balance"].Float64(),
		CreatedAt:     userRecord["created_at"].Time().Format("2006-01-02 15:04:05"),
		UpdatedAt:     userRecord["updated_at"].Time().Format("2006-01-02 15:04:05"),
	}
	return
}

func (s *sUser) UserUpdateProfile(ctx context.Context, req *entity.UserUpdateProfileReq) (res *entity.UserUpdateProfileRes, err error) {
	userIDStr := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userIDStr == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	if !middleware.CheckResourceOwnership(g.RequestFromCtx(ctx), userIDStr) {
		return nil, gerror.NewCode(consts.CodeNotOwner, "You can only update your own profile")
	}

	// Validate unique fields
	if req.Email != "" || req.Phone != "" || req.Username != "" {
		m := dao.Users.Ctx(ctx).Where("id <> ?", userIDStr)
		if req.Username != "" {
			m = m.WhereOr("username", req.Username)
		}
		if req.Email != "" {
			m = m.WhereOr("email", req.Email)
		}
		if req.Phone != "" {
			m = m.WhereOr("phone", req.Phone)
		}
		count, err := m.Count()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking unique fields")
		}
		if count > 0 {
			return nil, gerror.NewCode(consts.CodeEmailExists, "Username, email, or phone already exists")
		}
	}

	data := g.Map{
		"updated_at": gtime.Now(),
	}
	if req.Username != "" {
		data["username"] = req.Username
	}
	if req.FullName != "" {
		data["full_name"] = req.FullName
	}
	if req.Email != "" {
		data["email"] = req.Email
	}
	if req.Phone != "" {
		data["phone"] = req.Phone
	}
	if req.Gender != "" {
		data["gender"] = req.Gender
	}
	if req.BirthDate != "" {
		data["birth_date"] = req.BirthDate
	}
	if req.AvatarUrl != "" {
		data["avatar_url"] = req.AvatarUrl
	}
	if req.Password != "" {
		hashedPwd, err := s.HashPassword(req.Password)
		if err != nil {
			return nil, gerror.NewCode(consts.CodeHashPasswordFailed, "Error hashing password")
		}
		data["password_hash"] = hashedPwd
	}

	_, err = dao.Users.Ctx(ctx).Where("id", userIDStr).Data(data).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToUpdate, "Error updating profile")
	}

	res = &entity.UserUpdateProfileRes{
		Success: true,
		Message: "Profile updated successfully",
	}
	return
}

func (s *sUser) GetAllUsers(ctx context.Context, req *entity.UserListReq) (res *entity.UserListRes, err error) {
	userIDStr := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userIDStr == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	currentUser, err := dao.Users.Ctx(ctx).Where("id", userIDStr).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if currentUser.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	if currentUser["role"].String() != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeNotAdmin, "Only admins can list users")
	}

	m := dao.Users.Ctx(ctx)
	if req.Username != "" {
		m = m.WhereLike("username", "%"+req.Username+"%")
	}
	if req.Email != "" {
		m = m.WhereLike("email", "%"+req.Email+"%")
	}
	if req.Role != "" {
		m = m.Where("role", req.Role)
	}

	total, err := m.Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error counting users")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	offset := (req.Page - 1) * req.PageSize

	users, err := dao.Users.Ctx(ctx).
		Fields("id", "username", "email", "phone", "full_name", "gender", "birth_date", "role", "avatar_url", "wallet_balance", "created_at", "updated_at").
		Offset(offset).
		Limit(req.PageSize).
		Order("created_at DESC").
		All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving users")
	}

	userList := make([]entity.UserItem, 0, len(users))
	for _, u := range users {
		userList = append(userList, entity.UserItem{
			UserId:        u["id"].Int64(),
			Username:      u["username"].String(),
			Email:         u["email"].String(),
			Phone:         u["phone"].String(),
			FullName:      u["full_name"].String(),
			Gender:        u["gender"].String(),
			BirthDate:     u["birth_date"].String(),
			Role:          u["role"].String(),
			AvatarUrl:     u["avatar_url"].String(),
			WalletBalance: u["wallet_balance"].Float64(),
			CreatedAt:     u["created_at"].Time().Format("2006-01-02 15:04:05"),
			UpdatedAt:     u["updated_at"].Time().Format("2006-01-02 15:04:05"),
		})
	}

	res = &entity.UserListRes{
		Users: userList,
		Total: total,
		Page:  req.Page,
		Size:  req.PageSize,
	}
	return
}

func (s *sUser) DeleteUser(ctx context.Context, req *entity.UserDeleteReq) (res *entity.UserDeleteRes, err error) {
	userIDStr := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userIDStr == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	currentUser, err := dao.Users.Ctx(ctx).Where("id", userIDStr).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if currentUser.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	if currentUser["role"].String() != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeNotAdmin, "Only admins can delete users")
	}

	userRecord, err := dao.Users.Ctx(ctx).Where("id", req.UserId).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking target user")
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Target user not found")
	}
	if userRecord["id"].String() == userIDStr {
		return nil, gerror.NewCode(consts.CodeCannotDeleteSelf, "Cannot delete your own account")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error starting transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Delete related data
	_, err = dao.WalletTransactions.Ctx(ctx).TX(tx).Where("user_id", req.UserId).Delete()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error deleting wallet transactions")
	}
	_, err = dao.Vehicles.Ctx(ctx).TX(tx).Where("user_id", req.UserId).Delete()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error deleting vehicles")
	}
	_, err = dao.Favorites.Ctx(ctx).TX(tx).Where("user_id", req.UserId).Delete()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error deleting favorites")
	}
	_, err = dao.ParkingOrders.Ctx(ctx).TX(tx).Where("user_id", req.UserId).Delete()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error deleting parking orders")
	}
	_, err = dao.OthersServiceOrders.Ctx(ctx).TX(tx).Where("user_id", req.UserId).Delete()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error deleting service orders")
	}
	_, err = dao.ApiTokens.Ctx(ctx).TX(tx).Where("user_id", req.UserId).Delete()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error deleting user tokens")
	}

	// Delete user
	_, err = dao.Users.Ctx(ctx).TX(tx).Where("id", req.UserId).Delete()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error deleting user")
	}

	// Notify admins
	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving admins")
	}
	for _, admin := range adminUsers {
		notiData := do.Notifications{
			UserId:         admin["id"].Int64(),
			Type:           "user_deleted",
			Content:        fmt.Sprintf("User #%d (%s) deleted.", req.UserId, userRecord["username"].String()),
			RelatedOrderId: req.UserId,
			IsRead:         false,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	res = &entity.UserDeleteRes{
		Message: "User deleted successfully",
	}
	return
}

func (s *sUser) UpdateUserRole(ctx context.Context, req *entity.UserUpdateRoleReq) (res *entity.UserUpdateRoleRes, err error) {
	userIDStr := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userIDStr == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	currentUser, err := dao.Users.Ctx(ctx).Where("id", userIDStr).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if currentUser.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	if currentUser["role"].String() != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeNotAdmin, "Only admins can update roles")
	}

	userRecord, err := dao.Users.Ctx(ctx).Where("id", req.UserId).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking target user")
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Target user not found")
	}
	if userRecord["id"].String() == userIDStr {
		return nil, gerror.NewCode(consts.CodeCannotDeleteSelf, "Cannot update your own role")
	}

	// Validate role
	isValidRole := false
	for _, validRole := range consts.ValidRoles {
		if req.Role == validRole {
			isValidRole = true
			break
		}
	}
	if !isValidRole {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Invalid role")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error starting transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = dao.Users.Ctx(ctx).TX(tx).Where("id", req.UserId).Data(g.Map{
		"role":       req.Role,
		"updated_at": gtime.Now(),
	}).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToUpdate, "Error updating user role")
	}

	_, err = dao.ApiTokens.Ctx(ctx).TX(tx).Where("user_id", req.UserId).Data(g.Map{
		"is_active": false,
	}).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error invalidating user tokens")
	}

	// Notify admins
	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving admins")
	}
	for _, admin := range adminUsers {
		notiData := do.Notifications{
			UserId:         admin["id"].Int64(),
			Type:           "user_role_updated",
			Content:        fmt.Sprintf("User #%d (%s) role updated to %s.", req.UserId, userRecord["username"].String(), req.Role),
			RelatedOrderId: req.UserId,
			IsRead:         false,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	res = &entity.UserUpdateRoleRes{
		Message: "User role updated successfully",
	}
	return
}

func (s *sUser) UpdateWalletBalance(ctx context.Context, req *entity.UserUpdateWalletBalanceReq) (res *entity.UserUpdateWalletBalanceRes, err error) {
	userIDStr := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userIDStr == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "User not authenticated")
	}

	currentUser, err := dao.Users.Ctx(ctx).Where("id", userIDStr).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking user")
	}
	if currentUser.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "User not found")
	}
	if currentUser["role"].String() != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeNotAdmin, "Only admins can update wallet balance")
	}

	userRecord, err := dao.Users.Ctx(ctx).Where("id", req.UserId).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error checking target user")
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Target user not found")
	}

	if req.WalletBalance < 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Wallet balance cannot be negative")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error starting transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	currentBalance := userRecord["wallet_balance"].Float64()
	amount := req.WalletBalance - currentBalance
	if amount != 0 {
		transactionType := consts.TransactionTypeDeposit
		if amount < 0 {
			transactionType = consts.TransactionTypeWithdrawal
		}
		wtData := do.WalletTransactions{
			UserId:         req.UserId,
			Amount:         amount,
			Type:           transactionType,
			Description:    fmt.Sprintf("Admin updated wallet balance for user #%d", req.UserId),
			RelatedOrderId: 0,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.WalletTransactions.Ctx(ctx).TX(tx).Data(wtData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating wallet transaction")
		}
	}

	_, err = dao.Users.Ctx(ctx).TX(tx).Where("id", req.UserId).Data(g.Map{
		"wallet_balance": req.WalletBalance,
		"updated_at":     gtime.Now(),
	}).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToUpdate, "Error updating wallet balance")
	}

	// Notify admins
	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error retrieving admins")
	}
	for _, admin := range adminUsers {
		notiData := do.Notifications{
			UserId:         admin["id"].Int64(),
			Type:           "user_wallet_updated",
			Content:        fmt.Sprintf("User #%d (%s) wallet balance updated to %.2f.", req.UserId, userRecord["username"].String(), req.WalletBalance),
			RelatedOrderId: req.UserId,
			IsRead:         false,
			CreatedAt:      gtime.Now(),
		}
		_, err = dao.Notifications.Ctx(ctx).TX(tx).Data(notiData).Insert()
		if err != nil {
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Error creating notification")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Error committing transaction")
	}

	res = &entity.UserUpdateWalletBalanceRes{
		Message: "Wallet balance updated successfully",
	}
	return
}

func (s *sUser) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
