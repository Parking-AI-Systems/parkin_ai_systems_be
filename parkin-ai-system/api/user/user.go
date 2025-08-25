// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package user

import (
	"context"

	"parkin-ai-system/api/user/user"
)

type IUserUser interface {
	RefreshToken(ctx context.Context, req *user.RefreshTokenReq) (res *user.RefreshTokenRes, err error)
	UserLogin(ctx context.Context, req *user.UserLoginReq) (res *user.UserLoginRes, err error)
	UserLogout(ctx context.Context, req *user.UserLogoutReq) (res *user.UserLogoutRes, err error)
	Register(ctx context.Context, req *user.RegisterReq) (res *user.RegisterRes, err error)
}
