package notification

import "time"

type CreateAnnouncementRequest struct {
	Title          string `json:"title" form:"title"`
	Content        string `json:"content" form:"content"`
	Type           string `json:"type" form:"type"`
	Priority       int    `json:"priority" form:"priority"`
	TargetAudience string `json:"target_audience" form:"target_audience"`
	TargetUids     string `json:"target_uids" form:"target_uids"`
}

type UpdateAnnouncementRequest struct {
	ID             int64  `json:"id" form:"id"`
	Title          string `json:"title" form:"title"`
	Content        string `json:"content" form:"content"`
	Type           string `json:"type" form:"type"`
	Priority       int    `json:"priority" form:"priority"`
	TargetAudience string `json:"target_audience" form:"target_audience"`
	TargetUids     string `json:"target_uids" form:"target_uids"`
}

type DeleteAnnouncementRequest struct {
	ID int64 `json:"id" form:"id"`
}

type PublishAnnouncementRequest struct {
	ID int64 `json:"id" form:"id"`
}

type ArchiveAnnouncementRequest struct {
	ID int64 `json:"id" form:"id"`
}

type ListAnnouncementsRequest struct {
	Status   int `json:"status" form:"status"`
	Page     int `json:"page" form:"page"`
	PageSize int `json:"page_size" form:"page_size"`
}

type AnnouncementItem struct {
	ID             int64      `json:"id"`
	Title          string     `json:"title"`
	Content        string     `json:"content"`
	Type           string     `json:"type"`
	Priority       int        `json:"priority"`
	TargetAudience string     `json:"target_audience"`
	TargetUids     string     `json:"target_uids"`
	Status         int        `json:"status"`
	CreatedBy      int64      `json:"created_by"`
	PublishedAt    *time.Time `json:"published_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type ListAnnouncementsResponse struct {
	List     []*AnnouncementItem `json:"list"`
	Total    int64               `json:"total"`
	Page     int                 `json:"page"`
	PageSize int                 `json:"page_size"`
}

type AnnouncementDetailResponse struct {
	*AnnouncementItem
}
