package other_service_orders

import (
	"context"

	"parkin-ai-system/api/other_service_orders/other_service_orders"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerOther_service_orders) OthersServicePopular(ctx context.Context, req *other_service_orders.OthersServicePopularReq) (res *other_service_orders.OthersServicePopularRes, err error) {
	input := &entity.OthersServicePopularReq{
		Period:    req.Period,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	// Call service
	popularRes, err := service.OthersServiceOrder().OthersServicePopular(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &other_service_orders.OthersServicePopularRes{
		Services: make([]other_service_orders.OthersServicePopularItem, 0, len(popularRes.Services)),
	}
	for _, item := range popularRes.Services {
		res.Services = append(res.Services, other_service_orders.OthersServicePopularItem{
			ServiceId:  item.ServiceId,
			Name:       item.Name,
			OrderCount: item.OrderCount,
		})
	}

	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
