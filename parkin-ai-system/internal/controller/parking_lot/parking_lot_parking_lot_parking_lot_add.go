package parking_lot

import (
	"context"

	"parkin-ai-system/api/parking_lot/parking_lot"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

// Name         string               `json:"name"`
//
//	Address      string               `json:"address"`
//	Latitude     float64              `json:"latitude"`
//	Longitude    float64              `json:"longitude"`
//	IsVerified   bool                 `json:"isVerified"`
//	IsActive     bool                 `json:"isActive"`
//	TotalSlots   int                  `json:"totalSlots"`
//	PricePerHour float64              `json:"pricePerHour"`
//	Description  string               `json:"description"`
//	OpenTime     *gtime.Time          `json:"openTime"`
//	CloseTime    *gtime.Time          `json:"closeTime"`
//	Images       []ParkingLotImageInput `json:"images"`
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
	return res, nil
}
