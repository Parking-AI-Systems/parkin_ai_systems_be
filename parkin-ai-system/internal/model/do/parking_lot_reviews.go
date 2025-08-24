// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ParkingLotReviews is the golang structure of table parking_lot_reviews for DAO operations like Where/Data.
type ParkingLotReviews struct {
	g.Meta    `orm:"table:parking_lot_reviews, do:true"`
	Id        interface{} //
	LotId     interface{} //
	UserId    interface{} //
	Rating    interface{} //
	Comment   interface{} //
	CreatedAt *gtime.Time //
}
