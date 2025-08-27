package parking_lot

import (
	"github.com/gogf/gf/v2/frame/g"
)

type ParkingLotDetailReq struct {
	g.Meta `path:"/parking-lots/{id}" method:"get" tags:"ParkingLot" summary:"Chi tiết bãi" security:"BearerAuth"`
	Id     string `json:"id" v:"required"`
}

type ParkingLotDetailRes struct {
	Lot     *ParkingLotInfo    `json:"lot"`
	Slots   []ParkingSlotInfo  `json:"slots"`
	Images  []ParkingLotImage  `json:"images"`
	Reviews []ParkingLotReview `json:"reviews"`
}

type ParkingLotInfo struct {
	Id             string  `json:"id"`
	Name           string  `json:"name"`
	Address        string  `json:"address"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	TotalSlots     int     `json:"total_slots"`
	AvailableSlots int     `json:"available_slots"`
	Description    string  `json:"description"`
}

type ParkingSlotInfo struct {
	Id         string `json:"id"`
	SlotNumber string `json:"slot_number"`
	Status     string `json:"status"`
}

type ParkingLotImage struct {
	Id       string `json:"id"`
	ImageUrl string `json:"image_url"`
}

type ParkingLotReview struct {
	Id      string `json:"id"`
	Score   int    `json:"score"`
	Comment string `json:"comment"`
}
