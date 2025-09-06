package others_service

import (
	"context"

	"parkin-ai-system/api/others_service/others_service"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerOthers_service) OthersServiceAdd(ctx context.Context, req *others_service.OthersServiceAddReq) (res *others_service.OthersServiceAddRes, err error) {
	// Map API request to entity request
	input := &entity.OthersServiceAddReq{
		LotId:           req.LotId,
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		DurationMinutes: req.DurationMinutes,
		IsActive:        req.IsActive,
	}

	// Call service
	addRes, err := service.OthersService().OthersServiceAdd(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &others_service.OthersServiceAddRes{
		Id: addRes.Id,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
