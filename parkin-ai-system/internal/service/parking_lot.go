// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================


package service

import (
	"context"
	"parkin-ai-system/api/parking_lot/parking_lot"
)

type IParkingLot interface {
	ParkingLotAdd(ctx context.Context, req *parking_lot.ParkingLotAddReq) (res *parking_lot.ParkingLotAddRes, err error)
	ParkingLotDetail(ctx context.Context, req *parking_lot.ParkingLotDetailReq) (res *parking_lot.ParkingLotDetailRes, err error)
	ParkingLotList(ctx context.Context, req *parking_lot.ParkingLotListReq) (res *parking_lot.ParkingLotListRes, err error)
	ParkingLotUpdate(ctx context.Context, req *parking_lot.ParkingLotUpdateReq) (res *parking_lot.ParkingLotUpdateRes, err error)
	ParkingLotDelete(ctx context.Context, req *parking_lot.ParkingLotDeleteReq) (res *parking_lot.ParkingLotDeleteRes, err error)
}

var localParkingLot IParkingLot

func ParkingLot() IParkingLot {
	if localParkingLot == nil {
		panic("implement not found for interface IParkingLot, forgot register?")
	}
	return localParkingLot
}

func RegisterParkingLot(i IParkingLot) {
	localParkingLot = i
}
