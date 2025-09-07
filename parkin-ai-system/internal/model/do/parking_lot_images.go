// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ParkingLotImages is the golang structure of table parking_lot_images for DAO operations like Where/Data.
type ParkingLotImages struct {
	g.Meta     `orm:"table:parking_lot_images, do:true"`
	Id         interface{} //
	LotId      interface{} //
	ImageUrl   interface{} //
	UploadedBy interface{} //
	CreatedAt  *gtime.Time //
	UpdatedAt  *gtime.Time //
	DeletedAt  *gtime.Time //
}
