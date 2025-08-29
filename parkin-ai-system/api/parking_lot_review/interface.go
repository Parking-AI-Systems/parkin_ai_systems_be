package parking_lot_review

import (
	"context"
	"parkin-ai-system/api/parking_lot_review"
)

type IParkingLotReview interface {
	ParkingLotReviewAdd(ctx context.Context, req *parking_lot_review.ParkingLotReviewAddReq) (res *parking_lot_review.ParkingLotReviewAddRes, err error)
	ParkingLotReviewUpdate(ctx context.Context, req *parking_lot_review.ParkingLotReviewUpdateReq) (res *parking_lot_review.ParkingLotReviewUpdateRes, err error)
	ParkingLotReviewDelete(ctx context.Context, req *parking_lot_review.ParkingLotReviewDeleteReq) (res *parking_lot_review.ParkingLotReviewDeleteRes, err error)
	ParkingLotReviewDetail(ctx context.Context, req *parking_lot_review.ParkingLotReviewDetailReq) (res *parking_lot_review.ParkingLotReviewDetailRes, err error)
	ParkingLotReviewList(ctx context.Context, req *parking_lot_review.ParkingLotReviewListReq) (res *parking_lot_review.ParkingLotReviewListRes, err error)
}
