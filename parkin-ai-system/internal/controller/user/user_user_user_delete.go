package user

import (
	"context"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerUser) UserDelete(ctx context.Context, req *user.UserDeleteReq) (res *user.UserDeleteRes, err error) {
	// Map API request to entity request
	input := &entity.UserDeleteReq{
		UserId: req.UserId,
	}

	// Call service
	deleteRes, err := service.User().DeleteUser(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &user.UserDeleteRes{
		Message: deleteRes.Message,
	}
	return res, nil
}
