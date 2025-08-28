package vehicle

import (
	"context"
	"fmt"
	"parkin-ai-system/api/vehicle/vehicle"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type sVehicle struct{}

func Init() {
	service.RegisterVehicle(&sVehicle{})
}

func init() {
	Init()
}

func (s *sVehicle) Add(ctx context.Context, req *vehicle.VehicleAddReq) (res *vehicle.VehicleAddRes, err error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.New("Unauthorized")
	}

	count, err := dao.Vehicles.Ctx(ctx).Where("user_id", userID).Count()
	if err != nil {
		return nil, gerror.New("Database error")
	}
	if count >= 5 {
		return nil, gerror.New("Max vehicles reached (5)")
	}

	exists, err := dao.Vehicles.Ctx(ctx).Where("license_plate", req.LicensePlate).Count()
	if err != nil {
		return nil, gerror.New("Database error")
	}
	if exists > 0 {
		return nil, gerror.New("License plate already exists")
	}

	result, err := dao.Vehicles.Ctx(ctx).Data(g.Map{
		"user_id":       userID,
		"license_plate": req.LicensePlate,
		"model":         req.Model,
		"color":         req.Color,
		"brand":         req.Brand,
		"type":          req.Type,
		"created_at":    gtime.Now(),
	}).Insert()
	if err != nil {
		return nil, gerror.New("Failed to add vehicle")
	}

	vehicleID, err := result.LastInsertId()
	if err != nil {
		return nil, gerror.New("Failed to get vehicle id")
	}

	res = &vehicle.VehicleAddRes{
		VehicleID: fmt.Sprintf("%d", vehicleID),
	}
	return
}

func (s *sVehicle) List(ctx context.Context, req *vehicle.VehicleListReq) (res *vehicle.VehicleListRes, err error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.New("Unauthorized")
	}

	records, err := dao.Vehicles.Ctx(ctx).
		Where("user_id", userID).
		Order("created_at DESC").
		All()
	if err != nil {
		return nil, gerror.New("Database error")
	}

	vehicles := make([]vehicle.VehicleListItem, 0, len(records))
	for _, v := range records {
		vehicles = append(vehicles, vehicle.VehicleListItem{
			ID:           v["id"].String(),
			LicensePlate: v["license_plate"].String(),
			Model:        v["model"].String(),
			Color:        v["color"].String(),
			Brand:        v["brand"].String(),
			Type:         v["type"].String(),
			CreatedAt:    v["created_at"].String(),
		})
	}

	res = &vehicle.VehicleListRes{Vehicles: vehicles}
	return
}

func (s *sVehicle) Detail(ctx context.Context, req *vehicle.VehicleDetailReq) (res *vehicle.VehicleDetailRes, err error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.New("Unauthorized")
	}

	v, err := dao.Vehicles.Ctx(ctx).Where("id", req.ID).One()
	if err != nil {
		return nil, gerror.New("Database error")
	}
	if v.IsEmpty() {
		return nil, gerror.New("Vehicle not found")
	}
	if v["user_id"].String() != userID {
		return nil, gerror.New("Not owner")
	}

	res = &vehicle.VehicleDetailRes{
		ID:           v["id"].String(),
		LicensePlate: v["license_plate"].String(),
		Model:        v["model"].String(),
		Color:        v["color"].String(),
		Brand:        v["brand"].String(),
		Type:         v["type"].String(),
		CreatedAt:    v["created_at"].String(),
	}
	return
}

func (s *sVehicle) Update(ctx context.Context, req *vehicle.VehicleUpdateReq) (res *vehicle.VehicleUpdateRes, err error) {
	userID := g.RequestFromCtx(ctx).GetCtxVar("user_id").String()
	if userID == "" {
		return nil, gerror.New("Unauthorized")
	}

	v, err := dao.Vehicles.Ctx(ctx).Where("id", req.ID).One()
	if err != nil {
		return nil, gerror.New("Database error")
	}
	if v.IsEmpty() {
		return nil, gerror.New("Vehicle not found")
	}
	if v["user_id"].String() != userID {
		return nil, gerror.New("Not owner")
	}

	_, err = dao.Vehicles.Ctx(ctx).
		Where("id", req.ID).
		Data(g.Map{
			"model": req.Model,
			"color": req.Color,
			"brand": req.Brand,
			"type":  req.Type,
		}).Update()
	if err != nil {
		return nil, gerror.New("Failed to update vehicle")
	}

	// Return updated vehicle
	v, err = dao.Vehicles.Ctx(ctx).Where("id", req.ID).One()
	if err != nil {
		return nil, gerror.New("Database error")
	}

	res = &vehicle.VehicleUpdateRes{
		ID:           v["id"].String(),
		LicensePlate: v["license_plate"].String(),
		Model:        v["model"].String(),
		Color:        v["color"].String(),
		Brand:        v["brand"].String(),
		Type:         v["type"].String(),
		CreatedAt:    v["created_at"].String(),
	}
	return
}
