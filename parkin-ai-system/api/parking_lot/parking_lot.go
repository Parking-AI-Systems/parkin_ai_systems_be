// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================
package parking_lot

import (
	"context"
	"parkin-ai-system/api/parking_lot/parking_lot"
)

type IParkingLotParkingLot interface {
	ParkingLotAdd(ctx context.Context, req *parking_lot.ParkingLotAddReq) (res *parking_lot.ParkingLotAddRes, err error)
	ParkingLotDetail(ctx context.Context, req *parking_lot.ParkingLotDetailReq) (res *parking_lot.ParkingLotDetailRes, err error)
}
