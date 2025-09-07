package user

import (
	"context"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerUser) UserRoleDistribution(ctx context.Context, req *user.UserRoleDistributionReq) (res *user.UserRoleDistributionRes, err error) {
	// Map API request to entity request
	input := &entity.UserRoleDistributionReq{
		Period:    req.Period,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// Call service
	distRes, err := service.User().UserRoleDistribution(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &user.UserRoleDistributionRes{
		Roles: make([]user.UserRoleItem, 0, len(distRes.Roles)),
	}
	for _, item := range distRes.Roles {
		res.Roles = append(res.Roles, user.UserRoleItem{
			Role:  item.Role,
			Count: item.Count,
		})
	}

	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
