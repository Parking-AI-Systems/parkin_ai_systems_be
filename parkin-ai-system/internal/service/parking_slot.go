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
	IParkingSlot interface {
		ParkingSlotAdd(ctx context.Context, req *entity.ParkingSlotAddReq) (*entity.ParkingSlotAddRes, error)
		ParkingSlotList(ctx context.Context, req *entity.ParkingSlotListReq) (*entity.ParkingSlotListRes, error)
		ParkingSlotGet(ctx context.Context, req *entity.ParkingSlotGetReq) (*entity.ParkingSlotItem, error)
		ParkingSlotUpdate(ctx context.Context, req *entity.ParkingSlotUpdateReq) (*entity.ParkingSlotItem, error)
		ParkingSlotDelete(ctx context.Context, req *entity.ParkingSlotDeleteReq) (*entity.ParkingSlotDeleteRes, error)
	}
)

var (
	localParkingSlot IParkingSlot
)

func ParkingSlot() IParkingSlot {
	if localParkingSlot == nil {
		panic("implement not found for interface IParkingSlot, forgot register?")
	}
	return localParkingSlot
}

func RegisterParkingSlot(i IParkingSlot) {
	localParkingSlot = i
}
