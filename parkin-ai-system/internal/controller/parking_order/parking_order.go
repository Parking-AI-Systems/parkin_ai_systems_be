package parking_order

import (
	"context"
	"parkin-ai-system/api/parking_order"
	"parkin-ai-system/internal/service"
)

type ControllerParkingOrder struct{}

func NewParkingOrder() *ControllerParkingOrder {
	return &ControllerParkingOrder{}
}

func (c *ControllerParkingOrder) ParkingOrderAdd(ctx context.Context, req *parking_order.ParkingOrderAddReq) (res *parking_order.ParkingOrderAddRes, err error) {
	return service.ParkingOrder().ParkingOrderAddWithUser(ctx, req)
}

func (c *ControllerParkingOrder) ParkingOrderList(ctx context.Context, req *parking_order.ParkingOrderListReq) (res *parking_order.ParkingOrderListRes, err error) {
	return service.ParkingOrder().ParkingOrderList(req)
}

func (c *ControllerParkingOrder) ParkingOrderUpdate(ctx context.Context, req *parking_order.ParkingOrderUpdateReq) (res *parking_order.ParkingOrderUpdateRes, err error) {
	return service.ParkingOrder().ParkingOrderUpdate(req)
}

func (c *ControllerParkingOrder) ParkingOrderDelete(ctx context.Context, req *parking_order.ParkingOrderDeleteReq) (res *parking_order.ParkingOrderDeleteRes, err error) {
	return service.ParkingOrder().ParkingOrderDelete(req)
}
