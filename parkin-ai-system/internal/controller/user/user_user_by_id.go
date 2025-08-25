package user

import (
	"context"
	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/service"
)

func (c *ControllerUser) UserById(ctx context.Context, req *user.UserByIdReq) (res *user.UserByIdRes, err error) {
	res, err = service.User().UserById(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
