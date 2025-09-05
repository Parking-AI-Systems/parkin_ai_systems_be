// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package parking_order

import (
	"context"

	"parkin-ai-system/api/parking_order/parking_order"
)

type IParkingOrderParking_order interface {
	ParkingOrderAdd(ctx context.Context, req *parking_order.ParkingOrderAddReq) (res *parking_order.ParkingOrderAddRes, err error)
	ParkingOrderList(ctx context.Context, req *parking_order.ParkingOrderListReq) (res *parking_order.ParkingOrderListRes, err error)
	ParkingOrderGet(ctx context.Context, req *parking_order.ParkingOrderGetReq) (res *parking_order.ParkingOrderGetRes, err error)
	ParkingOrderUpdate(ctx context.Context, req *parking_order.ParkingOrderUpdateReq) (res *parking_order.ParkingOrderUpdateRes, err error)
	ParkingOrderCancel(ctx context.Context, req *parking_order.ParkingOrderCancelReq) (res *parking_order.ParkingOrderCancelRes, err error)
	ParkingOrderDelete(ctx context.Context, req *parking_order.ParkingOrderDeleteReq) (res *parking_order.ParkingOrderDeleteRes, err error)
	ParkingOrderPayment(ctx context.Context, req *parking_order.ParkingOrderPaymentReq) (res *parking_order.ParkingOrderPaymentRes, err error)
}
