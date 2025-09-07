package parking_order

import (
	"context"

	"parkin-ai-system/api/parking_order/parking_order"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

// ParkingOrderCancel cancels a parking order by calling the service.
func (c *ControllerParking_order) ParkingOrderCancel(ctx context.Context, req *parking_order.ParkingOrderCancelReq) (res *parking_order.ParkingOrderCancelRes, err error) {
	// Map request to service input
	input := &entity.ParkingOrderCancelReq{
		Id: req.Id,
	}

	// Call service function
	order, err := service.ParkingOrder().ParkingOrderCancel(ctx, input)
	if err != nil {
		return nil, err
	}

	// Create response
	res = &parking_order.ParkingOrderCancelRes{
		Order: entityToApiParkingOrderItem(order),
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}

func entityToApiParkingOrderItem(item *entity.ParkingOrderItem) parking_order.ParkingOrderItem {
	return parking_order.ParkingOrderItem{
		Id:            item.Id,
		UserId:        item.UserId,
		VehicleId:     item.VehicleId,
		LotName:       item.LotName,
		SlotCode:      item.SlotCode,
		LotId:         item.LotId,
		SlotId:        item.SlotId,
		VehiclePlate:  item.VehiclePlate,
		Price:         item.Price,
		StartTime:     item.StartTime,
		EndTime:       item.EndTime,
		Status:        item.Status,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
		PaymentStatus: item.PaymentStatus,
	}
}
