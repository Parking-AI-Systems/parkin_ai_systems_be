// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package parking_lot_review

import (
	"context"

	"parkin-ai-system/api/parking_lot_review/parking_lot_review"
)

type IParkingLotReviewParking_lot_review interface {
	ParkingLotReviewAdd(ctx context.Context, req *parking_lot_review.ParkingLotReviewAddReq) (res *parking_lot_review.ParkingLotReviewAddRes, err error)
	ParkingLotReviewList(ctx context.Context, req *parking_lot_review.ParkingLotReviewListReq) (res *parking_lot_review.ParkingLotReviewListRes, err error)
	ParkingLotReviewGet(ctx context.Context, req *parking_lot_review.ParkingLotReviewGetReq) (res *parking_lot_review.ParkingLotReviewGetRes, err error)
	ParkingLotReviewUpdate(ctx context.Context, req *parking_lot_review.ParkingLotReviewUpdateReq) (res *parking_lot_review.ParkingLotReviewUpdateRes, err error)
	ParkingLotReviewDelete(ctx context.Context, req *parking_lot_review.ParkingLotReviewDeleteReq) (res *parking_lot_review.ParkingLotReviewDeleteRes, err error)
	MyParkingLotReviewList(ctx context.Context, req *parking_lot_review.MyParkingLotReviewListReq) (res *parking_lot_review.MyParkingLotReviewListRes, err error)
}
