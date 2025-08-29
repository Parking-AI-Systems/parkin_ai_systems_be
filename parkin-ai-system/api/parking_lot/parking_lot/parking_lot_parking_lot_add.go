package parking_lot

import (
	"github.com/gogf/gf/v2/frame/g"
)

type ParkingLotAddReq struct {
	g.Meta       `path:"/parking-lots" method:"post" tags:"ParkingLot" summary:"Tạo bãi mới" security:"BearerAuth"`
	Name         string  `json:"name" v:"required|length:1,64"`
	Address      string  `json:"address" v:"required|length:1,128"`
	Latitude     float64 `json:"latitude" v:"required"`
	Longitude    float64 `json:"longitude" v:"required"`
	TotalSlots   int     `json:"total_slots" v:"required|min:1"`
	PricePerHour float64 `json:"price_per_hour"`
	Description  string  `json:"description"`
	OpenTime     string  `json:"open_time"`
	CloseTime    string  `json:"close_time"`
	ImageUrl     string  `json:"image_url"`
}

type ParkingLotAddRes struct {
	LotID string `json:"lot_id"`
}
