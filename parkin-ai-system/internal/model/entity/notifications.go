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
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"       description:""`
	DeletedAt      *gtime.Time `json:"deletedAt"      orm:"deleted_at"       description:""`
}

type NotificationListReq struct {
	IsRead   *bool `json:"isRead"`
	Page     int   `json:"page"`
	PageSize int   `json:"pageSize"`
}

type NotificationItem struct {
	Id             int64  `json:"id"`
	UserId         int64  `json:"user_id"`
	Type           string `json:"type"`
	Content        string `json:"content"`
	RelatedOrderId int64  `json:"related_order_id"`
	IsRead         bool   `json:"is_read"`
	CreatedAt      string `json:"created_at"`
	UpdatedAt      string `json:"updated_at"`
	DeletedAt      string `json:"deleted_at"`
	RelatedInfo    string `json:"related_info"`
}

type NotificationListRes struct {
	List  []NotificationItem `json:"list"`
	Total int                `json:"total"`
}

type NotificationGetReq struct {
	Id int64 `json:"id"`
}

type NotificationGetRes struct {
	Notification NotificationItem `json:"notification"`
}

type NotificationMarkReadReq struct {
	Ids []int64 `json:"ids"`
}

type NotificationMarkReadRes struct {
	Message string `json:"message"`
}

type NotificationDeleteReq struct {
	Id int64 `json:"id"`
}

type NotificationDeleteRes struct {
	Message string `json:"message"`
}