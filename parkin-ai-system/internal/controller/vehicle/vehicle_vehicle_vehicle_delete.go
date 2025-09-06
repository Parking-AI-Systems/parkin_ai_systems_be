package vehicle

import (
	"context"

	"parkin-ai-system/api/vehicle/vehicle"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerVehicle) VehicleDelete(ctx context.Context, req *vehicle.VehicleDeleteReq) (res *vehicle.VehicleDeleteRes, err error) {
	// Map API request to entity request
	input := &entity.VehicleDeleteReq{
		Id: req.Id,
	}

	// Call service
	deleteRes, err := service.Vehicle().VehicleDelete(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &vehicle.VehicleDeleteRes{
		Message: deleteRes.Message,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
