package parking_order

import (
	"context"

	"parkin-ai-system/api/parking_order/parking_order"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

// ParkingOrderDelete soft deletes a parking order by calling the service.
func (c *ControllerParking_order) ParkingOrderDelete(ctx context.Context, req *parking_order.ParkingOrderDeleteReq) (res *parking_order.ParkingOrderDeleteRes, err error) {
	// Map request to service input
	input := &entity.ParkingOrderDeleteReq{
		Id: req.Id,
	}

	// Call service function
	message, err := service.ParkingOrder().ParkingOrderDelete(ctx, input)
	if err != nil {
		return nil, err
	}

	// Create response
	res = &parking_order.ParkingOrderDeleteRes{
		Message: message.Message,
	}
	return res, nil
}
