package service

import (
	"context"
	"parkin-ai-system/api/vehicles"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/do"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type sVehicles struct{}

var (
	vehiclesService = sVehicles{}
)

func Vehicles() *sVehicles {
	return &vehiclesService
}

// Create vehicle
func (s *sVehicles) Create(ctx context.Context, req *vehicles.CreateVehicleReq) (*vehicles.CreateVehicleRes, error) {
	data := do.Vehicles{}
	gconv.Struct(req, &data)
	id, err := dao.Vehicles.Ctx(ctx).Data(data).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	return &vehicles.CreateVehicleRes{Id: gconv.Int64(id)}, nil
}

// List vehicles
func (s *sVehicles) List(ctx context.Context, req *vehicles.ListVehicleReq) (*vehicles.ListVehicleRes, error) {
	var items []vehicles.VehicleItem
	err := dao.Vehicles.Ctx(ctx).Where("user_id", req.UserId).Scan(&items)
	if err != nil {
		return nil, err
	}
	return &vehicles.ListVehicleRes{Vehicles: items}, nil
}

// Get vehicle detail
func (s *sVehicles) Get(ctx context.Context, req *vehicles.GetVehicleReq) (*vehicles.GetVehicleRes, error) {
	var item vehicles.VehicleItem
	err := dao.Vehicles.Ctx(ctx).Where("id", req.Id).Scan(&item)
	if err != nil {
		return nil, err
	}
	return &vehicles.GetVehicleRes{Vehicle: item}, nil
}

// Update vehicle
func (s *sVehicles) Update(ctx context.Context, req *vehicles.UpdateVehicleReq) (*vehicles.UpdateVehicleRes, error) {
	data := do.Vehicles{}
	gconv.Struct(req, &data)
	_, err := dao.Vehicles.Ctx(ctx).Where("id", req.Id).Data(data).Update()
	if err != nil {
		return nil, err
	}
	return &vehicles.UpdateVehicleRes{Success: true}, nil
}

// Delete vehicle
func (s *sVehicles) Delete(ctx context.Context, req *vehicles.DeleteVehicleReq) (*vehicles.DeleteVehicleRes, error) {
	_, err := dao.Vehicles.Ctx(ctx).Where("id", req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &vehicles.DeleteVehicleRes{Success: true}, nil
}
