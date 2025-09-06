package user

import (
	"context"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerUser) UserLogout(ctx context.Context, req *user.UserLogoutReq) (res *user.UserLogoutRes, err error) {
	// Map API request to entity request
	input := &entity.UserLogoutReq{}

	// Call service
	logoutRes, err := service.User().Logout(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &user.UserLogoutRes{
		Message: logoutRes.Message,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}

	return res, nil
}
