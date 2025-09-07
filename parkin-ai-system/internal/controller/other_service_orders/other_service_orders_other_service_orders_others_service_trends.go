package other_service_orders

import (
	"context"

	"parkin-ai-system/api/other_service_orders/other_service_orders"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerOther_service_orders) OthersServiceTrends(ctx context.Context, req *other_service_orders.OthersServiceTrendsReq) (res *other_service_orders.OthersServiceTrendsRes, err error) {
	// Map API request to entity request
	input := &entity.OthersServiceTrendsReq{
		Period:    req.Period,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// Call service
	trendsRes, err := service.OthersServiceOrder().OthersServiceTrends(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &other_service_orders.OthersServiceTrendsRes{
		Orders: make([]other_service_orders.OthersServiceTrendsItem, 0, len(trendsRes.Orders)),
		Total:  trendsRes.Total,
	}
	for _, item := range trendsRes.Orders {
		res.Orders = append(res.Orders, other_service_orders.OthersServiceTrendsItem{
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
