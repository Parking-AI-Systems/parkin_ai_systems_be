package user

import (
	"context"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerUser) UserCount(ctx context.Context, req *user.UserCountReq) (res *user.UserCountRes, err error) {
	// Map API request to entity request
	input := &entity.UserCountReq{
		Period:    req.Period,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// Call service
	countRes, err := service.User().UserCount(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &user.UserCountRes{
		TotalUsers: countRes.TotalUsers,
	}

	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
