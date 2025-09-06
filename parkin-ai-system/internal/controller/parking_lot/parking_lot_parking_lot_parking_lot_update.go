package parking_lot

import (
	"context"

	"parkin-ai-system/api/parking_lot/parking_lot"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerParking_lot) ParkingLotUpdate(ctx context.Context, req *parking_lot.ParkingLotUpdateReq) (res *parking_lot.ParkingLotUpdateRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingLotUpdateReq{
		Id:           req.Id,
		Name:         req.Name,
		Address:      req.Address,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		IsVerified:   req.IsVerified,
		IsActive:     req.IsActive,
		TotalSlots:   req.TotalSlots,
		PricePerHour: req.PricePerHour,
		Description:  req.Description,
		OpenTime:     req.OpenTime,
		CloseTime:    req.CloseTime,
	}

	// Call service
	updateRes, err := service.ParkingLot().ParkingLotUpdate(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_lot.ParkingLotUpdateRes{
		Lot: entityToApiParkingLotItem(updateRes),
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
