package vehicle

import (
	"github.com/gogf/gf/v2/frame/g"
)

type VehicleDetailReq struct {
	g.Meta `path:"/vehicles/{id}" method:"get" tags:"Vehicle" summary:"Get vehicle detail" security:"BearerAuth"`
	ID     string `json:"id" in:"path"`
}

type VehicleDetailRes struct {
	ID           string `json:"id"`
	LicensePlate string `json:"license_plate"`
	Model        string `json:"model"`
	Color        string `json:"color"`
	Brand        string `json:"brand"`
	Type         string `json:"type"`
	CreatedAt    string `json:"created_at"`
}
