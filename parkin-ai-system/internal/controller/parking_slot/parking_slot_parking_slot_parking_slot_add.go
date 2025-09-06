package parking_slot

import (
	"context"

	"parkin-ai-system/api/parking_slot/parking_slot"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerParking_slot) ParkingSlotAdd(ctx context.Context, req *parking_slot.ParkingSlotAddReq) (res *parking_slot.ParkingSlotAddRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingSlotAddReq{
		LotId:       req.LotId,
		Code:        req.Code,
		IsAvailable: req.IsAvailable,
		SlotType:    req.SlotType,
		Floor:       req.Floor,
	}

	// Call service
	addRes, err := service.ParkingSlot().ParkingSlotAdd(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_slot.ParkingSlotAddRes{
		Id: addRes.Id,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
