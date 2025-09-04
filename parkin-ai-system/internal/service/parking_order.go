// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"parkin-ai-system/api/parking_order"
)

type (
	IParkingOrder interface {
		ParkingOrderAdd(req *parking_order.ParkingOrderAddReq) (*parking_order.ParkingOrderAddRes, error)
		ParkingOrderAddWithUser(ctx context.Context, req *parking_order.ParkingOrderAddReq) (*parking_order.ParkingOrderAddRes, error)
		ParkingOrderList(req *parking_order.ParkingOrderListReq) (*parking_order.ParkingOrderListRes, error)
		ParkingOrderUpdate(req *parking_order.ParkingOrderUpdateReq) (*parking_order.ParkingOrderUpdateRes, error)
		ParkingOrderDelete(req *parking_order.ParkingOrderDeleteReq) (*parking_order.ParkingOrderDeleteRes, error)
	}
)

var (
	localParkingOrder IParkingOrder
)

func ParkingOrder() IParkingOrder {
	if localParkingOrder == nil {
		panic("implement not found for interface IParkingOrder, forgot register?")
	}
	return localParkingOrder
}

func RegisterParkingOrder(i IParkingOrder) {
	localParkingOrder = i
}
