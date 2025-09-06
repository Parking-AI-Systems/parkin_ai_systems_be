package others_service

import (
	"context"

	"parkin-ai-system/api/others_service/others_service"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerOthers_service) OthersServiceGet(ctx context.Context, req *others_service.OthersServiceGetReq) (res *others_service.OthersServiceGetRes, err error) {
	// Map API request to entity request
	input := &entity.OthersServiceGetReq{
		Id: req.Id,
	}

	// Call service
	serviceInfo, err := service.OthersService().OthersServiceGet(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &others_service.OthersServiceGetRes{
		Service: entityToApiOthersServiceItem(serviceInfo),
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
func entityToApiOthersServiceItem(e *entity.OthersServiceItem) others_service.OthersServiceItem {
	if e == nil {
		return others_service.OthersServiceItem{}
	}
	return others_service.OthersServiceItem{
		Id:              e.Id,
		LotId:           e.LotId,
		Name:            e.Name,
		Description:     e.Description,
		Price:           e.Price,
		DurationMinutes: e.DurationMinutes,
		IsActive:        e.IsActive,
		CreatedAt:       e.CreatedAt,
	}
}
