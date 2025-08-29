package parking_lot_review

import "context"

type IParkingLotReview interface {
	ParkingLotReviewAdd(ctx context.Context, req *ParkingLotReviewAddReq) (res *ParkingLotReviewAddRes, err error)
	ParkingLotReviewUpdate(ctx context.Context, req *ParkingLotReviewUpdateReq) (res *ParkingLotReviewUpdateRes, err error)
	ParkingLotReviewDelete(ctx context.Context, req *ParkingLotReviewDeleteReq) (res *ParkingLotReviewDeleteRes, err error)
	ParkingLotReviewDetail(ctx context.Context, req *ParkingLotReviewDetailReq) (res *ParkingLotReviewDetailRes, err error)
	ParkingLotReviewList(ctx context.Context, req *ParkingLotReviewListReq) (res *ParkingLotReviewListRes, err error)
}
