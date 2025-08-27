package vehicle

import (
	"context"
	"parkin-ai-system/api/vehicle/vehicle"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerVehicle) VehicleAdd(ctx context.Context, req *vehicle.VehicleAddReq) (res *vehicle.VehicleAddRes, err error) {
	res, err = service.Vehicle().Add(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}
