package parking_lot

import "github.com/gogf/gf/v2/frame/g"

type ParkingLotUpdateReq struct {
	g.Meta         `path:"/parking-lots/{id}" method:"put" tags:"ParkingLot" summary:"Cập nhật bãi" security:"BearerAuth"`
	Id             string  `json:"id" v:"required"`
	Name           string  `json:"name"`
	Address        string  `json:"address"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	TotalSlots     int     `json:"total_slots"`
	AvailableSlots int     `json:"available_slots"`
	PricePerHour   float64 `json:"price_per_hour"`
	Description    string  `json:"description"`
	OpenTime       string  `json:"open_time"`
	CloseTime      string  `json:"close_time"`
	ImageUrl       string  `json:"image_url"`
	IsActive       bool    `json:"is_active"`
	IsVerified     bool    `json:"is_verified"`
}

type ParkingLotUpdateRes struct {
	Success bool `json:"success"`
}
