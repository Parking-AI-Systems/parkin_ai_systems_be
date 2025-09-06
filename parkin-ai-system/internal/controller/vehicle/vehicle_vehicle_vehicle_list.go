package vehicle

import (
	"context"

	"parkin-ai-system/api/vehicle/vehicle"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerVehicle) VehicleList(ctx context.Context, req *vehicle.VehicleListReq) (res *vehicle.VehicleListRes, err error) {
	// Map API request to entity request
	input := &entity.VehicleListReq{
		Type:     req.Type,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	// Call service
	listRes, err := service.Vehicle().VehicleList(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity list to API response
	res = &vehicle.VehicleListRes{
		List:  make([]vehicle.VehicleItem, 0, len(listRes.List)),
		Total: listRes.Total,
	}
	for _, item := range listRes.List {
		res.List = append(res.List, vehicle.VehicleItem{
			Id:           item.Id,
			UserId:       item.UserId,
			Username:     item.Username,
			LicensePlate: item.LicensePlate,
			Brand:        item.Brand,
			Model:        item.Model,
			Color:        item.Color,
			Type:         item.Type,
			CreatedAt:    item.CreatedAt,
		})
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
