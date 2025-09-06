package user

import (
	"context"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
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
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}

	return res, nil
}
