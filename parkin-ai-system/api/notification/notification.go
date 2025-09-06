// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package notification

import (
	"context"

	"parkin-ai-system/api/notification/notification"
)

type INotificationNotification interface {
	NotificationList(ctx context.Context, req *notification.NotificationListReq) (res *notification.NotificationListRes, err error)
	NotificationGet(ctx context.Context, req *notification.NotificationGetReq) (res *notification.NotificationGetRes, err error)
	NotificationMarkRead(ctx context.Context, req *notification.NotificationMarkReadReq) (res *notification.NotificationMarkReadRes, err error)
	NotificationDelete(ctx context.Context, req *notification.NotificationDeleteReq) (res *notification.NotificationDeleteRes, err error)
}
