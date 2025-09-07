// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ParkingLotImages is the golang structure for table parking_lot_images.
type ParkingLotImages struct {
	Id         int64       `json:"id"         orm:"id"          description:""`
	LotId      int64       `json:"lotId"      orm:"lot_id"      description:""`
	ImageUrl   string      `json:"imageUrl"   orm:"image_url"   description:""`
	UploadedBy int64       `json:"uploadedBy" orm:"uploaded_by" description:""`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"  description:""`
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at"  description:""`
}
