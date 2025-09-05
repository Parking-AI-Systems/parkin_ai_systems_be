package parking_lot

import (
	"context"

	"parkin-ai-system/api/parking_lot/parking_lot"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerParking_lot) ParkingLotAdd(ctx context.Context, req *parking_lot.ParkingLotAddReq) (res *parking_lot.ParkingLotAddRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingLotAddReq{
		Name:         req.Name,
		Address:      req.Address,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		TotalSlots:   req.TotalSlots,
		PricePerHour: req.PricePerHour,
		Images:       make([]entity.ParkingLotImageInput, len(req.Images)),
	}
	for i, img := range req.Images {
		input.Images[i] = entity.ParkingLotImageInput{
			ImageUrl:    img.ImageUrl,
			Description: img.Description,
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
	return res, nil
}
