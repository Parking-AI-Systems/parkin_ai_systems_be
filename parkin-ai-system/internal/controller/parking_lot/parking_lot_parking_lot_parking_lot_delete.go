package parking_lot

import (
	"context"

	"parkin-ai-system/api/parking_lot/parking_lot"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerParking_lot) ParkingLotDelete(ctx context.Context, req *parking_lot.ParkingLotDeleteReq) (res *parking_lot.ParkingLotDeleteRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingLotDeleteReq{
		Id: req.Id,
	}

	// Call service
	deleteRes, err := service.ParkingLot().ParkingLotDelete(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_lot.ParkingLotDeleteRes{
		Message: deleteRes.Message,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
