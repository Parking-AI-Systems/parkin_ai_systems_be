package user

import (
	"context"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerUser) UserUpdateRole(ctx context.Context, req *user.UserUpdateRoleReq) (res *user.UserUpdateRoleRes, err error) {
	// Map API request to entity request
	input := &entity.UserUpdateRoleReq{
		UserId: req.UserId,
		Role:   req.Role,
	}

	// Call service
	updateRes, err := service.User().UpdateUserRole(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &user.UserUpdateRoleRes{
		Message: updateRes.Message,
	}
	return res, nil
}
