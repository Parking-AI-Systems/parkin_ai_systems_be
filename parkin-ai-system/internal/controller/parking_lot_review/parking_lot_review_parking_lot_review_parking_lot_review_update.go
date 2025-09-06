package parking_lot_review

import (
	"context"

	"parkin-ai-system/api/parking_lot_review/parking_lot_review"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerParking_lot_review) ParkingLotReviewUpdate(ctx context.Context, req *parking_lot_review.ParkingLotReviewUpdateReq) (res *parking_lot_review.ParkingLotReviewUpdateRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingLotReviewUpdateReq{
		Id:      req.Id,
		Rating:  req.Rating,
		Comment: req.Comment,
	}

	// Call service
	updateRes, err := service.ParkingLotReview().ParkingLotReviewUpdate(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_lot_review.ParkingLotReviewUpdateRes{
		Review: parking_lot_review.ParkingLotReviewItem{
			Id:        updateRes.Id,
			LotId:     updateRes.LotId,
			UserId:    updateRes.UserId,
			Rating:    updateRes.Rating,
			Comment:   updateRes.Comment,
			CreatedAt: updateRes.CreatedAt,
		},
	}
	return res, nil
}
