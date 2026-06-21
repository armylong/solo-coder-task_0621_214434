package notification

import (
	"errors"

	notificationCs "github.com/armylong/armylong-go/internal/cs/notification"
	notificationModel "github.com/armylong/armylong-go/internal/model/notification"
)

var DefaultNotificationTypes = []string{
	notificationModel.AnnouncementTypeSystem,
	notificationModel.AnnouncementTypeMaintenance,
	notificationModel.AnnouncementTypeActivity,
	notificationModel.AnnouncementTypeUpdate,
}

type preferenceBusiness struct{}

var PreferenceBusiness = &preferenceBusiness{}

func (b *preferenceBusiness) GetPreferences(uid int64) (*notificationCs.GetNotificationPreferencesResponse, error) {
	if uid == 0 {
		return nil, errors.New("请先登录")
	}

	existing, err := notificationModel.TbNotificationPreferenceModel.GetByUid(uid)
	if err != nil {
		return nil, errors.New("获取通知偏好失败: " + err.Error())
	}

	existingMap := make(map[string]*notificationModel.TbNotificationPreference)
	for _, p := range existing {
		existingMap[p.NotificationType] = p
	}

	items := make([]*notificationCs.NotificationPreferenceItem, 0, len(DefaultNotificationTypes))
	for _, nType := range DefaultNotificationTypes {
		if p, ok := existingMap[nType]; ok {
			items = append(items, &notificationCs.NotificationPreferenceItem{
				ID:               p.ID,
				NotificationType: p.NotificationType,
				Enabled:          p.Enabled,
			})
		} else {
			items = append(items, &notificationCs.NotificationPreferenceItem{
				NotificationType: nType,
				Enabled:          1,
			})
		}
	}

	return &notificationCs.GetNotificationPreferencesResponse{
		Preferences: items,
	}, nil
}

func (b *preferenceBusiness) UpdatePreference(uid int64, req *notificationCs.UpdateNotificationPreferenceRequest) error {
	if uid == 0 {
		return errors.New("请先登录")
	}
	if req.NotificationType == "" {
		return errors.New("通知类型不能为空")
	}

	pref := &notificationModel.TbNotificationPreference{
		Uid:              uid,
		NotificationType: req.NotificationType,
		Enabled:          req.Enabled,
	}

	return notificationModel.TbNotificationPreferenceModel.Upsert(pref)
}

func (b *preferenceBusiness) BatchUpdatePreferences(uid int64, req *notificationCs.BatchUpdateNotificationPreferencesRequest) (*notificationCs.BatchUpdateNotificationPreferencesResponse, error) {
	if uid == 0 {
		return nil, errors.New("请先登录")
	}
	if len(req.Preferences) == 0 {
		return nil, errors.New("请选择要更新的偏好")
	}

	preferences := make([]*notificationModel.TbNotificationPreference, 0, len(req.Preferences))
	for _, p := range req.Preferences {
		preferences = append(preferences, &notificationModel.TbNotificationPreference{
			Uid:              uid,
			NotificationType: p.NotificationType,
			Enabled:          p.Enabled,
		})
	}

	err := notificationModel.TbNotificationPreferenceModel.BatchUpsert(preferences)
	if err != nil {
		return nil, errors.New("批量更新偏好失败: " + err.Error())
	}

	prefs, err := b.GetPreferences(uid)
	if err != nil {
		return nil, err
	}

	return &notificationCs.BatchUpdateNotificationPreferencesResponse{
		Preferences: prefs.Preferences,
	}, nil
}

func (b *preferenceBusiness) GetPreferencesAsResponse(uid int64) (*notificationCs.GetNotificationPreferencesResponse, error) {
	return b.GetPreferences(uid)
}
