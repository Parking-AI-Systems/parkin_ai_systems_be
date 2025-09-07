// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"parkin-ai-system/internal/model/entity"
)

type (
	IUser interface {
		SignUp(ctx context.Context, req *entity.UserRegisterReq) (res *entity.UserRegisterRes, err error)
		Login(ctx context.Context, req *entity.UserLoginReq) (res *entity.UserLoginRes, err error)
		RefreshToken(ctx context.Context, req *entity.UserRefreshTokenReq) (res *entity.UserRefreshTokenRes, err error)
		Logout(ctx context.Context, req *entity.UserLogoutReq) (res *entity.UserLogoutRes, err error)
		UserProfile(ctx context.Context, req *entity.UserProfileReq) (res *entity.UserProfileRes, err error)
		UserById(ctx context.Context, req *entity.UserByIdReq) (res *entity.UserByIdRes, err error)
		UserUpdateProfile(ctx context.Context, req *entity.UserUpdateProfileReq) (res *entity.UserUpdateProfileRes, err error)
		GetAllUsers(ctx context.Context, req *entity.UserListReq) (res *entity.UserListRes, err error)
		DeleteUser(ctx context.Context, req *entity.UserDeleteReq) (res *entity.UserDeleteRes, err error)
		UpdateUserRole(ctx context.Context, req *entity.UserUpdateRoleReq) (res *entity.UserUpdateRoleRes, err error)
		UpdateWalletBalance(ctx context.Context, req *entity.UserUpdateWalletBalanceReq) (res *entity.UserUpdateWalletBalanceRes, err error)
		HashPassword(password string) (string, error)
		UserCount(ctx context.Context, req *entity.UserCountReq) (*entity.UserCountRes, error)
		UserRoleDistribution(ctx context.Context, req *entity.UserRoleDistributionReq) (*entity.UserRoleDistributionRes, error)
		UserRecentRegistrations(ctx context.Context, req *entity.UserRecentRegistrationsReq) (*entity.UserRecentRegistrationsRes, error)
	}
)

var (
	localUser IUser
)

func User() IUser {
	if localUser == nil {
		panic("implement not found for interface IUser, forgot register?")
	}
	return localUser
}

func RegisterUser(i IUser) {
	localUser = i
}
