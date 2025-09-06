package parking_lot

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type ParkingLotImageInput struct {
	ImageUrl string `json:"image_url" v:"required#Image URL is required"`
}

type ParkingLotAddReq struct {
	g.Meta       `path:"/parking-lots" tags:"Parking Lot" method:"POST" summary:"Create a new parking lot" description:"Creates a parking lot with the specified details and optional images. Admin only." middleware:"middleware.Auth"`
	Name         string                 `json:"name" v:"required#Name is required"`
	Address      string                 `json:"address" v:"required#Address is required"`
	Latitude     float64                `json:"latitude" v:"required|min:-90|max:90#Latitude is required|Invalid latitude"`
	Longitude    float64                `json:"longitude" v:"required|min:-180|max:180#Longitude is required|Invalid longitude"`
	OwnerId      int64                  `json:"owner_id" v:"min:0#Owner ID must be non-negative"`
	IsVerified   bool                   `json:"is_verified"`
	IsActive     bool                   `json:"is_active"`
	TotalSlots   int                    `json:"total_slots" v:"required|min:1#Total slots is required|Total slots must be positive"`
	PricePerHour float64                `json:"price_per_hour" v:"required|min:0#Price per hour is required|Price per hour must be non-negative"`
	Description  string                 `json:"description" v:"length:0,1000#Description must be less than 1000 characters"`
	OpenTime     *gtime.Time            `json:"open_time" v:"required#Open time is required"`
	CloseTime    *gtime.Time            `json:"close_time" v:"required#Close time is required"`
	ImageUrl     string                 `json:"image_url" v:"required#Main image URL is required"`
	Images       []ParkingLotImageInput `json:"images"`
}

type ParkingLotAddRes struct {
	Id int64 `json:"id"`
}

type ParkingLotListReq struct {
	g.Meta    `path:"/parking-lots" tags:"Parking Lot" method:"GET" summary:"List parking lots" description:"Retrieves a paginated list of parking lots with optional filters for location and active status." middleware:"middleware.Auth"`
	Latitude  float64 `json:"latitude" v:"min:-90|max:90#Invalid latitude"`
	Longitude float64 `json:"longitude" v:"min:-180|max:180#Invalid longitude"`
	Radius    float64 `json:"radius" v:"min:0#Radius must be non-negative"`
	IsActive  bool    `json:"is_active"`
	Page      int     `json:"page" v:"min:1#Page must be at least 1"`
	PageSize  int     `json:"page_size" v:"min:1|max:100#Page size must be between 1 and 100"`
}

type ParkingLotImageItem struct {
	Id           int64  `json:"id"`
	ParkingLotId int64  `json:"parking_lot_id"`
	LotName      string `json:"lot_name"`
	ImageUrl     string `json:"image_url"`
	Description  string `json:"description"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type ParkingLotItem struct {
	Id             int64                 `json:"id"`
	Name           string                `json:"name"`
	Address        string                `json:"address"`
	Latitude       float64               `json:"latitude"`
	Longitude      float64               `json:"longitude"`
	OwnerId        int64                 `json:"owner_id"`
	IsVerified     bool                  `json:"is_verified"`
	IsActive       bool                  `json:"is_active"`
	TotalSlots     int                   `json:"total_slots"`
	AvailableSlots int                   `json:"available_slots"`
	PricePerHour   float64               `json:"price_per_hour"`
	Description    string                `json:"description"`
	OpenTime       string                `json:"open_time"`
	CloseTime      string                `json:"close_time"`
	ImageUrl       string                `json:"image_url"`
	Images         []ParkingLotImageItem `json:"images"`
	CreatedAt      string                `json:"created_at"`
	UpdatedAt      string                `json:"updated_at"`
}

type ParkingLotListRes struct {
	List  []ParkingLotItem `json:"list"`
	Total int              `json:"total"`
}

type ParkingLotGetReq struct {
	g.Meta `path:"/parking-lots/:id" tags:"Parking Lot" method:"GET" summary:"Get parking lot details" description:"Retrieves details of a specific parking lot by ID, including images." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Parking lot ID is required|Parking lot ID must be positive"`
}

type ParkingLotGetRes struct {
	Lot ParkingLotItem `json:"lot"`
}

type ParkingLotUpdateReq struct {
	g.Meta       `path:"/parking-lots/:id" tags:"Parking Lot" method:"PATCH" summary:"Update a parking lot" description:"Updates the details of a parking lot. Admin only." middleware:"middleware.Auth"`
	Id           int64       `json:"id" v:"required|min:1#Parking lot ID is required|Parking lot ID must be positive"`
	Name         string      `json:"name" v:"length:0,255#Name must be less than 255 characters"`
	Address      string      `json:"address" v:"length:0,255#Address must be less than 255 characters"`
	Latitude     float64     `json:"latitude" v:"min:-90|max:90#Invalid latitude"`
	Longitude    float64     `json:"longitude" v:"min:-180|max:180#Invalid longitude"`
	OwnerId      int64       `json:"owner_id" v:"min:0#Owner ID must be non-negative"`
	IsVerified   *bool       `json:"is_verified"`
	IsActive     *bool       `json:"is_active"`
	TotalSlots   int         `json:"total_slots" v:"min:0#Total slots must be non-negative"`
	PricePerHour float64     `json:"price_per_hour" v:"min:0#Price per hour must be non-negative"`
	Description  string      `json:"description" v:"length:0,1000#Description must be less than 1000 characters"`
	OpenTime     *gtime.Time `json:"open_time"`
	CloseTime    *gtime.Time `json:"close_time"`
	ImageUrl     string      `json:"image_url" v:"length:0,255#Main image URL must be less than 255 characters"`
}

type ParkingLotUpdateRes struct {
	Lot ParkingLotItem `json:"lot"`
}

type ParkingLotDeleteReq struct {
	g.Meta `path:"/parking-lots/:id" tags:"Parking Lot" method:"DELETE" summary:"Delete a parking lot" description:"Soft deletes a parking lot and its images. Admin only." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Parking lot ID is required|Parking lot ID must be positive"`
}

type ParkingLotDeleteRes struct {
	Message string `json:"message"`
}

type ParkingLotImageDeleteReq struct {
	g.Meta `path:"/parking-lot-images/:id" tags:"Parking Lot Image" method:"DELETE" summary:"Delete a parking lot image" description:"Permanently deletes a parking lot image. Admin only." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Image ID is required|Image ID must be positive"`
}

type ParkingLotImageDeleteRes struct {
	Message string `json:"message"`
}
