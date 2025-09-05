package parking_lot

import (
	"context"

	"parkin-ai-system/api/parking_lot/parking_lot"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerParking_lot) ParkingLotGet(ctx context.Context, req *parking_lot.ParkingLotGetReq) (res *parking_lot.ParkingLotGetRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingLotGetReq{
		Id: req.Id,
	}

	// Call service
	lot, err := service.ParkingLot().ParkingLotGet(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &parking_lot.ParkingLotGetRes{
		Lot: entityToApiParkingLotItem(lot),
	}
	return res, nil
}
func entityToApiParkingLotItem(item *entity.ParkingLotItem) parking_lot.ParkingLotItem {
	apiItem := parking_lot.ParkingLotItem{
		Id:           item.Id,
		Name:         item.Name,
		Address:      item.Address,
		Latitude:     item.Latitude,
		Longitude:    item.Longitude,
		TotalSlots:   item.TotalSlots,
		PricePerHour: item.PricePerHour,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
		Images:       make([]parking_lot.ParkingLotImageItem, len(item.Images)),
	}
	for i, img := range item.Images {
		apiItem.Images[i] = parking_lot.ParkingLotImageItem{
			Id:        img.Id,
			ImageUrl:  img.ImageUrl,
			CreatedAt: img.CreatedAt,
		}
	}
	return apiItem
}
