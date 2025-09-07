package other_service_orders

import (
	"context"

	"parkin-ai-system/api/other_service_orders/other_service_orders"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerOther_service_orders) OthersServiceOrderCancel(ctx context.Context, req *other_service_orders.OthersServiceOrderCancelReq) (res *other_service_orders.OthersServiceOrderCancelRes, err error) {
	// Map API request to entity request
	input := &entity.OthersServiceOrderCancelReq{
		Id: req.Id,
	}

	// Call service
	cancelRes, err := service.OthersServiceOrder().OthersServiceOrderCancel(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &other_service_orders.OthersServiceOrderCancelRes{
		Order: entityToApiServiceOrderItem(cancelRes),
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
func entityToApiServiceOrderItem(item *entity.OthersServiceOrderItem) other_service_orders.OthersServiceOrderItem {
	return other_service_orders.OthersServiceOrderItem{
		Id:            item.Id,
		VehicleId:     item.VehicleId,
		LotId:         item.LotId,
		ServiceId:     item.ServiceId,
		ScheduledTime: item.ScheduledTime,
		VehiclePlate:  item.VehiclePlate,
		Status:        item.Status,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
		Price:         item.Price,
		PaymentStatus: item.PaymentStatus,
		DeletedAt:     item.DeletedAt,
	}
}
