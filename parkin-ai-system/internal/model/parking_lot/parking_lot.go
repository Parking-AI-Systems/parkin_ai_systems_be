package parking_lot

// ParkingLotListReq is the request struct for listing parking lots
// You can add filter fields if needed
// Example: Name string `json:"name"`
type ParkingLotListReq struct{}

type ParkingLotListRes struct {
	Lots []ParkingLotInfo `json:"lots"`
}

type ParkingLotAddReq struct {
	Name           string  `json:"name"`
	Address        string  `json:"address"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	TotalSlots     int     `json:"total_slots"`
	PricePerHour   float64 `json:"price_per_hour"`
	Description    string  `json:"description"`
	OpenTime       string  `json:"open_time"`
	CloseTime      string  `json:"close_time"`
	ImageUrl       string  `json:"image_url"`
}

type ParkingLotAddRes struct {
	LotID string `json:"lot_id"`
}

type ParkingLotUpdateReq struct {
	Id             string  `json:"id"`
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

type ParkingLotDeleteReq struct {
	Id string `json:"id"`
}

type ParkingLotDeleteRes struct {
	Success bool `json:"success"`
}

type ParkingLotDetailReq struct {
	Id string `json:"id"`
}

type ParkingLotDetailRes struct {
	Lot     *ParkingLotInfo        `json:"lot"`
	Slots   []ParkingSlotInfo      `json:"slots"`
	Images  []ParkingLotImage      `json:"images"`
	Reviews []ParkingLotReview     `json:"reviews"`
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
