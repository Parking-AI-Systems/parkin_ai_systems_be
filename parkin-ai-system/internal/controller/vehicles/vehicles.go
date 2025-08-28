package vehicles

import (
	"context"
	"parkin-ai-system/api/vehicles"
	"parkin-ai-system/internal/service"
)

type ControllerVehicles struct{}

// Tạo mới vehicle
func (c *ControllerVehicles) Create(ctx context.Context, req *vehicles.CreateVehicleReq) (res *vehicles.CreateVehicleRes, err error) {
	return service.Vehicles().Create(ctx, req)
}

// Lấy danh sách vehicles
func (c *ControllerVehicles) List(ctx context.Context, req *vehicles.ListVehicleReq) (res *vehicles.ListVehicleRes, err error) {
	return service.Vehicles().List(ctx, req)
}

// Lấy chi tiết vehicle
func (c *ControllerVehicles) Get(ctx context.Context, req *vehicles.GetVehicleReq) (res *vehicles.GetVehicleRes, err error) {
	return service.Vehicles().Get(ctx, req)
}

// Cập nhật vehicle
func (c *ControllerVehicles) Update(ctx context.Context, req *vehicles.UpdateVehicleReq) (res *vehicles.UpdateVehicleRes, err error) {
	return service.Vehicles().Update(ctx, req)
}

// Xóa vehicle
func (c *ControllerVehicles) Delete(ctx context.Context, req *vehicles.DeleteVehicleReq) (res *vehicles.DeleteVehicleRes, err error) {
	return service.Vehicles().Delete(ctx, req)
}
