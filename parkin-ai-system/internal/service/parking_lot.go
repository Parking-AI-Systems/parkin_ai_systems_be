// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	api_add "parkin-ai-system/api/parking_lot/parking_lot"
	api_delete "parkin-ai-system/api/parking_lot/parking_lot"
	api_detail "parkin-ai-system/api/parking_lot/parking_lot"
	api_list "parkin-ai-system/api/parking_lot/parking_lot"
	api_update "parkin-ai-system/api/parking_lot/parking_lot"
)

type (
	IParkingLot interface {
		ParkingLotList(ctx context.Context, req *api_list.ParkingLotListReq) (res *api_list.ParkingLotListRes, err error)
		ParkingLotUpdate(ctx context.Context, req *api_update.ParkingLotUpdateReq) (res *api_update.ParkingLotUpdateRes, err error)
		ParkingLotDelete(ctx context.Context, req *api_delete.ParkingLotDeleteReq) (res *api_delete.ParkingLotDeleteRes, err error)
		ParkingLotAdd(ctx context.Context, req *api_add.ParkingLotAddReq) (res *api_add.ParkingLotAddRes, err error)
		ParkingLotDetail(ctx context.Context, req *api_detail.ParkingLotDetailReq) (res *api_detail.ParkingLotDetailRes, err error)
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
