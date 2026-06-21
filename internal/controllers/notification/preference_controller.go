package notification

import (
	"context"
	"errors"

	preferenceBiz "github.com/armylong/armylong-go/internal/business/notification"
	notificationCs "github.com/armylong/armylong-go/internal/cs/notification"
	"github.com/armylong/armylong-go/internal/middlewares"
	"github.com/gin-gonic/gin"
)

type PreferenceController struct{}

func (c *PreferenceController) ActionGetPreferences(ctx *gin.Context, req *notificationCs.GetNotificationPreferencesRequest) (*notificationCs.GetNotificationPreferencesResponse, error) {
	uid := middlewares.GetLoginUID(ctx)
	if uid == 0 {
		return nil, errors.New("请先登录")
	}
	return preferenceBiz.PreferenceBusiness.GetPreferences(uid)
}

func (c *PreferenceController) ActionUpdatePreference(ctx context.Context, req *notificationCs.UpdateNotificationPreferenceRequest) error {
	uid := middlewares.GetLoginUIDFromContext(ctx)
	if uid == 0 {
		return errors.New("请先登录")
	}
	return preferenceBiz.PreferenceBusiness.UpdatePreference(uid, req)
}

func (c *PreferenceController) ActionBatchUpdatePreferences(ctx context.Context, req *notificationCs.BatchUpdateNotificationPreferencesRequest) (*notificationCs.BatchUpdateNotificationPreferencesResponse, error) {
	uid := middlewares.GetLoginUIDFromContext(ctx)
	if uid == 0 {
		return nil, errors.New("请先登录")
	}
	return preferenceBiz.PreferenceBusiness.BatchUpdatePreferences(uid, req)
}
