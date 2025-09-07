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
	DeletedAt      *gtime.Time `json:"deletedAt"      orm:"deleted_at"      description:""`
}

type ParkingLotImageInput struct {
	ImageUrl    string `json:"imageUrl"`
	Description string `json:"description"`
}

type ParkingLotAddReq struct {
	Name         string               `json:"name"`
	Address      string               `json:"address"`
	Latitude     float64              `json:"latitude"`
	Longitude    float64              `json:"longitude"`
	IsVerified   bool                 `json:"isVerified"`
	IsActive     bool                 `json:"isActive"`
	TotalSlots   int                  `json:"totalSlots"`
	PricePerHour float64              `json:"pricePerHour"`
	Description  string               `json:"description"`
	OpenTime     *gtime.Time          `json:"openTime"`
	CloseTime    *gtime.Time          `json:"closeTime"`
	Images       []ParkingLotImageInput `json:"images"`
}

type ParkingLotAddRes struct {
	Id int64 `json:"id"`
}

type ParkingLotListReq struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
	IsActive  bool    `json:"isActive"`
	Page      int     `json:"page"`
	PageSize  int     `json:"pageSize"`
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
	Id            int64                `json:"id"`
	Name          string               `json:"name"`
	Address       string               `json:"address"`
	Latitude      float64              `json:"latitude"`
	Longitude     float64              `json:"longitude"`
	OwnerId       int64                `json:"owner_id"`
	IsVerified    bool                 `json:"is_verified"`
	IsActive      bool                 `json:"is_active"`
	TotalSlots    int                  `json:"total_slots"`
	AvailableSlots int                  `json:"available_slots"`
	PricePerHour  float64              `json:"price_per_hour"`
	Description   string               `json:"description"`
	OpenTime      string               `json:"open_time"`
	CloseTime     string               `json:"close_time"`
	Images        []ParkingLotImageItem `json:"images"`
	CreatedAt     string               `json:"created_at"`
	UpdatedAt     string               `json:"updated_at"`
	DeletedAt     string               `json:"deleted_at,omitempty"`
}

type ParkingLotListRes struct {
	List  []ParkingLotItem `json:"list"`
	Total int              `json:"total"`
}

type ParkingLotGetReq struct {
	Id int64 `json:"id"`
}

type ParkingLotGetRes struct {
	Lot ParkingLotItem `json:"lot"`
}

type ParkingLotUpdateReq struct {
	Id           int64       `json:"id"`
	Name         string      `json:"name"`
	Address      string      `json:"address"`
	Latitude     float64     `json:"latitude"`
	Longitude    float64     `json:"longitude"`
	IsVerified   *bool       `json:"is_verified"`
	IsActive     *bool       `json:"is_active"`
	TotalSlots   int         `json:"total_slots"`
	PricePerHour float64     `json:"price_per_hour"`
	Description  string      `json:"description"`
	OpenTime     *gtime.Time `json:"open_time"`
	CloseTime    *gtime.Time `json:"close_time"`
	ImageUrl     string      `json:"image_url"`
}

type ParkingLotUpdateRes struct {
	Lot ParkingLotItem `json:"lot"`
}

type ParkingLotDeleteReq struct {
	Id int64 `json:"id"`
}

type ParkingLotDeleteRes struct {
	Message string `json:"message"`
}

type ParkingLotImageDeleteReq struct {
	Id int64 `json:"id"`
}

type ParkingLotImageDeleteRes struct {
	Message string `json:"message"`
}