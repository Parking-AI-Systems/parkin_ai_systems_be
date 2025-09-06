package parking_lot

import (
	"context"

	"parkin-ai-system/api/parking_lot/parking_lot"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerParking_lot) ParkingLotList(ctx context.Context, req *parking_lot.ParkingLotListReq) (res *parking_lot.ParkingLotListRes, err error) {
	// Map API request to entity request
	input := &entity.ParkingLotListReq{
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Radius:    req.Radius,
		Page:      req.Page,
		PageSize:  req.PageSize,
	}

	// Call service
	listRes, err := service.ParkingLot().ParkingLotList(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity list to API response
	res = &parking_lot.ParkingLotListRes{
		List:  make([]parking_lot.ParkingLotItem, 0, len(listRes.List)),
		Total: listRes.Total,
	}
	for _, item := range listRes.List {
		res.List = append(res.List, entityToApiParkingLotItem(&item))
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
