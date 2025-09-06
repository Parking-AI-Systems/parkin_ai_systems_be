package parking_lot_review

import (
	"context"

	"parkin-ai-system/api/parking_lot_review/parking_lot_review"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerParking_lot_review) ParkingLotReviewAdd(ctx context.Context, req *parking_lot_review.ParkingLotReviewAddReq) (res *parking_lot_review.ParkingLotReviewAddRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingLotReviewAddReq{
		LotId:   req.LotId,
		Rating:  req.Rating,
		Comment: req.Comment,
	}

	// Call service
	addRes, err := service.ParkingLotReview().ParkingLotReviewAdd(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_lot_review.ParkingLotReviewAddRes{
		Id: addRes.Id,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
