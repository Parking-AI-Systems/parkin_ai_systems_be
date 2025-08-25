package user

import (
	"context"
	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/service"
)

func (c *ControllerUser) UserProfile(ctx context.Context, req *user.UserProfileReq) (res *user.UserProfileRes, err error) {
	return service.User().UserProfile(ctx, req)
}
