// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ParkingLots is the golang structure for table parking_lots.
type ParkingLots struct {
	Id             int64       `json:"id"             orm:"id"              description:""`
	Name           string      `json:"name"           orm:"name"            description:""`
	Address        string      `json:"address"        orm:"address"         description:""`
	Latitude       float64     `json:"latitude"       orm:"latitude"        description:""`
	Longitude      float64     `json:"longitude"      orm:"longitude"       description:""`
	OwnerId        int64       `json:"ownerId"        orm:"owner_id"        description:""`
	IsVerified     bool        `json:"isVerified"     orm:"is_verified"     description:""`
	IsActive       bool        `json:"isActive"       orm:"is_active"       description:""`
	TotalSlots     int         `json:"totalSlots"     orm:"total_slots"     description:""`
	AvailableSlots int         `json:"availableSlots" orm:"available_slots" description:""`
	PricePerHour   float64     `json:"pricePerHour"   orm:"price_per_hour"  description:""`
	Description    string      `json:"description"    orm:"description"     description:""`
	OpenTime       *gtime.Time `json:"openTime"       orm:"open_time"       description:""`
	CloseTime      *gtime.Time `json:"closeTime"      orm:"close_time"      description:""`
	ImageUrl       string      `json:"imageUrl"       orm:"image_url"       description:""`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:""`
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:""`
}
