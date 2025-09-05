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
	IParkingLot interface {
		ParkingLotAdd(ctx context.Context, req *entity.ParkingLotAddReq) (*entity.ParkingLotAddRes, error)
		ParkingLotList(ctx context.Context, req *entity.ParkingLotListReq) (*entity.ParkingLotListRes, error)
		ParkingLotGet(ctx context.Context, req *entity.ParkingLotGetReq) (*entity.ParkingLotItem, error)
		ParkingLotUpdate(ctx context.Context, req *entity.ParkingLotUpdateReq) (*entity.ParkingLotItem, error)
		ParkingLotDelete(ctx context.Context, req *entity.ParkingLotDeleteReq) (*entity.ParkingLotDeleteRes, error)
		ParkingLotImageDelete(ctx context.Context, req *entity.ParkingLotImageDeleteReq) (*entity.ParkingLotImageDeleteRes, error)
	}
)

var (
	localParkingLot IParkingLot
)

func ParkingLot() IParkingLot {
	if localParkingLot == nil {
		panic("implement not found for interface IParkingLot, forgot register?")
	}
	return localParkingLot
}

func RegisterParkingLot(i IParkingLot) {
	localParkingLot = i
}
