package notification

import (
	"context"

	"parkin-ai-system/api/notification/notification"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerNotification) NotificationMarkRead(ctx context.Context, req *notification.NotificationMarkReadReq) (res *notification.NotificationMarkReadRes, err error) {
	// Map API request to entity request
	input := &entity.NotificationMarkReadReq{
		Ids: req.Ids,
	}

	// Call service
	markRes, err := service.Notification().NotificationMarkRead(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &notification.NotificationMarkReadRes{
		Message: markRes.Message,
	}
	return res, nil
}
