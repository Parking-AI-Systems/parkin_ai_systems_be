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
}
