package parking_order

import (
	"context"

	"parkin-ai-system/api/parking_order/parking_order"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerParking_order) ParkingOrderUpdate(ctx context.Context, req *parking_order.ParkingOrderUpdateReq) (res *parking_order.ParkingOrderUpdateRes, err error) {
	// Map request to service input
	input := &entity.ParkingOrderUpdateReq{
		Id:        req.Id,
		Status:    req.Status,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// Call service function
	updateRes, err := service.ParkingOrder().ParkingOrderUpdate(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_order.ParkingOrderUpdateRes{
		Order: entityToApiParkingOrderItem(updateRes),
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
