package others_service

import (
	"context"

	"parkin-ai-system/api/others_service/others_service"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerOthers_service) OthersServiceList(ctx context.Context, req *others_service.OthersServiceListReq) (res *others_service.OthersServiceListRes, err error) {
	// Map API request to entity request
	input := &entity.OthersServiceListReq{
		LotId:    req.LotId,
		IsActive: req.IsActive,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// Call service
	listRes, err := service.OthersService().OthersServiceList(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity list to API response
	res = &others_service.OthersServiceListRes{
		List:  make([]others_service.OthersServiceItem, 0, len(listRes.List)),
		Total: listRes.Total,
	}
	for _, item := range listRes.List {
		res.List = append(res.List, entityToApiOthersServiceItem(&item))
	}
	return res, nil
}
