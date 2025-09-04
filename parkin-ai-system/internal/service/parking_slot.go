// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"parkin-ai-system/api/parking_slot"
)

type (
	IParkingSlot interface {
		ParkingSlotAdd(req *parking_slot.ParkingSlotAddReq) (*parking_slot.ParkingSlotAddRes, error)
		ParkingSlotList(req *parking_slot.ParkingSlotListReq) (*parking_slot.ParkingSlotListRes, error)
		ParkingSlotUpdate(req *parking_slot.ParkingSlotUpdateReq) (*parking_slot.ParkingSlotUpdateRes, error)
		ParkingSlotDelete(req *parking_slot.ParkingSlotDeleteReq) (*parking_slot.ParkingSlotDeleteRes, error)
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
