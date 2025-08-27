// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"parkin-ai-system/api/vehicle/vehicle"
)

type (
	IVehicle interface {
		Add(ctx context.Context, req *vehicle.VehicleAddReq) (res *vehicle.VehicleAddRes, err error)
		List(ctx context.Context, req *vehicle.VehicleListReq) (res *vehicle.VehicleListRes, err error)
		Detail(ctx context.Context, req *vehicle.VehicleDetailReq) (res *vehicle.VehicleDetailRes, err error)
	}
)

var (
	localVehicle IVehicle
)

func Vehicle() IVehicle {
	if localVehicle == nil {
		panic("implement not found for interface IVehicle, forgot register?")
	}
	return localVehicle
}

func RegisterVehicle(i IVehicle) {
	localVehicle = i
}
