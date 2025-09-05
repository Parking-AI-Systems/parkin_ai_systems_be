package parking_order

import (
	"context"

	"parkin-ai-system/api/parking_order/parking_order"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

// ParkingOrderAdd creates a new parking order by calling the service.
func (c *ControllerParking_order) ParkingOrderAdd(ctx context.Context, req *parking_order.ParkingOrderAddReq) (res *parking_order.ParkingOrderAddRes, err error) {
	// Map request to service input
	input := &entity.ParkingOrderAddReq{
		VehicleId: req.VehicleId,
		LotId:     req.LotId,
		SlotId:    req.SlotId,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	orderId, err := service.ParkingOrder().ParkingOrderAddWithUser(ctx, input)
	if err != nil {
		return nil, err
	}

	// Create response
	res = &parking_order.ParkingOrderAddRes{
		Id: orderId.Id,
	}
	return res, nil
}
