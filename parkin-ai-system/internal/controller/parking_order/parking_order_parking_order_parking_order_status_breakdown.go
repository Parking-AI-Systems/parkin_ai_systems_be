package parking_order

import (
	"context"

	"parkin-ai-system/api/parking_order/parking_order"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerParking_order) ParkingOrderStatusBreakdown(ctx context.Context, req *parking_order.ParkingOrderStatusBreakdownReq) (res *parking_order.ParkingOrderStatusBreakdownRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingOrderStatusBreakdownReq{
		Period:    req.Period,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// Call service
	breakdownRes, err := service.ParkingOrder().ParkingOrderStatusBreakdown(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_order.ParkingOrderStatusBreakdownRes{
		Statuses: make([]parking_order.ParkingOrderStatusItem, 0, len(breakdownRes.Statuses)),
		Total:    breakdownRes.Total,
	}
	for _, item := range breakdownRes.Statuses {
		res.Statuses = append(res.Statuses, parking_order.ParkingOrderStatusItem{
			Status: item.Status,
			Count:  item.Count,
		})
	}

	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
