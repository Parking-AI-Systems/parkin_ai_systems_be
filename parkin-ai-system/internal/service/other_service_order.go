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
	IOthersServiceOrder interface {
		OthersServiceOrderAddWithUser(ctx context.Context, req *entity.OthersServiceOrderAddReq) (*entity.OthersServiceOrderAddRes, error)
		OthersServiceOrderList(ctx context.Context, req *entity.OthersServiceOrderListReq) (*entity.OthersServiceOrderListRes, error)
		OthersServiceOrderGet(ctx context.Context, req *entity.OthersServiceOrderGetReq) (*entity.OthersServiceOrderItem, error)
		OthersServiceOrderUpdate(ctx context.Context, req *entity.OthersServiceOrderUpdateReq) (*entity.OthersServiceOrderItem, error)
		OthersServiceOrderCancel(ctx context.Context, req *entity.OthersServiceOrderCancelReq) (*entity.OthersServiceOrderItem, error)
		OthersServiceOrderDelete(ctx context.Context, req *entity.OthersServiceOrderDeleteReq) (*entity.OthersServiceOrderDeleteRes, error)
		OthersServiceOrderPayment(ctx context.Context, req *entity.OthersServiceOrderPaymentReq) (*entity.OthersServiceOrderItem, error)
		OthersServiceRevenue(ctx context.Context, req *entity.OthersServiceRevenueReq) (*entity.OthersServiceRevenueRes, error)
		OthersServicePopular(ctx context.Context, req *entity.OthersServicePopularReq) (*entity.OthersServicePopularRes, error)
		OthersServiceTrends(ctx context.Context, req *entity.OthersServiceTrendsReq) (*entity.OthersServiceTrendsRes, error)
	}
)

var (
	localOthersServiceOrder IOthersServiceOrder
)

func OthersServiceOrder() IOthersServiceOrder {
	if localOthersServiceOrder == nil {
		panic("implement not found for interface IOthersServiceOrder, forgot register?")
	}
	return localOthersServiceOrder
}

func RegisterOthersServiceOrder(i IOthersServiceOrder) {
	localOthersServiceOrder = i
}
