package vehicle

import (
	"context"
	"parkin-ai-system/api/vehicle/vehicle"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
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

	vehicleID := guid.S()
	_, err = dao.Vehicles.Ctx(ctx).Data(g.Map{
		"id":            vehicleID,
		"user_id":       userID,
		"license_plate": req.LicensePlate,
		"model":         req.Model,
		"color":         req.Color,
		"brand":         req.Brand,
		"type":          req.Type,
		"user":          userID,
	}).Insert()
	if err != nil {
		return nil, gerror.New("Failed to add vehicle")
	}

	res = &vehicle.VehicleAddRes{
		VehicleID: vehicleID,
	}
	return
}
