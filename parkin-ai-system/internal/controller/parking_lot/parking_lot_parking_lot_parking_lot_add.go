package parking_lot

import (
	"context"

	"parkin-ai-system/api/parking_lot/parking_lot"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerParking_lot) ParkingLotAdd(ctx context.Context, req *parking_lot.ParkingLotAddReq) (res *parking_lot.ParkingLotAddRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingLotAddReq{
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
		Images:       make([]entity.ParkingLotImageInput, len(req.Images)),
	}
	for i, img := range req.Images {
		input.Images[i] = entity.ParkingLotImageInput{
			ImageUrl: img.ImageUrl,
		}
	}

	// Call service
	addRes, err := service.ParkingLot().ParkingLotAdd(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_lot.ParkingLotAddRes{
		Id: addRes.Id,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
