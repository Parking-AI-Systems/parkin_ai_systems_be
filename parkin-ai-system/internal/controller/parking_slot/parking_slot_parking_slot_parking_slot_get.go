package parking_slot

import (
	"context"

	"parkin-ai-system/api/parking_slot/parking_slot"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerParking_slot) ParkingSlotGet(ctx context.Context, req *parking_slot.ParkingSlotGetReq) (res *parking_slot.ParkingSlotGetRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingSlotGetReq{
		Id: req.Id,
	}

	// Call service
	slot, err := service.ParkingSlot().ParkingSlotGet(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_slot.ParkingSlotGetRes{
		Slot: parking_slot.ParkingSlotItem{
			Id:          slot.Id,
			LotId:       slot.LotId,
			Code:        slot.Code,
			IsAvailable: slot.IsAvailable,
			SlotType:    slot.SlotType,
			Floor:       slot.Floor,
			CreatedAt:   slot.CreatedAt,
		},
	}
	return res, nil
}
