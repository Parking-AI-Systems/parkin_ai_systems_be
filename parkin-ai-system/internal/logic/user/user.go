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
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
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
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Please provide your email, username, password, and phone number.")
	}

	count, err := dao.Users.Ctx(ctx).
		WhereOr("email", req.Email).
		WhereOr("username", req.Username).
		WhereOr("phone", req.Phone).
		Where("deleted_at IS NULL").
		Count()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your information. Please try again later.")
	}
	if count > 0 {
		return nil, gerror.NewCode(consts.CodeEmailExists, "This username, email, or phone number is already in use.")
	}

	hashedPwd, err := s.HashPassword(req.Password)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeHashPasswordFailed, "Unable to process your password. Please try again.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while creating your account. Please try again later.")
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
		return nil, gerror.NewCode(consts.CodeFailedToCreate, "Something went wrong while creating your account. Please try again later.")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).Where("deleted_at IS NULL").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while creating your account. Please try again later.")
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
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while creating your account. Please try again later.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while creating your account. Please try again later.")
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
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Please provide your email and password.")
	}

	userRecord, err := dao.Users.Ctx(ctx).Where("email", req.Account).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "We couldn’t find an account with that email.")
	}

	password := userRecord["password_hash"].String()
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(req.Password)); err != nil {
		return nil, gerror.NewCode(consts.CodeIncorrectPassword, "The password you entered is incorrect.")
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
		return nil, gerror.NewCode(consts.CodeFailedToCreate, "Something went wrong during login. Please try again later.")
	}

	accessToken := &service.AccessToken{
		Iss: "parkin-ai-system",
		Sub: fmt.Sprintf("%d", userId),
		Exp: gtime.Now().Add(time.Hour * 24).Unix(),
	}

	accessTokenStr, err := accessToken.Gen()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToCreate, "Something went wrong during login. Please try again later.")
	}

	res = &entity.UserLoginRes{
		AccessToken:   accessTokenStr,
		RefreshToken:  refreshTokenStr,
		UserId:        userId,
		Username:      userRecord["username"].String(),
		Role:          userRecord["role"].String(),
		WalletBalance: userRecord["wallet_balance"].Float64(),
	}
	return res, nil
}

func (s *sUser) RefreshToken(ctx context.Context, req *entity.UserRefreshTokenReq) (res *entity.UserRefreshTokenRes, err error) {
	if req.RefreshToken == "" {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Please provide a refresh token.")
	}

	tokenRecord, err := dao.ApiTokens.Ctx(ctx).Where("token", req.RefreshToken).Where("is_active", true).One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify the token. Please try again.")
	}
	if tokenRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeInvalidToken, "The refresh token is invalid or no longer active.")
	}

	userId := tokenRecord["user_id"].Int64()
	userRecord, err := dao.Users.Ctx(ctx).Where("id", userId).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your account. Please try again later.")
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}

	accessToken := &service.AccessToken{
		Iss: "parkin-ai-system",
		Sub: fmt.Sprintf("%d", userId),
		Exp: gtime.Now().Add(time.Hour * 24).Unix(),
	}

	accessTokenStr, err := accessToken.Gen()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToCreate, "Something went wrong while refreshing your session. Please try again later.")
	}

	newRefreshTokenStr := guid.S()

	_, err = dao.ApiTokens.Ctx(ctx).Where("token", req.RefreshToken).Data(g.Map{
		"is_active": false,
	}).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToUpdate, "Something went wrong while refreshing your session. Please try again later.")
	}

	_, err = dao.ApiTokens.Ctx(ctx).Data(g.Map{
		"user_id":     userId,
		"token":       newRefreshTokenStr,
		"description": "Refreshed token",
		"is_active":   true,
		"created_at":  gtime.Now(),
	}).Insert()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToCreate, "Something went wrong while refreshing your session. Please try again later.")
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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view your profile.")
	}

	userRecord, err := dao.Users.Ctx(ctx).Where("id", userIDStr).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load your profile. Please try again later.")
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view the profile.")
	}

	currentUser, err := dao.Users.Ctx(ctx).Where("id", userIDStr).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load the profile. Please try again later.")
	}
	if currentUser.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "The account could not be found. Please contact support.")
	}

	isAdmin := currentUser["role"].String() == consts.RoleAdmin
	if !isAdmin && gconv.Int64(userIDStr) != req.Id {
		return nil, gerror.NewCode(consts.CodeNotOwner, "Only admins or the account owner can view this profile.")
	}

	userRecord, err := dao.Users.Ctx(ctx).Where("id", req.Id).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load the profile. Please try again later.")
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "The account could not be found. Please contact support.")
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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to update your profile.")
	}

	if !middleware.CheckResourceOwnership(g.RequestFromCtx(ctx), userIDStr) {
		return nil, gerror.NewCode(consts.CodeNotOwner, "You can only update your own profile.")
	}

	if req.Email != "" || req.Phone != "" || req.Username != "" {
		m := dao.Users.Ctx(ctx).Where("id <> ?", userIDStr).Where("deleted_at IS NULL")
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
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to verify your information. Please try again.")
		}
		if count > 0 {
			return nil, gerror.NewCode(consts.CodeEmailExists, "This username, email, or phone number is already in use.")
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
			return nil, gerror.NewCode(consts.CodeHashPasswordFailed, "Unable to process your password. Please try again.")
		}
		data["password_hash"] = hashedPwd
	}

	_, err = dao.Users.Ctx(ctx).Where("id", userIDStr).Where("deleted_at IS NULL").Data(data).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToUpdate, "Something went wrong while updating your profile. Please try again later.")
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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to view users.")
	}

	currentUser, err := dao.Users.Ctx(ctx).Where("id", userIDStr).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load users. Please try again later.")
	}
	if currentUser.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}
	if currentUser["role"].String() != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeNotAdmin, "Only admins can view the user list.")
	}

	m := dao.Users.Ctx(ctx).Where("deleted_at IS NULL")
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
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load users. Please try again later.")
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
		Where("deleted_at IS NULL").
		Offset(offset).
		Limit(req.PageSize).
		Order("created_at DESC").
		All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Unable to load users. Please try again later.")
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
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to delete a user.")
	}

	currentUser, err := dao.Users.Ctx(ctx).Where("id", userIDStr).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the user. Please try again later.")
	}
	if currentUser.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}
	if currentUser["role"].String() != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeNotAdmin, "Only admins can delete users.")
	}

	userRecord, err := dao.Users.Ctx(ctx).Where("id", req.UserId).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the user. Please try again later.")
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "The account could not be found. Please contact support.")
	}
	if userRecord["id"].String() == userIDStr {
		return nil, gerror.NewCode(consts.CodeCannotDeleteSelf, "You cannot delete your own account.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the user. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = dao.WalletTransactions.Ctx(ctx).TX(tx).Where("user_id", req.UserId).Data(g.Map{
		"deleted_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the user. Please try again later.")
	}
	_, err = dao.Vehicles.Ctx(ctx).TX(tx).Where("user_id", req.UserId).Data(g.Map{
		"deleted_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the user. Please try again later.")
	}
	_, err = dao.Favorites.Ctx(ctx).TX(tx).Where("user_id", req.UserId).Data(g.Map{
		"deleted_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the user. Please try again later.")
	}
	_, err = dao.ParkingOrders.Ctx(ctx).TX(tx).Where("user_id", req.UserId).Data(g.Map{
		"deleted_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the user. Please try again later.")
	}
	_, err = dao.OthersServiceOrders.Ctx(ctx).TX(tx).Where("user_id", req.UserId).Data(g.Map{
		"deleted_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the user. Please try again later.")
	}
	_, err = dao.ApiTokens.Ctx(ctx).TX(tx).Where("user_id", req.UserId).Delete()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the user. Please try again later.")
	}

	_, err = dao.Users.Ctx(ctx).TX(tx).Where("id", req.UserId).Data(g.Map{
		"deleted_at": gtime.Now(),
		"updated_at": gtime.Now(),
	}).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the user. Please try again later.")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).Where("deleted_at IS NULL").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the user. Please try again later.")
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
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the user. Please try again later.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while deleting the user. Please try again later.")
	}

	res = &entity.UserDeleteRes{
		Message: "User deleted successfully",
	}
	return
}

func (s *sUser) UpdateUserRole(ctx context.Context, req *entity.UserUpdateRoleReq) (res *entity.UserUpdateRoleRes, err error) {
	userIDStr := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userIDStr == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to update a user’s role.")
	}

	currentUser, err := dao.Users.Ctx(ctx).Where("id", userIDStr).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the user’s role. Please try again later.")
	}
	if currentUser.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}
	if currentUser["role"].String() != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeNotAdmin, "Only admins can update user roles.")
	}

	userRecord, err := dao.Users.Ctx(ctx).Where("id", req.UserId).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the user’s role. Please try again later.")
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "The account could not be found. Please contact support.")
	}
	if userRecord["id"].String() == userIDStr {
		return nil, gerror.NewCode(consts.CodeCannotDeleteSelf, "You cannot update your own role.")
	}

	isValidRole := false
	for _, validRole := range consts.ValidRoles {
		if req.Role == validRole {
			isValidRole = true
			break
		}
	}
	if !isValidRole {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "Please select a valid role.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the user’s role. Please try again later.")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	_, err = dao.Users.Ctx(ctx).TX(tx).Where("id", req.UserId).Where("deleted_at IS NULL").Data(g.Map{
		"role":       req.Role,
		"updated_at": gtime.Now(),
	}).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToUpdate, "Something went wrong while updating the user’s role. Please try again later.")
	}

	_, err = dao.ApiTokens.Ctx(ctx).TX(tx).Where("user_id", req.UserId).Data(g.Map{
		"is_active": false,
	}).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the user’s role. Please try again later.")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).Where("deleted_at IS NULL").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the user’s role. Please try again later.")
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
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the user’s role. Please try again later.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the user’s role. Please try again later.")
	}

	res = &entity.UserUpdateRoleRes{
		Message: "User role updated successfully",
	}
	return
}

func (s *sUser) UpdateWalletBalance(ctx context.Context, req *entity.UserUpdateWalletBalanceReq) (res *entity.UserUpdateWalletBalanceRes, err error) {
	userIDStr := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userIDStr == "" {
		return nil, gerror.NewCode(consts.CodeUnauthorized, "Please log in to update a user’s wallet balance.")
	}

	currentUser, err := dao.Users.Ctx(ctx).Where("id", userIDStr).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the wallet balance. Please try again later.")
	}
	if currentUser.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "Your account could not be found. Please contact support.")
	}
	if currentUser["role"].String() != consts.RoleAdmin {
		return nil, gerror.NewCode(consts.CodeNotAdmin, "Only admins can update wallet balances.")
	}

	userRecord, err := dao.Users.Ctx(ctx).Where("id", req.UserId).Where("deleted_at IS NULL").One()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the wallet balance. Please try again later.")
	}
	if userRecord.IsEmpty() {
		return nil, gerror.NewCode(consts.CodeUserNotFound, "The account could not be found. Please contact support.")
	}

	if req.WalletBalance < 0 {
		return nil, gerror.NewCode(consts.CodeInvalidInput, "The wallet balance cannot be negative.")
	}

	tx, err := g.DB().Begin(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the wallet balance. Please try again later.")
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
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the wallet balance. Please try again later.")
		}
	}

	_, err = dao.Users.Ctx(ctx).TX(tx).Where("id", req.UserId).Where("deleted_at IS NULL").Data(g.Map{
		"wallet_balance": req.WalletBalance,
		"updated_at":     gtime.Now(),
	}).Update()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeFailedToUpdate, "Something went wrong while updating the wallet balance. Please try again later.")
	}

	adminUsers, err := dao.Users.Ctx(ctx).Where("role", consts.RoleAdmin).Where("deleted_at IS NULL").All()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the wallet balance. Please try again later.")
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
			return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the wallet balance. Please try again later.")
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, gerror.NewCode(consts.CodeDatabaseError, "Something went wrong while updating the wallet balance. Please try again later.")
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
