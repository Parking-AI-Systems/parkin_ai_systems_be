// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"parkin-ai-system/internal/model/entity"
)

type (
	IVehicle interface {
		VehicleAdd(ctx context.Context, req *entity.VehicleAddReq) (*entity.VehicleAddRes, error)
		VehicleList(ctx context.Context, req *entity.VehicleListReq) (*entity.VehicleListRes, error)
		VehicleGet(ctx context.Context, req *entity.VehicleGetReq) (*entity.VehicleItem, error)
		VehicleUpdate(ctx context.Context, req *entity.VehicleUpdateReq) (*entity.VehicleItem, error)
		VehicleDelete(ctx context.Context, req *entity.VehicleDeleteReq) (*entity.VehicleDeleteRes, error)
		// CheckVehicleSlotCompatibility checks if a vehicle type is compatible with a parking slot type
		CheckVehicleSlotCompatibility(ctx context.Context, vehicleID int64, slotID int64) error
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
