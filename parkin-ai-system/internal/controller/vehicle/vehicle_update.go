package vehicle

import (
	"context"
	"parkin-ai-system/api/vehicle/vehicle"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerVehicle) VehicleUpdate(ctx context.Context, req *vehicle.VehicleUpdateReq) (res *vehicle.VehicleUpdateRes, err error) {
	res, err = service.Vehicle().Update(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return res, nil
}
