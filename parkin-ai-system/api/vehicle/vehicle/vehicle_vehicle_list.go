package vehicle

import (
	"github.com/gogf/gf/v2/frame/g"
)

type VehicleListReq struct {
	g.Meta `path:"/vehicles" method:"get" tags:"Vehicle" summary:"List user vehicles" security:"BearerAuth"`
}

type VehicleListItem struct {
	ID           string `json:"id"`
	LicensePlate string `json:"license_plate"`
	Model        string `json:"model"`
	Color        string `json:"color"`
	Brand        string `json:"brand"`
	Type         string `json:"type"`
	CreatedAt    string `json:"created_at"`
}

type VehicleListRes struct {
	Vehicles []VehicleListItem `json:"vehicles"`
}
