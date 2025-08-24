// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Vehicles is the golang structure of table vehicles for DAO operations like Where/Data.
type Vehicles struct {
	g.Meta       `orm:"table:vehicles, do:true"`
	Id           interface{} //
	UserId       interface{} //
	LicensePlate interface{} //
	Brand        interface{} //
	Model        interface{} //
	Color        interface{} //
	Type         interface{} //
	CreatedAt    *gtime.Time //
}
