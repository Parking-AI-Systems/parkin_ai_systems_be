package notification

import (
	"context"

	"parkin-ai-system/api/notification/notification"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerNotification) NotificationDelete(ctx context.Context, req *notification.NotificationDeleteReq) (res *notification.NotificationDeleteRes, err error) {
	// Map API request to entity request
	input := &entity.NotificationDeleteReq{
		Id: req.Id,
	}

	// Call service
	deleteRes, err := service.Notification().NotificationDelete(ctx, input)
	if err != nil {
		return nil, err
	}

	// Map entity response to API response
	res = &notification.NotificationDeleteRes{
		Message: deleteRes.Message,
	}
	if r := g.RequestFromCtx(ctx); r != nil {
		r.Response.WriteJson(res)
		return nil, nil
	}
	return res, nil
}
