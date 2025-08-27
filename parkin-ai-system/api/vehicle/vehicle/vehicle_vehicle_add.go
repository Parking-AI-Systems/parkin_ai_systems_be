package vehicle

import (
	"github.com/gogf/gf/v2/frame/g"
)

type VehicleAddReq struct {
	g.Meta       `path:"/vehicles" method:"post" tags:"Vehicle" summary:"Add vehicle" security:"BearerAuth"`
	LicensePlate string `json:"license_plate" v:"required|length:6,12"`
	Model        string `json:"model" v:"required|length:1,32"`
	Color        string `json:"color" v:"required|length:1,32"`
	Brand        string `json:"brand" v:"required|length:1,32"`
	Type         string `json:"type" v:"required|length:1,32"`
}

type VehicleAddRes struct {
	VehicleID string `json:"vehicle_id"`
}
