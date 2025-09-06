// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package others_service

import (
	"context"

	"parkin-ai-system/api/others_service/others_service"
)

type IOthersServiceOthers_service interface {
	OthersServiceAdd(ctx context.Context, req *others_service.OthersServiceAddReq) (res *others_service.OthersServiceAddRes, err error)
	OthersServiceList(ctx context.Context, req *others_service.OthersServiceListReq) (res *others_service.OthersServiceListRes, err error)
	OthersServiceGet(ctx context.Context, req *others_service.OthersServiceGetReq) (res *others_service.OthersServiceGetRes, err error)
	OthersServiceUpdate(ctx context.Context, req *others_service.OthersServiceUpdateReq) (res *others_service.OthersServiceUpdateRes, err error)
	OthersServiceDelete(ctx context.Context, req *others_service.OthersServiceDeleteReq) (res *others_service.OthersServiceDeleteRes, err error)
}
