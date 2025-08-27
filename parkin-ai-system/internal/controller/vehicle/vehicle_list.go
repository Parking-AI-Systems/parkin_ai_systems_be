package vehicle

import (
	"context"
	"parkin-ai-system/api/vehicle/vehicle"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerVehicle) VehicleList(ctx context.Context, req *vehicle.VehicleListReq) (res *vehicle.VehicleListRes, err error) {
	res, err = service.Vehicle().List(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return res, nil
}
