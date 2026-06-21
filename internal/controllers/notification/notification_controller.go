package notification

import (
	"context"
	"errors"

	notificationBiz "github.com/armylong/armylong-go/internal/business/notification"
	notificationCs "github.com/armylong/armylong-go/internal/cs/notification"
	"github.com/armylong/armylong-go/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type NotificationController struct{}

func (c *NotificationController) ActionList(ctx *gin.Context, req *notificationCs.ListNotificationsRequest) (*notificationCs.ListNotificationsResponse, error) {
	uid := middlewares.GetLoginUID(ctx)
	if uid == 0 {
		return nil, errors.New("请先登录")
	}
	return notificationBiz.NotificationBusiness.List(uid, req)
}

func (c *NotificationController) ActionGetDetail(ctx *gin.Context, req *notificationCs.GetNotificationDetailRequest) (*notificationCs.NotificationDetailResponse, error) {
	uid := middlewares.GetLoginUID(ctx)
	if uid == 0 {
		return nil, errors.New("请先登录")
	}
	return notificationBiz.NotificationBusiness.GetDetail(uid, req.ID)
}

func (c *NotificationController) ActionMarkAsRead(ctx context.Context, req *notificationCs.MarkAsReadRequest) error {
	uid := middlewares.GetLoginUIDFromContext(ctx)
	if uid == 0 {
		return errors.New("请先登录")
	}
	return notificationBiz.NotificationBusiness.MarkAsRead(uid, req.IDs)
}

func (c *NotificationController) ActionMarkAllAsRead(ctx context.Context, req *notificationCs.MarkAllAsReadRequest) error {
	uid := middlewares.GetLoginUIDFromContext(ctx)
	if uid == 0 {
		return errors.New("请先登录")
	}
	return notificationBiz.NotificationBusiness.MarkAllAsRead(uid)
}

func (c *NotificationController) ActionDeleteRead(ctx context.Context, req *notificationCs.DeleteReadNotificationsRequest) error {
	uid := middlewares.GetLoginUIDFromContext(ctx)
	if uid == 0 {
		return errors.New("请先登录")
	}
	return notificationBiz.NotificationBusiness.DeleteRead(uid)
}

func (c *NotificationController) ActionDelete(ctx context.Context, req *notificationCs.DeleteNotificationsRequest) error {
	uid := middlewares.GetLoginUIDFromContext(ctx)
	if uid == 0 {
		return errors.New("请先登录")
	}
	return notificationBiz.NotificationBusiness.Delete(uid, req.IDs)
}

func (c *NotificationController) ActionUnreadCount(ctx context.Context, req *notificationCs.UnreadCountRequest) (*notificationCs.UnreadCountResponse, error) {
	uid := middlewares.GetLoginUIDFromContext(ctx)
	if uid == 0 {
		return nil, errors.New("请先登录")
	}
	return notificationBiz.NotificationBusiness.UnreadCount(uid)
}
