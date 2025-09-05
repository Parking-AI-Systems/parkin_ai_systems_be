package notification

import (
	"context"
	"parkin-ai-system/api/notification"
	"parkin-ai-system/internal/dao"
	"parkin-ai-system/internal/model/do"
	"parkin-ai-system/internal/model/entity"
	"parkin-ai-system/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

type sNotification struct{}

func init() {
	service.RegisterNotification(&sNotification{})
}

func (s *sNotification) NotificationAdd(ctx context.Context, req *notification.NotificationAddReq) (*notification.NotificationAddRes, error) {
	userId := ctx.Value("user_id")
	n := do.Notifications{}
	gconv.Struct(req, &n)
	n.UserId = userId
	n.CreatedAt = gtime.Now()
	lastId, err := dao.Notifications.Ctx(ctx).Data(n).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	return &notification.NotificationAddRes{Id: gconv.Int64(lastId)}, nil
}

func (s *sNotification) NotificationList(ctx context.Context, req *notification.NotificationListReq) (*notification.NotificationListRes, error) {
	userId := ctx.Value("user_id")
	var notis []entity.Notifications
	err := dao.Notifications.Ctx(ctx).Where("user_id", userId).Order("id desc").Scan(&notis)
	if err != nil {
		return nil, err
	}
	var list []notification.NotificationItem
	for _, n := range notis {
		item := notification.NotificationItem{}
		gconv.Struct(n, &item)
		item.CreatedAt = n.CreatedAt.Format("2006-01-02 15:04:05")
		list = append(list, item)
	}
	return &notification.NotificationListRes{List: list}, nil
}

func (s *sNotification) NotificationUpdate(ctx context.Context, req *notification.NotificationUpdateReq) (*notification.NotificationUpdateRes, error) {
	_, err := dao.Notifications.Ctx(ctx).Where("id", req.Id).Data(g.Map{"is_read": req.IsRead}).Update()
	if err != nil {
		return nil, err
	}
	return &notification.NotificationUpdateRes{Success: true}, nil
}

func (s *sNotification) NotificationDelete(ctx context.Context, req *notification.NotificationDeleteReq) (*notification.NotificationDeleteRes, error) {
	_, err := dao.Notifications.Ctx(ctx).Where("id", req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return &notification.NotificationDeleteRes{Success: true}, nil
}
