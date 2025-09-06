package parking_lot_review

import (
	"context"

	"parkin-ai-system/api/parking_lot_review/parking_lot_review"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerParking_lot_review) ParkingLotReviewDelete(ctx context.Context, req *parking_lot_review.ParkingLotReviewDeleteReq) (res *parking_lot_review.ParkingLotReviewDeleteRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingLotReviewDeleteReq{
		Id: req.Id,
	}

	// Call service
	deleteRes, err := service.ParkingLotReview().ParkingLotReviewDelete(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_lot_review.ParkingLotReviewDeleteRes{
		Message: deleteRes.Message,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
