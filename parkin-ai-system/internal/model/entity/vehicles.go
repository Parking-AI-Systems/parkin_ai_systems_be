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
}
