package parking_order

import (
	"context"

	"parkin-ai-system/api/parking_order/parking_order"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerParking_order) ParkingOrderRevenue(ctx context.Context, req *parking_order.ParkingOrderRevenueReq) (res *parking_order.ParkingOrderRevenueRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingOrderRevenueReq{
		Period:    req.Period,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// Call service
	revenueRes, err := service.ParkingOrder().ParkingOrderRevenue(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_order.ParkingOrderRevenueRes{
		TotalRevenue: revenueRes.TotalRevenue,
	}

	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
