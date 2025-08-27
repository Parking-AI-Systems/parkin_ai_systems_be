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
	VehicleDetail(ctx context.Context, req *vehicle.VehicleDetailReq) (res *vehicle.VehicleDetailRes, err error)
	VehicleList(ctx context.Context, req *vehicle.VehicleListReq) (res *vehicle.VehicleListRes, err error)
	VehicleUpdate(ctx context.Context, req *vehicle.VehicleUpdateReq) (res *vehicle.VehicleUpdateRes, err error)
}
