package parking_lot

import (
	"context"

	"parkin-ai-system/api/parking_lot/parking_lot"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
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
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
func entityToApiParkingLotItem(item *entity.ParkingLotItem) parking_lot.ParkingLotItem {
	apiItem := parking_lot.ParkingLotItem{
		Id:             item.Id,
		Name:           item.Name,
		Address:        item.Address,
		Latitude:       item.Latitude,
		Longitude:      item.Longitude,
		OwnerId:        item.OwnerId,
		IsVerified:     item.IsVerified,
		IsActive:       item.IsActive,
		AvailableSlots: item.AvailableSlots,
		TotalSlots:     item.TotalSlots,
		PricePerHour:   item.PricePerHour,
		Description:    item.Description,
		OpenTime:       item.OpenTime,
		CloseTime:      item.CloseTime,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
		Images:         make([]parking_lot.ParkingLotImageItem, len(item.Images)),
	}
	for i, img := range item.Images {
		apiItem.Images[i] = parking_lot.ParkingLotImageItem{
			Id:           img.Id,
			ParkingLotId: img.ParkingLotId,
			LotName:      item.Name,
			ImageUrl:     img.ImageUrl,
			CreatedAt:    img.CreatedAt,
			Description:  img.Description,
			UpdatedAt:    img.UpdatedAt,
		}
	}
	return apiItem
}
