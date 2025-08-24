// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Favorites is the golang structure for table favorites.
type Favorites struct {
	Id        int64       `json:"id"        orm:"id"         description:""`
	UserId    int64       `json:"userId"    orm:"user_id"    description:""`
	LotId     int64       `json:"lotId"     orm:"lot_id"     description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
}
