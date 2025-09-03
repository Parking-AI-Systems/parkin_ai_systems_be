package service

import (
	"parkin-ai-system/api/notification"
)

type INotification interface {
	notification.INotification
}

var localNotification INotification

func Notification() INotification {
	if localNotification == nil {
		panic("implement not found for interface INotification, forgot register?")
	}
	return localNotification
}

func RegisterNotification(i INotification) {
	localNotification = i
}
