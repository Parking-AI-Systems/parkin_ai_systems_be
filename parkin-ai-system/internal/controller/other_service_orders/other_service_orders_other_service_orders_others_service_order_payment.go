package other_service_orders

import (
	"context"

	"parkin-ai-system/api/other_service_orders/other_service_orders"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerOther_service_orders) OthersServiceOrderPayment(ctx context.Context, req *other_service_orders.OthersServiceOrderPaymentReq) (res *other_service_orders.OthersServiceOrderPaymentRes, err error) {
	// Map API request to entity request
	input := &entity.OthersServiceOrderPaymentReq{
		Id: req.Id,
	}

	// Call service
	paymentRes, err := service.OthersServiceOrder().OthersServiceOrderPayment(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &other_service_orders.OthersServiceOrderPaymentRes{
		Order: entityToApiServiceOrderItem(paymentRes),
	}
	return res, nil
}
