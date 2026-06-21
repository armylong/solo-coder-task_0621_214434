package notification

type GetNotificationPreferencesRequest struct{}

type NotificationPreferenceItem struct {
	ID               int64  `json:"id"`
	NotificationType string `json:"notification_type"`
	Enabled          int    `json:"enabled"`
}

type GetNotificationPreferencesResponse struct {
	Preferences []*NotificationPreferenceItem `json:"preferences"`
}

type UpdateNotificationPreferenceRequest struct {
	NotificationType string `json:"notification_type" form:"notification_type"`
	Enabled          int    `json:"enabled" form:"enabled"`
}

type BatchUpdateNotificationPreferencesRequest struct {
	Preferences []*UpdateNotificationPreferenceRequest `json:"preferences"`
}

type BatchUpdateNotificationPreferencesResponse struct {
	Preferences []*NotificationPreferenceItem `json:"preferences"`
}
