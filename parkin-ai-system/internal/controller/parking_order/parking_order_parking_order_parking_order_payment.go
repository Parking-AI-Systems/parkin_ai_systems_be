package parking_order

import (
	"context"

	"parkin-ai-system/api/parking_order/parking_order"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerParking_order) ParkingOrderPayment(ctx context.Context, req *parking_order.ParkingOrderPaymentReq) (res *parking_order.ParkingOrderPaymentRes, err error) {
	// Map request to service input
	input := &entity.ParkingOrderPaymentReq{
		Id: req.Id,
	}

	// Call service function
	paymentRes, err := service.ParkingOrder().ParkingOrderPayment(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_order.ParkingOrderPaymentRes{
		Order: entityToApiParkingOrderItem(paymentRes),
	}
	return res, nil
}
