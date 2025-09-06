package vehicle

import (
	"context"

	"parkin-ai-system/api/vehicle/vehicle"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerVehicle) VehicleAdd(ctx context.Context, req *vehicle.VehicleAddReq) (res *vehicle.VehicleAddRes, err error) {
	// Map API request to entity request
	input := &entity.VehicleAddReq{
		LicensePlate: req.LicensePlate,
		Brand:        req.Brand,
		Model:        req.Model,
		Color:        req.Color,
		Type:         req.Type,
	}

	// Call service
	addRes, err := service.Vehicle().VehicleAdd(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &vehicle.VehicleAddRes{
		Id: addRes.Id,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
