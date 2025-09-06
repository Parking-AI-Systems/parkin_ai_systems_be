package parking_order

import (
	"context"

	"parkin-ai-system/api/parking_order/parking_order"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

// ParkingOrderGet retrieves a parking order by calling the service.
func (c *ControllerParking_order) ParkingOrderGet(ctx context.Context, req *parking_order.ParkingOrderGetReq) (res *parking_order.ParkingOrderGetRes, err error) {
	// Map request to service input
	input := &entity.ParkingOrderGetReq{
		Id: req.Id,
	}

	// Call service function
	order, err := service.ParkingOrder().ParkingOrderGet(ctx, input)
	if err != nil {
		return nil, err
	}

	// Create response with mapped order
	res = &parking_order.ParkingOrderGetRes{
		Order: entityToApiParkingOrderItem(order),
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
