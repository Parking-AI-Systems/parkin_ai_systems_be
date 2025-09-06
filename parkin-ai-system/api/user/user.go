// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package user

import (
	"context"

	"parkin-ai-system/api/user/user"
)

type IUserUser interface {
	UserRegister(ctx context.Context, req *user.UserRegisterReq) (res *user.UserRegisterRes, err error)
	UserLogin(ctx context.Context, req *user.UserLoginReq) (res *user.UserLoginRes, err error)
	UserRefreshToken(ctx context.Context, req *user.UserRefreshTokenReq) (res *user.UserRefreshTokenRes, err error)
	UserLogout(ctx context.Context, req *user.UserLogoutReq) (res *user.UserLogoutRes, err error)
	UserProfile(ctx context.Context, req *user.UserProfileReq) (res *user.UserProfileRes, err error)
	UserById(ctx context.Context, req *user.UserByIdReq) (res *user.UserByIdRes, err error)
	UserUpdateProfile(ctx context.Context, req *user.UserUpdateProfileReq) (res *user.UserUpdateProfileRes, err error)
	UserList(ctx context.Context, req *user.UserListReq) (res *user.UserListRes, err error)
	UserDelete(ctx context.Context, req *user.UserDeleteReq) (res *user.UserDeleteRes, err error)
	UserUpdateRole(ctx context.Context, req *user.UserUpdateRoleReq) (res *user.UserUpdateRoleRes, err error)
	UserUpdateWalletBalance(ctx context.Context, req *user.UserUpdateWalletBalanceReq) (res *user.UserUpdateWalletBalanceRes, err error)
}
