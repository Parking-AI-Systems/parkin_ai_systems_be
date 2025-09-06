package other_service_orders

import (
	"context"

	"parkin-ai-system/api/other_service_orders/other_service_orders"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerOther_service_orders) OthersServiceOrderList(ctx context.Context, req *other_service_orders.OthersServiceOrderListReq) (res *other_service_orders.OthersServiceOrderListRes, err error) {
	// Map API request to entity request
	input := &entity.OthersServiceOrderListReq{
		UserId:   req.UserId,
		LotId:    req.LotId,
		Page:     req.Page,
		PageSize: req.PageSize,
		Status:   req.Status,
	}

	// Call service
	listRes, err := service.OthersServiceOrder().OthersServiceOrderList(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity list to API response
	res = &other_service_orders.OthersServiceOrderListRes{
		List:  make([]other_service_orders.OthersServiceOrderItem, 0, len(listRes.List)),
		Total: listRes.Total,
	}
	for _, item := range listRes.List {
		res.List = append(res.List, entityToApiServiceOrderItem(&item))
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
