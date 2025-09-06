package other_service_orders

import (
	"context"

	"parkin-ai-system/api/other_service_orders/other_service_orders"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerOther_service_orders) OthersServiceOrderGet(ctx context.Context, req *other_service_orders.OthersServiceOrderGetReq) (res *other_service_orders.OthersServiceOrderGetRes, err error) {
	// Map API request to entity request
	input := &entity.OthersServiceOrderGetReq{
		Id: req.Id,
	}

	// Call service
	order, err := service.OthersServiceOrder().OthersServiceOrderGet(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &other_service_orders.OthersServiceOrderGetRes{
		Order: entityToApiServiceOrderItem(order),
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
