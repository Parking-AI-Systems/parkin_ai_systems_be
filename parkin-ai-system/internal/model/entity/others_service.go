// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// OthersService is the golang structure for table others_service.
type OthersService struct {
	Id              int64       `json:"id"              orm:"id"               description:""`
	LotId           int64       `json:"lotId"           orm:"lot_id"           description:""`
	Name            string      `json:"name"            orm:"name"             description:""`
	Description     string      `json:"description"     orm:"description"      description:""`
	Price           float64     `json:"price"           orm:"price"            description:""`
	DurationMinutes int         `json:"durationMinutes" orm:"duration_minutes" description:""`
	IsActive        bool        `json:"isActive"        orm:"is_active"        description:""`
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       description:""`
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       description:""`
	DeletedAt      *gtime.Time `json:"deletedAt"      orm:"deleted_at"      description:""`
}

type OthersServiceAddReq struct {
	LotId           int64   `json:"lotId"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	DurationMinutes int     `json:"durationMinutes"`
	IsActive        bool    `json:"isActive"`
}

type OthersServiceAddRes struct {
	Id int64 `json:"id"`
}

type OthersServiceListReq struct {
	LotId    int64 `json:"lotId"`
	IsActive bool  `json:"isActive"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

type OthersServiceItem struct {
	Id              int64   `json:"id"`
	LotId           int64   `json:"lot_id"`
	LotName         string  `json:"lot_name"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	DurationMinutes int     `json:"duration_minutes"`
	IsActive        bool    `json:"is_active"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	DeletedAt       string  `json:"deleted_at"`
}

type OthersServiceListRes struct {
	List  []OthersServiceItem `json:"list"`
	Total int                 `json:"total"`
}

type OthersServiceGetReq struct {
	Id int64 `json:"id"`
}

type OthersServiceGetRes struct {
	Service OthersServiceItem `json:"service"`
}

type OthersServiceUpdateReq struct {
	Id              int64   `json:"id"`
	LotId           int64   `json:"lotId"`
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	DurationMinutes int     `json:"durationMinutes"`
	IsActive        *bool   `json:"isActive"`
}

type OthersServiceUpdateRes struct {
	Service OthersServiceItem `json:"service"`
}

type OthersServiceDeleteReq struct {
	Id int64 `json:"id"`
}

type OthersServiceDeleteRes struct {
	Message string `json:"message"`
}