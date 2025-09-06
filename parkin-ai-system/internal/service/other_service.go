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
	IOthersService interface {
		OthersServiceAdd(ctx context.Context, req *entity.OthersServiceAddReq) (*entity.OthersServiceAddRes, error)
		OthersServiceList(ctx context.Context, req *entity.OthersServiceListReq) (*entity.OthersServiceListRes, error)
		OthersServiceGet(ctx context.Context, req *entity.OthersServiceGetReq) (*entity.OthersServiceItem, error)
		OthersServiceUpdate(ctx context.Context, req *entity.OthersServiceUpdateReq) (*entity.OthersServiceItem, error)
		OthersServiceDelete(ctx context.Context, req *entity.OthersServiceDeleteReq) (*entity.OthersServiceDeleteRes, error)
	}
)

var (
	localOthersService IOthersService
)

func OthersService() IOthersService {
	if localOthersService == nil {
		panic("implement not found for interface IOthersService, forgot register?")
	}
	return localOthersService
}

func RegisterOthersService(i IOthersService) {
	localOthersService = i
}
