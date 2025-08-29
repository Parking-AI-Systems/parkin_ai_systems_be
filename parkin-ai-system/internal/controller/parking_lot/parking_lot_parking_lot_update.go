package parking_lot

import (
	"context"
	"parkin-ai-system/api/parking_lot/parking_lot"
	"parkin-ai-system/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerParkingLot) ParkingLotUpdate(ctx context.Context, req *parking_lot.ParkingLotUpdateReq) (res *parking_lot.ParkingLotUpdateRes, err error) {
	res, err = service.ParkingLot().ParkingLotUpdate(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}
