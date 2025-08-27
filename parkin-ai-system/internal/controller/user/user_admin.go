package user

import (
	"context"
	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/service"
)

func (c *ControllerUser) GetAllUsers(ctx context.Context, req *user.GetAllUsersReq) (res *user.GetAllUsersRes, err error) {
	return service.User().GetAllUsers(ctx, req)
}

func (c *ControllerUser) DeleteUser(ctx context.Context, req *user.DeleteUserReq) (res *user.DeleteUserRes, err error) {
	return service.User().DeleteUser(ctx, req)
}

func (c *ControllerUser) UpdateUserRole(ctx context.Context, req *user.UpdateUserRoleReq) (res *user.UpdateUserRoleRes, err error) {
	return service.User().UpdateUserRole(ctx, req)
}

func (c *ControllerUser) UserUpdateProfileWithRBAC(ctx context.Context, req *user.UserUpdateProfileReq) (res *user.UserUpdateProfileRes, err error) {
	return service.User().UserUpdateProfileWithRBAC(ctx, req)
}
