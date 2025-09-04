// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"parkin-ai-system/api/notification"
)

type (
	INotification interface {
		NotificationAdd(ctx context.Context, req *notification.NotificationAddReq) (*notification.NotificationAddRes, error)
		NotificationList(ctx context.Context, req *notification.NotificationListReq) (*notification.NotificationListRes, error)
		NotificationUpdate(ctx context.Context, req *notification.NotificationUpdateReq) (*notification.NotificationUpdateRes, error)
		NotificationDelete(ctx context.Context, req *notification.NotificationDeleteReq) (*notification.NotificationDeleteRes, error)
	}
)

var (
	localNotification INotification
)

func Notification() INotification {
	if localNotification == nil {
		panic("implement not found for interface INotification, forgot register?")
	}
	return localNotification
}

func RegisterNotification(i INotification) {
	localNotification = i
}
