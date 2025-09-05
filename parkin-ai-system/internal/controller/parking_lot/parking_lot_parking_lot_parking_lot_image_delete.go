package parking_lot

import (
	"context"

	"parkin-ai-system/api/parking_lot/parking_lot"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerParking_lot) ParkingLotImageDelete(ctx context.Context, req *parking_lot.ParkingLotImageDeleteReq) (res *parking_lot.ParkingLotImageDeleteRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingLotImageDeleteReq{
		Id: req.Id,
	}

	// Call service
	deleteRes, err := service.ParkingLot().ParkingLotImageDelete(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_lot.ParkingLotImageDeleteRes{
		Message: deleteRes.Message,
	}
	return res, nil
}
