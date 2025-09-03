package other_service_order

import (
	"context"
	"parkin-ai-system/api/other_service_order"
	"parkin-ai-system/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
)

type ControllerOtherServiceOrder struct{}

func NewOtherServiceOrder() *ControllerOtherServiceOrder {
	return &ControllerOtherServiceOrder{}
}

func (c *ControllerOtherServiceOrder) OtherServiceOrderAdd(ctx context.Context, req *other_service_order.OtherServiceOrderAddReq) (res *other_service_order.OtherServiceOrderAddRes, err error) {
	res, err = service.OtherServiceOrder().OtherServiceOrderAdd(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerOtherServiceOrder) OtherServiceOrderUpdate(ctx context.Context, req *other_service_order.OtherServiceOrderUpdateReq) (res *other_service_order.OtherServiceOrderUpdateRes, err error) {
	res, err = service.OtherServiceOrder().OtherServiceOrderUpdate(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerOtherServiceOrder) OtherServiceOrderDelete(ctx context.Context, req *other_service_order.OtherServiceOrderDeleteReq) (res *other_service_order.OtherServiceOrderDeleteRes, err error) {
	res, err = service.OtherServiceOrder().OtherServiceOrderDelete(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerOtherServiceOrder) OtherServiceOrderList(ctx context.Context, req *other_service_order.OtherServiceOrderListReq) (res *other_service_order.OtherServiceOrderListRes, err error) {
	res, err = service.OtherServiceOrder().OtherServiceOrderList(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}
