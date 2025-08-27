// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"parkin-ai-system/api/user/user"
)

type (
	IUser interface {
		SignUp(ctx context.Context, req *user.RegisterReq) (res *user.RegisterRes, err error)
		Login(ctx context.Context, req *user.UserLoginReq) (res *user.UserLoginRes, err error)
		RefreshToken(ctx context.Context, req *user.RefreshTokenReq) (res *user.RefreshTokenRes, err error)
		Logout(ctx context.Context, req *user.UserLogoutReq) (res *user.UserLogoutRes, err error)
		HashPassword(password string) (string, error)
		UserProfile(ctx context.Context, req *user.UserProfileReq) (res *user.UserProfileRes, err error)
		UserById(ctx context.Context, req *user.UserByIdReq) (res *user.UserByIdRes, err error)
		UserUpdateProfile(ctx context.Context, req *user.UserUpdateProfileReq) (res *user.UserUpdateProfileRes, err error)
		GetAllUsers(ctx context.Context, req *user.GetAllUsersReq) (res *user.GetAllUsersRes, err error)
		DeleteUser(ctx context.Context, req *user.DeleteUserReq) (res *user.DeleteUserRes, err error)
		UpdateUserRole(ctx context.Context, req *user.UpdateUserRoleReq) (res *user.UpdateUserRoleRes, err error)
		UserUpdateProfileWithRBAC(ctx context.Context, req *user.UserUpdateProfileReq) (res *user.UserUpdateProfileRes, err error)
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
