package other_service

import "context"

type IOtherService interface {
	OtherServiceAdd(ctx context.Context, req *OtherServiceAddReq) (res *OtherServiceAddRes, err error)
	OtherServiceUpdate(ctx context.Context, req *OtherServiceUpdateReq) (res *OtherServiceUpdateRes, err error)
	OtherServiceDelete(ctx context.Context, req *OtherServiceDeleteReq) (res *OtherServiceDeleteRes, err error)
	OtherServiceDetail(ctx context.Context, req *OtherServiceDetailReq) (res *OtherServiceDetailRes, err error)
	OtherServiceList(ctx context.Context, req *OtherServiceListReq) (res *OtherServiceListRes, err error)
}
