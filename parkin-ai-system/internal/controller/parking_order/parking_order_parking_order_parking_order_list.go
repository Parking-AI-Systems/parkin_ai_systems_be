package parking_order

import (
	"context"

	"parkin-ai-system/api/parking_order/parking_order"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerParking_order) ParkingOrderList(ctx context.Context, req *parking_order.ParkingOrderListReq) (res *parking_order.ParkingOrderListRes, err error) {
	// Map request to service input
	input := &entity.ParkingOrderListReq{
		UserId:   req.UserId,
		LotId:    req.LotId,
		Page:     req.Page,
		PageSize: req.PageSize,
		Status:   req.Status,
	}

	// Call service function
	listRes, err := service.ParkingOrder().ParkingOrderList(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity list to API response
	res = &parking_order.ParkingOrderListRes{
		List:  make([]parking_order.ParkingOrderItem, 0, len(listRes.List)),
		Total: listRes.Total,
	}
	for _, item := range listRes.List {
		res.List = append(res.List, entityToApiParkingOrderItem(&item))
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
