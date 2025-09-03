package notification

import (
	"github.com/gogf/gf/v2/frame/g"
)

type NotificationAddReq struct {
	g.Meta      `path:"/notifications" method:"post" tags:"Notification" summary:"Tạo thông báo" security:"BearerAuth"`
	Type        string `json:"type" v:"required"`
	Content     string `json:"content" v:"required"`
	RelatedOrderId int64 `json:"related_order_id"`
}

type NotificationAddRes struct {
	Id int64 `json:"id"`
}

type NotificationListReq struct {
	g.Meta `path:"/notifications" method:"get" tags:"Notification" summary:"Danh sách thông báo" security:"BearerAuth"`
}

type NotificationListRes struct {
	List []NotificationItem `json:"list"`
}

type NotificationItem struct {
	Id            int64  `json:"id"`
	Type          string `json:"type"`
	Content       string `json:"content"`
	RelatedOrderId int64 `json:"related_order_id"`
	IsRead        bool   `json:"is_read"`
	CreatedAt     string `json:"created_at"`
}

type NotificationUpdateReq struct {
	g.Meta   `path:"/notifications/{id}" method:"put" tags:"Notification" summary:"Cập nhật trạng thái thông báo" security:"BearerAuth"`
	Id       int64  `json:"id" v:"required"`
	IsRead   bool   `json:"is_read"`
}

type NotificationUpdateRes struct {
	Success bool `json:"success"`
}

type NotificationDeleteReq struct {
	g.Meta   `path:"/notifications/{id}" method:"delete" tags:"Notification" summary:"Xoá thông báo" security:"BearerAuth"`
	Id       int64  `json:"id" v:"required"`
}

type NotificationDeleteRes struct {
	Success bool `json:"success"`
}
