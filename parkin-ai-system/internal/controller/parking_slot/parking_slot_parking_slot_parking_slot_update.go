package parking_slot

import (
	"context"

	"parkin-ai-system/api/parking_slot/parking_slot"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerParking_slot) ParkingSlotUpdate(ctx context.Context, req *parking_slot.ParkingSlotUpdateReq) (res *parking_slot.ParkingSlotUpdateRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingSlotUpdateReq{
		Id:          req.Id,
		Code:        req.Code,
		IsAvailable: req.IsAvailable,
		SlotType:    req.SlotType,
		Floor:       req.Floor,
	}

	// Call service
	updateRes, err := service.ParkingSlot().ParkingSlotUpdate(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_slot.ParkingSlotUpdateRes{
		Slot: parking_slot.ParkingSlotItem{
			Id:          updateRes.Id,
			LotId:       updateRes.LotId,
			Code:        updateRes.Code,
			IsAvailable: updateRes.IsAvailable,
			SlotType:    updateRes.SlotType,
			Floor:       updateRes.Floor,
			CreatedAt:   updateRes.CreatedAt,
		},
	}
	return res, nil
}
