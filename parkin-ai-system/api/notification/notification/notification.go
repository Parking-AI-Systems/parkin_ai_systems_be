package notification

import "github.com/gogf/gf/v2/frame/g"

type NotificationListReq struct {
	g.Meta   `path:"/notifications" tags:"Notification" method:"GET" summary:"List notifications" description:"Retrieves a paginated list of notifications for the authenticated user." middleware:"middleware.Auth"`
	IsRead   *bool `json:"is_read" v:"boolean#IsRead must be a boolean"`
	Page     int   `json:"page" v:"min:1#Page must be at least 1"`
	PageSize int   `json:"page_size" v:"min:1|max:100#Page size must be between 1 and 100"`
}

type NotificationItem struct {
	Id             int64  `json:"id"`
	UserId         int64  `json:"user_id"`
	Type           string `json:"type"`
	Content        string `json:"content"`
	RelatedOrderId int64  `json:"related_order_id"`
	IsRead         bool   `json:"is_read"`
	CreatedAt      string `json:"created_at"`
	RelatedInfo    string `json:"related_info"`
}

type NotificationListRes struct {
	List  []NotificationItem `json:"list"`
	Total int                `json:"total"`
}

type NotificationGetReq struct {
	g.Meta `path:"/notifications/:id" tags:"Notification" method:"GET" summary:"Get notification details" description:"Retrieves details of a specific notification by ID for the authenticated user." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Notification ID is required|Notification ID must be positive"`
}

type NotificationGetRes struct {
	Notification NotificationItem `json:"notification"`
}

type NotificationMarkReadReq struct {
	g.Meta `path:"/notifications/mark-read" tags:"Notification" method:"POST" summary:"Mark notifications as read" description:"Marks one or more notifications as read for the authenticated user." middleware:"middleware.Auth"`
	Ids    []int64 `json:"ids" v:"required#At least one notification ID is required"`
}

type NotificationMarkReadRes struct {
	Message string `json:"message"`
}

type NotificationDeleteReq struct {
	g.Meta `path:"/notifications/:id" tags:"Notification" method:"DELETE" summary:"Delete a notification" description:"Permanently deletes a notification. Admin only." middleware:"middleware.Auth"`
	Id     int64 `json:"id" v:"required|min:1#Notification ID is required|Notification ID must be positive"`
}

type NotificationDeleteRes struct {
	Message string `json:"message"`
}
