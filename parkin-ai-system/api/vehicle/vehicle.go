// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package vehicle

import (
	"context"

	"parkin-ai-system/api/vehicle/vehicle"
)

type IVehicleVehicle interface {
	VehicleAdd(ctx context.Context, req *vehicle.VehicleAddReq) (res *vehicle.VehicleAddRes, err error)
}
