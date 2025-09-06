package parking_slot

import (
	"context"

	"parkin-ai-system/api/parking_slot/parking_slot"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerParking_slot) ParkingSlotDelete(ctx context.Context, req *parking_slot.ParkingSlotDeleteReq) (res *parking_slot.ParkingSlotDeleteRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingSlotDeleteReq{
		Id: req.Id,
	}

	// Call service
	deleteRes, err := service.ParkingSlot().ParkingSlotDelete(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_slot.ParkingSlotDeleteRes{
		Message: deleteRes.Message,
	}
	return res, nil
}
