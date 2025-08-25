package user

import (
	"context"
	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/service"
)

func (c *ControllerUser) UserProfile(ctx context.Context, req *user.UserProfileReq) (res *user.UserProfileRes, err error) {
	res, err = service.User().UserProfile(ctx, req)
	if err != nil {
		return nil, err
	}
	return
}
