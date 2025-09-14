package parking_lot_review

import (
	"context"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"parkin-ai-system/api/parking_lot_review/parking_lot_review"
)

func (c *ControllerParking_lot_review) MyParkingLotReviewList(ctx context.Context, req *parking_lot_review.MyParkingLotReviewListReq) (res *parking_lot_review.MyParkingLotReviewListRes, err error) {
	return nil, gerror.NewCode(gcode.CodeNotImplemented)
}
