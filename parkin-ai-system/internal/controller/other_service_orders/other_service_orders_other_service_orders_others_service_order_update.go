package other_service_orders

import (
	"context"

	"parkin-ai-system/api/other_service_orders/other_service_orders"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerOther_service_orders) OthersServiceOrderUpdate(ctx context.Context, req *other_service_orders.OthersServiceOrderUpdateReq) (res *other_service_orders.OthersServiceOrderUpdateRes, err error) {
	// Map API request to entity request
	input := &entity.OthersServiceOrderUpdateReq{
		Id:            req.Id,
		ScheduledTime: req.ScheduledTime,
		Status:        req.Status,
	}

	// Call service
	updateRes, err := service.OthersServiceOrder().OthersServiceOrderUpdate(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &other_service_orders.OthersServiceOrderUpdateRes{
		Order: entityToApiServiceOrderItem(updateRes),
	}
	return res, nil
}
