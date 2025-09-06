package vehicle

import (
	"context"

	"parkin-ai-system/api/vehicle/vehicle"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerVehicle) VehicleUpdate(ctx context.Context, req *vehicle.VehicleUpdateReq) (res *vehicle.VehicleUpdateRes, err error) {
	// Map API request to entity request
	input := &entity.VehicleUpdateReq{
		Id:           req.Id,
		LicensePlate: req.LicensePlate,
		Brand:        req.Brand,
		Model:        req.Model,
		Color:        req.Color,
		Type:         req.Type,
	}

	// Call service
	updateRes, err := service.Vehicle().VehicleUpdate(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &vehicle.VehicleUpdateRes{
		Vehicle: vehicle.VehicleItem{
			Id:           updateRes.Id,
			LicensePlate: updateRes.LicensePlate,
			Brand:        updateRes.Brand,
			Model:        updateRes.Model,
			Color:        updateRes.Color,
			Type:         updateRes.Type,
			CreatedAt:    updateRes.CreatedAt,
		},
	}
	return res, nil
}
