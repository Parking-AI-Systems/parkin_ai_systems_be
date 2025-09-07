package parking_order

import (
	"context"

	"parkin-ai-system/api/parking_order/parking_order"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerParking_order) ParkingOrderTrends(ctx context.Context, req *parking_order.ParkingOrderTrendsReq) (res *parking_order.ParkingOrderTrendsRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingOrderTrendsReq{
		Period:    req.Period,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// Call service
	trendsRes, err := service.ParkingOrder().ParkingOrderTrends(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_order.ParkingOrderTrendsRes{
		Orders: make([]parking_order.ParkingOrderTrendsItem, 0, len(trendsRes.Orders)),
		Total:  trendsRes.Total,
	}
	for _, item := range trendsRes.Orders {
		res.Orders = append(res.Orders, parking_order.ParkingOrderTrendsItem{
			Date:  item.Date,
			Count: item.Count,
		})
	}

	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
