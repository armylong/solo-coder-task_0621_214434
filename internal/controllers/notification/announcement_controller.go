package notification

import (
	"context"
	"errors"

	notificationBiz "github.com/armylong/armylong-go/internal/business/notification"
	notificationCs "github.com/armylong/armylong-go/internal/cs/notification"
	"github.com/armylong/armylong-go/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type AnnouncementController struct{}

func (c *AnnouncementController) ActionCreate(ctx *gin.Context, req *notificationCs.CreateAnnouncementRequest) (*notificationCs.AnnouncementDetailResponse, error) {
	uid := middlewares.GetLoginUID(ctx)
	if uid == 0 {
		return nil, errors.New("请先登录")
	}
	return notificationBiz.AnnouncementBusiness.Create(uid, req)
}

func (c *AnnouncementController) ActionUpdate(ctx *gin.Context, req *notificationCs.UpdateAnnouncementRequest) (*notificationCs.AnnouncementDetailResponse, error) {
	return notificationBiz.AnnouncementBusiness.Update(req)
}

func (c *AnnouncementController) ActionDelete(ctx *gin.Context, req *notificationCs.DeleteAnnouncementRequest) error {
	return notificationBiz.AnnouncementBusiness.Delete(req.ID)
}

func (c *AnnouncementController) ActionPublish(ctx context.Context, req *notificationCs.PublishAnnouncementRequest) (*notificationCs.AnnouncementDetailResponse, error) {
	return notificationBiz.AnnouncementBusiness.Publish(req.ID)
}

func (c *AnnouncementController) ActionArchive(ctx context.Context, req *notificationCs.ArchiveAnnouncementRequest) (*notificationCs.AnnouncementDetailResponse, error) {
	return notificationBiz.AnnouncementBusiness.Archive(req.ID)
}

func (c *AnnouncementController) ActionList(ctx context.Context, req *notificationCs.ListAnnouncementsRequest) (*notificationCs.ListAnnouncementsResponse, error) {
	return notificationBiz.AnnouncementBusiness.List(req)
}

func (c *AnnouncementController) ActionGetDetail(ctx context.Context, req *notificationCs.GetAnnouncementDetailRequest) (*notificationCs.AnnouncementDetailResponse, error) {
	return notificationBiz.AnnouncementBusiness.GetByID(req.ID)
}
