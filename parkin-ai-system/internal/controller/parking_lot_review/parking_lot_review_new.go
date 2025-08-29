package parking_lot_review

import (
	"parkin-ai-system/api/parking_lot_review"
)

type ControllerParkingLotReview struct{}

func NewParkingLotReview() parking_lot_review.IParkingLotReview {
	return &ControllerParkingLotReview{}
}
