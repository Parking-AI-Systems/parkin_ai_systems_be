package other_service_order

import "context"

type IOtherServiceOrder interface {
	OtherServiceOrderAdd(ctx context.Context, req *OtherServiceOrderAddReq) (*OtherServiceOrderAddRes, error)
	OtherServiceOrderUpdate(ctx context.Context, req *OtherServiceOrderUpdateReq) (*OtherServiceOrderUpdateRes, error)
	OtherServiceOrderDelete(ctx context.Context, req *OtherServiceOrderDeleteReq) (*OtherServiceOrderDeleteRes, error)
	OtherServiceOrderList(ctx context.Context, req *OtherServiceOrderListReq) (*OtherServiceOrderListRes, error)
}
