package notification

import (
	"context"

	"parkin-ai-system/api/notification/notification"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"
)

func (c *ControllerNotification) NotificationList(ctx context.Context, req *notification.NotificationListReq) (res *notification.NotificationListRes, err error) {
	// Map API request to entity request
	input := &entity.NotificationListReq{
		Page:     req.Page,
		PageSize: req.PageSize,
		IsRead:   req.IsRead,
	}

	// Call service
	listRes, err := service.Notification().NotificationList(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity list to API response
	res = &notification.NotificationListRes{
		List:  make([]notification.NotificationItem, 0, len(listRes.List)),
		Total: listRes.Total,
	}
	for _, item := range listRes.List {
		res.List = append(res.List, notification.NotificationItem{
			Id:             item.Id,
			UserId:         item.UserId,
			Type:           item.Type,
			Content:        item.Content,
			RelatedOrderId: item.RelatedOrderId,
			IsRead:         item.IsRead,
			CreatedAt:      item.CreatedAt,
			RelatedInfo:    item.RelatedInfo,
		})
	}
	return res, nil
}
