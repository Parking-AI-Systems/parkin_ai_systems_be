package notification

import (
	"context"

	"parkin-ai-system/api/notification/notification"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerNotification) NotificationGet(ctx context.Context, req *notification.NotificationGetReq) (res *notification.NotificationGetRes, err error) {
	// Map API request to entity request
	input := &entity.NotificationGetReq{
		Id: req.Id,
	}

	// Call service
	noti, err := service.Notification().NotificationGet(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &notification.NotificationGetRes{
		Notification: notification.NotificationItem{
			Id:             noti.Id,
			UserId:         noti.UserId,
			Type:           noti.Type,
			Content:        noti.Content,
			RelatedOrderId: noti.RelatedOrderId,
			IsRead:         noti.IsRead,
			CreatedAt:      noti.CreatedAt,
			RelatedInfo:    noti.RelatedInfo,
		},
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
