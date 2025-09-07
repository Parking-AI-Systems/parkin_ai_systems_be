package parking_lot_review

import (
	"context"

	"parkin-ai-system/api/parking_lot_review/parking_lot_review"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerParking_lot_review) ParkingLotReviewGet(ctx context.Context, req *parking_lot_review.ParkingLotReviewGetReq) (res *parking_lot_review.ParkingLotReviewGetRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingLotReviewGetReq{
		Id: req.Id,
	}

	// Call service
	review, err := service.ParkingLotReview().ParkingLotReviewGet(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_lot_review.ParkingLotReviewGetRes{
		Review: parking_lot_review.ParkingLotReviewItem{
			Id:        review.Id,
			LotId:     review.LotId,
			UserId:    review.UserId,
			LotName:   review.LotName,
			Username:  review.Username,
			Rating:    review.Rating,
			Comment:   review.Comment,
			CreatedAt: review.CreatedAt,
		},
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
