// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package parking_slot

import (
	"context"

	"parkin-ai-system/api/parking_slot/parking_slot"
)

type IParkingSlotParking_slot interface {
	ParkingSlotAdd(ctx context.Context, req *parking_slot.ParkingSlotAddReq) (res *parking_slot.ParkingSlotAddRes, err error)
	ParkingSlotList(ctx context.Context, req *parking_slot.ParkingSlotListReq) (res *parking_slot.ParkingSlotListRes, err error)
	ParkingSlotGet(ctx context.Context, req *parking_slot.ParkingSlotGetReq) (res *parking_slot.ParkingSlotGetRes, err error)
	ParkingSlotUpdate(ctx context.Context, req *parking_slot.ParkingSlotUpdateReq) (res *parking_slot.ParkingSlotUpdateRes, err error)
	ParkingSlotDelete(ctx context.Context, req *parking_slot.ParkingSlotDeleteReq) (res *parking_slot.ParkingSlotDeleteRes, err error)
}
