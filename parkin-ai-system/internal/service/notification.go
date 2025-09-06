// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"parkin-ai-system/internal/model/entity"
)

type (
	INotification interface {
		NotificationList(ctx context.Context, req *entity.NotificationListReq) (*entity.NotificationListRes, error)
		NotificationGet(ctx context.Context, req *entity.NotificationGetReq) (*entity.NotificationItem, error)
		NotificationMarkRead(ctx context.Context, req *entity.NotificationMarkReadReq) (*entity.NotificationMarkReadRes, error)
		NotificationDelete(ctx context.Context, req *entity.NotificationDeleteReq) (*entity.NotificationDeleteRes, error)
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
