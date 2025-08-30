package other_service

import (
	"context"
	"parkin-ai-system/api/other_service"
	"parkin-ai-system/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerOtherService) OtherServiceAdd(ctx context.Context, req *other_service.OtherServiceAddReq) (res *other_service.OtherServiceAddRes, err error) {
	res, err = service.OtherService().OtherServiceAdd(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerOtherService) OtherServiceUpdate(ctx context.Context, req *other_service.OtherServiceUpdateReq) (res *other_service.OtherServiceUpdateRes, err error) {
	res, err = service.OtherService().OtherServiceUpdate(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerOtherService) OtherServiceDelete(ctx context.Context, req *other_service.OtherServiceDeleteReq) (res *other_service.OtherServiceDeleteRes, err error) {
	res, err = service.OtherService().OtherServiceDelete(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerOtherService) OtherServiceDetail(ctx context.Context, req *other_service.OtherServiceDetailReq) (res *other_service.OtherServiceDetailRes, err error) {
	res, err = service.OtherService().OtherServiceDetail(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerOtherService) OtherServiceList(ctx context.Context, req *other_service.OtherServiceListReq) (res *other_service.OtherServiceListRes, err error) {
	res, err = service.OtherService().OtherServiceList(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}
