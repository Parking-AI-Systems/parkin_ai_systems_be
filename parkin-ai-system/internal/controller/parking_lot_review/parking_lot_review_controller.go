package parking_lot_review

import (
	"context"
	"parkin-ai-system/api/parking_lot_review"
	"parkin-ai-system/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerParkingLotReview) ParkingLotReviewAdd(ctx context.Context, req *parking_lot_review.ParkingLotReviewAddReq) (res *parking_lot_review.ParkingLotReviewAddRes, err error) {
	res, err = service.ParkingLotReview().ParkingLotReviewAdd(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerParkingLotReview) ParkingLotReviewUpdate(ctx context.Context, req *parking_lot_review.ParkingLotReviewUpdateReq) (res *parking_lot_review.ParkingLotReviewUpdateRes, err error) {
	res, err = service.ParkingLotReview().ParkingLotReviewUpdate(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerParkingLotReview) ParkingLotReviewDelete(ctx context.Context, req *parking_lot_review.ParkingLotReviewDeleteReq) (res *parking_lot_review.ParkingLotReviewDeleteRes, err error) {
	res, err = service.ParkingLotReview().ParkingLotReviewDelete(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerParkingLotReview) ParkingLotReviewDetail(ctx context.Context, req *parking_lot_review.ParkingLotReviewDetailReq) (res *parking_lot_review.ParkingLotReviewDetailRes, err error) {
	res, err = service.ParkingLotReview().ParkingLotReviewDetail(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerParkingLotReview) ParkingLotReviewList(ctx context.Context, req *parking_lot_review.ParkingLotReviewListReq) (res *parking_lot_review.ParkingLotReviewListRes, err error) {
	res, err = service.ParkingLotReview().ParkingLotReviewList(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}
