// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Notifications is the golang structure for table notifications.
type Notifications struct {
	Id             int64       `json:"id"             orm:"id"               description:""`
	UserId         int64       `json:"userId"         orm:"user_id"          description:""`
	Type           string      `json:"type"           orm:"type"             description:""`
	Content        string      `json:"content"        orm:"content"          description:""`
	RelatedOrderId int64       `json:"relatedOrderId" orm:"related_order_id" description:""`
	IsRead         bool        `json:"isRead"         orm:"is_read"          description:""`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"       description:""`
}
