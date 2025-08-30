package service

import (
	"context"
	"parkin-ai-system/api/other_service"
)

type IOtherService interface {
	OtherServiceAdd(ctx context.Context, req *other_service.OtherServiceAddReq) (res *other_service.OtherServiceAddRes, err error)
	OtherServiceUpdate(ctx context.Context, req *other_service.OtherServiceUpdateReq) (res *other_service.OtherServiceUpdateRes, err error)
	OtherServiceDelete(ctx context.Context, req *other_service.OtherServiceDeleteReq) (res *other_service.OtherServiceDeleteRes, err error)
	OtherServiceDetail(ctx context.Context, req *other_service.OtherServiceDetailReq) (res *other_service.OtherServiceDetailRes, err error)
	OtherServiceList(ctx context.Context, req *other_service.OtherServiceListReq) (res *other_service.OtherServiceListRes, err error)
}

var localOtherService IOtherService

func OtherService() IOtherService {
	if localOtherService == nil {
		panic("implement not found for interface IOtherService, forgot register?")
	}
	return localOtherService
}

func RegisterOtherService(i IOtherService) {
	localOtherService = i
}
