// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ParkingLotReviews is the golang structure for table parking_lot_reviews.
type ParkingLotReviews struct {
	Id        int64       `json:"id"        orm:"id"         description:""`
	LotId     int64       `json:"lotId"     orm:"lot_id"     description:""`
	UserId    int64       `json:"userId"    orm:"user_id"    description:""`
	Rating    int         `json:"rating"    orm:"rating"     description:""`
	Comment   string      `json:"comment"   orm:"comment"    description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
}
