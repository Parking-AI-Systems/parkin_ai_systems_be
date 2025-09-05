package other_service_orders

import (
	"context"

	"parkin-ai-system/api/other_service_orders/other_service_orders"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerOther_service_orders) OthersServiceOrderDelete(ctx context.Context, req *other_service_orders.OthersServiceOrderDeleteReq) (res *other_service_orders.OthersServiceOrderDeleteRes, err error) {
	// Map API request to entity request
	input := &entity.OthersServiceOrderDeleteReq{
		Id: req.Id,
	}

	// Call service
	deleteRes, err := service.OthersServiceOrder().OthersServiceOrderDelete(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &other_service_orders.OthersServiceOrderDeleteRes{
		Message: deleteRes.Message,
	}
	return res, nil
}
