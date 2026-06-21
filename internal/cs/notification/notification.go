package notification

import "time"

type ListNotificationsRequest struct {
	Type     string `json:"type" form:"type"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
}

type NotificationItem struct {
	ID             int64      `json:"id"`
	AnnouncementID int64      `json:"announcement_id"`
	Title          string     `json:"title"`
	Content        string     `json:"content"`
	Type           string     `json:"type"`
	IsRead         int        `json:"is_read"`
	ReadAt         *time.Time `json:"read_at"`
	CreatedAt      time.Time  `json:"created_at"`
}

type ListNotificationsResponse struct {
	List     []*NotificationItem `json:"list"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
}

type GetNotificationDetailRequest struct {
	ID int64 `json:"id" form:"id"`
}

type NotificationDetailResponse struct {
	*NotificationItem
}

type MarkAsReadRequest struct {
	IDs []int64 `json:"ids" form:"ids"`
}

type MarkAllAsReadRequest struct{}

type DeleteReadNotificationsRequest struct{}

type DeleteNotificationsRequest struct {
	IDs []int64 `json:"ids" form:"ids"`
}

type UnreadCountRequest struct{}

type UnreadCountByType struct {
	Type  string `json:"type"`
	Count int64  `json:"count"`
}

type UnreadCountResponse struct {
	Total  int64                `json:"total"`
	ByType []*UnreadCountByType `json:"by_type"`
}

type PushEvent struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
