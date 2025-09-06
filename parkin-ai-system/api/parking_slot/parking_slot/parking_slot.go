package parking_slot

import "github.com/gogf/gf/v2/frame/g"

type ParkingSlotAddReq struct {
	g.Meta      `path:"/parking-slots" tags:"Parking Slot" method:"POST" summary:"Add a new parking slot" description:"Creates a new parking slot for a parking lot. Admin only." middleware:"middleware.Auth"`
	LotId       int64  `json:"lot_id" v:"required|min:1#Parking lot ID is required|Parking lot ID must be positive"`
	Code        string `json:"code" v:"required|length:1,20#Code is required|Code must be between 1 and 20 characters"`
	IsAvailable bool   `json:"is_available" v:"boolean#IsAvailable must be a boolean"`
	SlotType    string `json:"slot_type" v:"required#Slot type is required"`
	Floor       string `json:"floor" v:"required|length:1,10#Floor is required|Floor must be between 1 and 10 characters"`
}

type ParkingSlotAddRes struct {
	Id int64 `json:"id"`
}

type ParkingSlotListReq struct {
	g.Meta      `path:"/parking-slots" tags:"Parking Slot" method:"GET" summary:"List parking slots" description:"Retrieves a paginated list of parking slots with optional filters." middleware:"middleware.Auth"`
	LotId       int64  `json:"lot_id" v:"min:0#Parking lot ID must be non-negative"`
	IsAvailable *bool  `json:"is_available" v:"boolean#IsAvailable must be a boolean"`
	SlotType    string `json:"slot_type"`
	Page        int    `json:"page" v:"min:1#Page must be at least 1"`
	PageSize    int    `json:"page_size" v:"min:1|max:100#Page size must be between 1 and 100"`
}

type ParkingSlotItem struct {
	Id          int64  `json:"id"`
	LotId       int64  `json:"lot_id"`
	LotName     string `json:"lot_name"`
	Code        string `json:"code"`
	IsAvailable bool   `json:"is_available"`
	SlotType    string `json:"slot_type"`
	Floor       string `json:"floor"`
	CreatedAt   string `json:"created_at"`
}

type ParkingSlotListRes struct {
	List  []ParkingSlotItem `json:"list"`
	Total int               `json:"total"`
}

type ParkingSlotGetReq struct {
	g.Meta `path:"/parking-slots/:id" tags:"Parking Slot" method:"GET" summary:"Get parking slot details" description:"Retrieves details of a specific parking slot by ID." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Parking slot ID is required|Parking slot ID must be positive"`
}

type ParkingSlotGetRes struct {
	Slot ParkingSlotItem `json:"slot"`
}

type ParkingSlotUpdateReq struct {
	g.Meta      `path:"/parking-slots/:id" tags:"Parking Slot" method:"PATCH" summary:"Update a parking slot" description:"Updates the details of a parking slot. Admin only." middleware:"middleware.Auth"`
	Id          int64  `json:"id" v:"required|min:1#Parking slot ID is required|Parking slot ID must be positive"`
	LotId       int64  `json:"lot_id" v:"min:0#Parking lot ID must be non-negative"`
	Code        string `json:"code" v:"length:0,20#Code must be less than 20 characters"`
	IsAvailable *bool  `json:"is_available" v:"boolean#IsAvailable must be a boolean"`
	SlotType    string `json:"slot_type"`
	Floor       string `json:"floor" v:"length:0,10#Floor must be less than 10 characters"`
}

type ParkingSlotUpdateRes struct {
	Slot ParkingSlotItem `json:"slot"`
}

type ParkingSlotDeleteReq struct {
	g.Meta `path:"/parking-slots/:id" tags:"Parking Slot" method:"DELETE" summary:"Delete a parking slot" description:"Permanently deletes a parking slot. Admin only." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Parking slot ID is required|Parking slot ID must be positive"`
}

type ParkingSlotDeleteRes struct {
	Message string `json:"message"`
}
