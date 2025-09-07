package user

import (
	"context"

	"parkin-ai-system/api/user/user"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerUser) UserRecentRegistrations(ctx context.Context, req *user.UserRecentRegistrationsReq) (res *user.UserRecentRegistrationsRes, err error) {
	// Map API request to entity request
	input := &entity.UserRecentRegistrationsReq{
		Period:    req.Period,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// Call service
	regRes, err := service.User().UserRecentRegistrations(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &user.UserRecentRegistrationsRes{
		Registrations: make([]user.UserRegistrationItem, 0, len(regRes.Registrations)),
	}
	for _, item := range regRes.Registrations {
		res.Registrations = append(res.Registrations, user.UserRegistrationItem{
			Date:  item.Date,
			Count: item.Count,
		})
	}

	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
