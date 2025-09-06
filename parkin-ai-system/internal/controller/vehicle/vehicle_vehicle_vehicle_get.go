package vehicle

import (
	"context"

	"parkin-ai-system/api/vehicle/vehicle"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerVehicle) VehicleGet(ctx context.Context, req *vehicle.VehicleGetReq) (res *vehicle.VehicleGetRes, err error) {
	// Map API request to entity request
	input := &entity.VehicleGetReq{
		Id: req.Id,
	}

	// Call service
	veh, err := service.Vehicle().VehicleGet(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &vehicle.VehicleGetRes{
		Vehicle: vehicle.VehicleItem{
			Id:           veh.Id,
			LicensePlate: veh.LicensePlate,
			Brand:        veh.Brand,
			Model:        veh.Model,
			Color:        veh.Color,
			Type:         veh.Type,
			CreatedAt:    veh.CreatedAt,
		},
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
