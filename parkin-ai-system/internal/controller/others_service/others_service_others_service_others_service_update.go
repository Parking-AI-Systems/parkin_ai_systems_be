package others_service

import (
	"context"

	"parkin-ai-system/api/others_service/others_service"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerOthers_service) OthersServiceUpdate(ctx context.Context, req *others_service.OthersServiceUpdateReq) (res *others_service.OthersServiceUpdateRes, err error) {
	// Map API request to entity request
	input := &entity.OthersServiceUpdateReq{
		Id:              req.Id,
		LotId:           req.LotId,
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		DurationMinutes: req.DurationMinutes,
		IsActive:        req.IsActive,
	}

	// Call service
	updateRes, err := service.OthersService().OthersServiceUpdate(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &others_service.OthersServiceUpdateRes{
		Service: entityToApiOthersServiceItem(updateRes),
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
