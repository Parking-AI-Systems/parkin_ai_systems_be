package notification

import (
	"context"
	"parkin-ai-system/api/notification"
	"parkin-ai-system/internal/service"
	"github.com/gogf/gf/v2/errors/gerror"
)

type ControllerNotification struct{}

func NewNotification() *ControllerNotification {
	return &ControllerNotification{}
}

func (c *ControllerNotification) NotificationAdd(ctx context.Context, req *notification.NotificationAddReq) (res *notification.NotificationAddRes, err error) {
	res, err = service.Notification().NotificationAdd(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerNotification) NotificationList(ctx context.Context, req *notification.NotificationListReq) (res *notification.NotificationListRes, err error) {
	res, err = service.Notification().NotificationList(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerNotification) NotificationUpdate(ctx context.Context, req *notification.NotificationUpdateReq) (res *notification.NotificationUpdateRes, err error) {
	res, err = service.Notification().NotificationUpdate(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}

func (c *ControllerNotification) NotificationDelete(ctx context.Context, req *notification.NotificationDeleteReq) (res *notification.NotificationDeleteRes, err error) {
	res, err = service.Notification().NotificationDelete(ctx, req)
	if err != nil {
		return nil, gerror.New(err.Error())
	}
	return
}
