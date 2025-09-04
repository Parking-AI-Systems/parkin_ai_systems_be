// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"parkin-ai-system/api/other_service_order"
)

type (
	IOtherServiceOrder interface {
		OtherServiceOrderAdd(ctx context.Context, req *other_service_order.OtherServiceOrderAddReq) (*other_service_order.OtherServiceOrderAddRes, error)
		OtherServiceOrderUpdate(ctx context.Context, req *other_service_order.OtherServiceOrderUpdateReq) (*other_service_order.OtherServiceOrderUpdateRes, error)
		OtherServiceOrderDelete(ctx context.Context, req *other_service_order.OtherServiceOrderDeleteReq) (*other_service_order.OtherServiceOrderDeleteRes, error)
		OtherServiceOrderList(ctx context.Context, req *other_service_order.OtherServiceOrderListReq) (*other_service_order.OtherServiceOrderListRes, error)
	}
)

var (
	localOtherServiceOrder IOtherServiceOrder
)

func OtherServiceOrder() IOtherServiceOrder {
	if localOtherServiceOrder == nil {
		panic("implement not found for interface IOtherServiceOrder, forgot register?")
	}
	return localOtherServiceOrder
}

func RegisterOtherServiceOrder(i IOtherServiceOrder) {
	localOtherServiceOrder = i
}
