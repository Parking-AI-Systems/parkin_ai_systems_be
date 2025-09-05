// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package other_service_orders

import (
	"context"

	"parkin-ai-system/api/other_service_orders/other_service_orders"
)

type IOtherServiceOrdersOther_service_orders interface {
	OthersServiceOrderAdd(ctx context.Context, req *other_service_orders.OthersServiceOrderAddReq) (res *other_service_orders.OthersServiceOrderAddRes, err error)
	OthersServiceOrderList(ctx context.Context, req *other_service_orders.OthersServiceOrderListReq) (res *other_service_orders.OthersServiceOrderListRes, err error)
	OthersServiceOrderGet(ctx context.Context, req *other_service_orders.OthersServiceOrderGetReq) (res *other_service_orders.OthersServiceOrderGetRes, err error)
	OthersServiceOrderUpdate(ctx context.Context, req *other_service_orders.OthersServiceOrderUpdateReq) (res *other_service_orders.OthersServiceOrderUpdateRes, err error)
	OthersServiceOrderCancel(ctx context.Context, req *other_service_orders.OthersServiceOrderCancelReq) (res *other_service_orders.OthersServiceOrderCancelRes, err error)
	OthersServiceOrderDelete(ctx context.Context, req *other_service_orders.OthersServiceOrderDeleteReq) (res *other_service_orders.OthersServiceOrderDeleteRes, err error)
	OthersServiceOrderPayment(ctx context.Context, req *other_service_orders.OthersServiceOrderPaymentReq) (res *other_service_orders.OthersServiceOrderPaymentRes, err error)
}
