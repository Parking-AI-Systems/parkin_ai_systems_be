// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Vehicles is the golang structure for table vehicles.
type Vehicles struct {
	Id           int64       `json:"id"           orm:"id"            description:""`
	UserId       int64       `json:"userId"       orm:"user_id"       description:""`
	LicensePlate string      `json:"licensePlate" orm:"license_plate" description:""`
	Brand        string      `json:"brand"        orm:"brand"         description:""`
	Model        string      `json:"model"        orm:"model"         description:""`
	Color        string      `json:"color"        orm:"color"         description:""`
	Type         string      `json:"type"         orm:"type"          description:""`
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:""`
	DeletedAt    *gtime.Time `json:"deletedAt"    orm:"deleted_at"    description:""`
}
type VehicleAddReq struct {
	LicensePlate string `json:"licensePlate"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Color        string `json:"color"`
	Type         string `json:"type"`
}

type VehicleAddRes struct {
	Id int64 `json:"id"`
}

type VehicleListReq struct {
	Type     string `json:"type"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
}

type VehicleItem struct {
	Id           int64  `json:"id"`
	UserId       int64  `json:"user_id"`
	Username     string `json:"username"`
	LicensePlate string `json:"license_plate"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Color        string `json:"color"`
	Type         string `json:"type"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	DeletedAt    string `json:"deleted_at"`
}

type VehicleListRes struct {
	List  []VehicleItem `json:"list"`
	Total int           `json:"total"`
}

type VehicleGetReq struct {
	Id int64 `json:"id"`
}

type VehicleGetRes struct {
	Vehicle VehicleItem `json:"vehicle"`
}

type VehicleUpdateReq struct {
	Id           int64  `json:"id"`
	LicensePlate string `json:"licensePlate"`
	Brand        string `json:"brand"`
	Model        string `json:"model"`
	Color        string `json:"color"`
	Type         string `json:"type"`
}

type VehicleUpdateRes struct {
	Vehicle VehicleItem `json:"vehicle"`
}

type VehicleDeleteReq struct {
	Id int64 `json:"id"`
}

type VehicleDeleteRes struct {
	Message string `json:"message"`
}