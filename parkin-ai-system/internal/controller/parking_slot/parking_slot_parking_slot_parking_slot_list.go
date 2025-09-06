package parking_slot

import (
	"context"

	"parkin-ai-system/api/parking_slot/parking_slot"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerParking_slot) ParkingSlotList(ctx context.Context, req *parking_slot.ParkingSlotListReq) (res *parking_slot.ParkingSlotListRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingSlotListReq{
		LotId:       req.LotId,
		IsAvailable: req.IsAvailable,
		SlotType:    req.SlotType,
		Page:        req.Page,
		PageSize:    req.PageSize,
	}

	// Call service
	listRes, err := service.ParkingSlot().ParkingSlotList(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity list to API response
	res = &parking_slot.ParkingSlotListRes{
		List:  make([]parking_slot.ParkingSlotItem, 0, len(listRes.List)),
		Total: listRes.Total,
	}
	for _, item := range listRes.List {
		res.List = append(res.List, parking_slot.ParkingSlotItem{
			Id:          item.Id,
			LotId:       item.LotId,
			Code:        item.Code,
			IsAvailable: item.IsAvailable,
			SlotType:    item.SlotType,
			Floor:       item.Floor,
			CreatedAt:   item.CreatedAt,
		})
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
