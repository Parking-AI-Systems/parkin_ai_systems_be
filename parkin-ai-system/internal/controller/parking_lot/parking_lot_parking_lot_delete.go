package parking_lot

import (
	"context"
	"parkin-ai-system/api/parking_lot/parking_lot"
	"parkin-ai-system/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerParkingLot) ParkingLotDelete(ctx context.Context, req *parking_lot.ParkingLotDeleteReq) (res *parking_lot.ParkingLotDeleteRes, err error) {
	res, err = service.ParkingLot().ParkingLotDelete(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}
