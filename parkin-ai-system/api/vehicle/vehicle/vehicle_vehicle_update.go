package vehicle

import (
	"github.com/gogf/gf/v2/frame/g"
)

type VehicleUpdateReq struct {
	g.Meta `path:"/vehicles/{id}" method:"patch" tags:"Vehicle" summary:"Update vehicle" security:"BearerAuth"`
	ID     string `json:"id" in:"path"`
	Model  string `json:"model" v:"required|length:1,32"`
	Color  string `json:"color" v:"required|length:1,32"`
	Brand  string `json:"brand" v:"length:0,32"`
	Type   string `json:"type" v:"length:0,32"`
}

type VehicleUpdateRes struct {
	ID           string `json:"id"`
	LicensePlate string `json:"license_plate"`
	Model        string `json:"model"`
	Color        string `json:"color"`
	Brand        string `json:"brand"`
	Type         string `json:"type"`
	CreatedAt    string `json:"created_at"`
}
