package other_service_orders

import (
	"context"

	"parkin-ai-system/api/other_service_orders/other_service_orders"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerOther_service_orders) OthersServiceOrderAdd(ctx context.Context, req *other_service_orders.OthersServiceOrderAddReq) (res *other_service_orders.OthersServiceOrderAddRes, err error) {
	// Map API request to entity request
	input := &entity.OthersServiceOrderAddReq{
		VehicleId:     req.VehicleId,
		LotId:         req.LotId,
		ServiceId:     req.ServiceId,
		ScheduledTime: req.ScheduledTime,
	}

	// Call service
	addRes, err := service.OthersServiceOrder().OthersServiceOrderAddWithUser(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &other_service_orders.OthersServiceOrderAddRes{
		Id: addRes.Id,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
