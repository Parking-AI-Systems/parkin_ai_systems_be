package parking_lot_review

import (
	"context"

	"parkin-ai-system/api/parking_lot_review/parking_lot_review"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerParking_lot_review) ParkingLotReviewList(ctx context.Context, req *parking_lot_review.ParkingLotReviewListReq) (res *parking_lot_review.ParkingLotReviewListRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingLotReviewListReq{
		LotId:    req.LotId,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// Call service
	listRes, err := service.ParkingLotReview().ParkingLotReviewList(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity list to API response
	res = &parking_lot_review.ParkingLotReviewListRes{
		List:  make([]parking_lot_review.ParkingLotReviewItem, 0, len(listRes.List)),
		Total: listRes.Total,
	}
	for _, item := range listRes.List {
		res.List = append(res.List, parking_lot_review.ParkingLotReviewItem{
			Id:        item.Id,
			LotId:     item.LotId,
			UserId:    item.UserId,
			Rating:    item.Rating,
			Comment:   item.Comment,
			CreatedAt: item.CreatedAt,
		})
	}
	return res, nil
}
