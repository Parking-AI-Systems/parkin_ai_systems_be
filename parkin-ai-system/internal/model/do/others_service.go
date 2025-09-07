// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// OthersService is the golang structure of table others_service for DAO operations like Where/Data.
type OthersService struct {
	g.Meta          `orm:"table:others_service, do:true"`
	Id              interface{} //
	LotId           interface{} //
	Name            interface{} //
	Description     interface{} //
	Price           interface{} //
	DurationMinutes interface{} //
	IsActive        interface{} //
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
	DeletedAt       *gtime.Time //
}
