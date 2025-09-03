package notification

import "context"

type INotification interface {
	NotificationAdd(ctx context.Context, req *NotificationAddReq) (*NotificationAddRes, error)
	NotificationList(ctx context.Context, req *NotificationListReq) (*NotificationListRes, error)
	NotificationUpdate(ctx context.Context, req *NotificationUpdateReq) (*NotificationUpdateRes, error)
	NotificationDelete(ctx context.Context, req *NotificationDeleteReq) (*NotificationDeleteRes, error)
}
