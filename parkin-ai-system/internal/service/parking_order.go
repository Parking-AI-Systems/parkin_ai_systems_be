// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"parkin-ai-system/internal/model/entity"
)

type (
	IParkingOrder interface {
		ParkingOrderAddWithUser(ctx context.Context, req *entity.ParkingOrderAddReq) (*entity.ParkingOrderAddRes, error)
		ParkingOrderList(ctx context.Context, req *entity.ParkingOrderListReq) (*entity.ParkingOrderListRes, error)
		ParkingOrderGet(ctx context.Context, req *entity.ParkingOrderGetReq) (*entity.ParkingOrderItem, error)
		ParkingOrderUpdate(ctx context.Context, req *entity.ParkingOrderUpdateReq) (*entity.ParkingOrderItem, error)
		ParkingOrderCancel(ctx context.Context, req *entity.ParkingOrderCancelReq) (*entity.ParkingOrderItem, error)
		ParkingOrderDelete(ctx context.Context, req *entity.ParkingOrderDeleteReq) (*entity.ParkingOrderDeleteRes, error)
		ParkingOrderPayment(ctx context.Context, req *entity.ParkingOrderPaymentReq) (*entity.ParkingOrderItem, error)
		ParkingOrderRevenue(ctx context.Context, req *entity.ParkingOrderRevenueReq) (*entity.ParkingOrderRevenueRes, error)
		ParkingOrderTrends(ctx context.Context, req *entity.ParkingOrderTrendsReq) (*entity.ParkingOrderTrendsRes, error)
		ParkingOrderStatusBreakdown(ctx context.Context, req *entity.ParkingOrderStatusBreakdownReq) (*entity.ParkingOrderStatusBreakdownRes, error)
		GetMyLotOrder(ctx context.Context, req *entity.GetMyLotOrderReq) (*entity.GetMyLotOrderRes, error)
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
