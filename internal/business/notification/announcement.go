package notification

import (
	"errors"
	"strconv"
	"strings"
	"time"

	notificationCs "github.com/armylong/armylong-go/internal/cs/notification"
	notificationModel "github.com/armylong/armylong-go/internal/model/notification"
	userModel "github.com/armylong/armylong-go/internal/model/user"
)

type announcementBusiness struct{}

var AnnouncementBusiness = &announcementBusiness{}

func (b *announcementBusiness) Create(createdBy int64, req *notificationCs.CreateAnnouncementRequest) (*notificationCs.AnnouncementDetailResponse, error) {
	if req.Title == "" {
		return nil, errors.New("标题不能为空")
	}
	if req.Type == "" {
		req.Type = notificationModel.AnnouncementTypeSystem
	}
	if req.TargetAudience == "" {
		req.TargetAudience = notificationModel.AnnouncementTargetAll
	}

	announcement := &notificationModel.TbAnnouncement{
		Title:          req.Title,
		Content:        req.Content,
		Type:           req.Type,
		Priority:       req.Priority,
		TargetAudience: req.TargetAudience,
		TargetUids:     req.TargetUids,
		Status:         notificationModel.AnnouncementStatusDraft,
		CreatedBy:      createdBy,
	}

	id, err := notificationModel.TbAnnouncementModel.Create(announcement)
	if err != nil {
		return nil, errors.New("创建公告失败: " + err.Error())
	}

	announcement.ID = id
	return b.toDetailResponse(announcement), nil
}

func (b *announcementBusiness) Update(req *notificationCs.UpdateAnnouncementRequest) (*notificationCs.AnnouncementDetailResponse, error) {
	if req.ID == 0 {
		return nil, errors.New("公告ID不能为空")
	}

	announcement, err := notificationModel.TbAnnouncementModel.GetByID(req.ID)
	if err != nil || announcement == nil {
		return nil, errors.New("公告不存在")
	}

	if announcement.Status != notificationModel.AnnouncementStatusDraft {
		return nil, errors.New("只能编辑草稿状态的公告")
	}

	if req.Title != "" {
		announcement.Title = req.Title
	}
	if req.Content != "" {
		announcement.Content = req.Content
	}
	if req.Type != "" {
		announcement.Type = req.Type
	}
	announcement.Priority = req.Priority
	if req.TargetAudience != "" {
		announcement.TargetAudience = req.TargetAudience
	}
	announcement.TargetUids = req.TargetUids

	err = notificationModel.TbAnnouncementModel.Update(announcement)
	if err != nil {
		return nil, errors.New("更新公告失败: " + err.Error())
	}

	return b.toDetailResponse(announcement), nil
}

func (b *announcementBusiness) Delete(id int64) error {
	if id == 0 {
		return errors.New("公告ID不能为空")
	}

	announcement, err := notificationModel.TbAnnouncementModel.GetByID(id)
	if err != nil || announcement == nil {
		return errors.New("公告不存在")
	}

	return notificationModel.TbAnnouncementModel.Delete(id)
}

func (b *announcementBusiness) Publish(id int64) (*notificationCs.AnnouncementDetailResponse, error) {
	if id == 0 {
		return nil, errors.New("公告ID不能为空")
	}

	announcement, err := notificationModel.TbAnnouncementModel.GetByID(id)
	if err != nil || announcement == nil {
		return nil, errors.New("公告不存在")
	}

	if announcement.Status != notificationModel.AnnouncementStatusDraft {
		return nil, errors.New("只能发布草稿状态的公告")
	}

	now := time.Now()
	announcement.Status = notificationModel.AnnouncementStatusPublished
	announcement.PublishedAt = &now

	err = notificationModel.TbAnnouncementModel.Update(announcement)
	if err != nil {
		return nil, errors.New("发布公告失败: " + err.Error())
	}

	go b.createNotificationsAndPush(announcement)

	return b.toDetailResponse(announcement), nil
}

func (b *announcementBusiness) Archive(id int64) (*notificationCs.AnnouncementDetailResponse, error) {
	if id == 0 {
		return nil, errors.New("公告ID不能为空")
	}

	announcement, err := notificationModel.TbAnnouncementModel.GetByID(id)
	if err != nil || announcement == nil {
		return nil, errors.New("公告不存在")
	}

	if announcement.Status != notificationModel.AnnouncementStatusPublished {
		return nil, errors.New("只能下架已发布的公告")
	}

	announcement.Status = notificationModel.AnnouncementStatusArchived
	err = notificationModel.TbAnnouncementModel.Update(announcement)
	if err != nil {
		return nil, errors.New("下架公告失败: " + err.Error())
	}

	return b.toDetailResponse(announcement), nil
}

func (b *announcementBusiness) List(req *notificationCs.ListAnnouncementsRequest) (*notificationCs.ListAnnouncementsResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	list, total, err := notificationModel.TbAnnouncementModel.List(req.Status, req.Page, req.PageSize)
	if err != nil {
		return nil, errors.New("获取公告列表失败: " + err.Error())
	}

	items := make([]*notificationCs.AnnouncementItem, 0, len(list))
	for _, a := range list {
		items = append(items, b.toItem(a))
	}

	return &notificationCs.ListAnnouncementsResponse{
		List:     items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (b *announcementBusiness) GetByID(id int64) (*notificationCs.AnnouncementDetailResponse, error) {
	announcement, err := notificationModel.TbAnnouncementModel.GetByID(id)
	if err != nil || announcement == nil {
		return nil, errors.New("公告不存在")
	}
	return b.toDetailResponse(announcement), nil
}

func (b *announcementBusiness) createNotificationsAndPush(announcement *notificationModel.TbAnnouncement) {
	targetUids := b.resolveTargetUids(announcement)
	if len(targetUids) == 0 {
		return
	}

	notifications := make([]*notificationModel.TbNotification, 0, len(targetUids))
	pushUids := make([]int64, 0, len(targetUids))

	for _, uid := range targetUids {
		if !notificationModel.TbNotificationPreferenceModel.IsEnabled(uid, announcement.Type) {
			continue
		}

		exists, _ := notificationModel.TbNotificationModel.ExistsByUidAndAnnouncementID(uid, announcement.ID)
		if exists {
			continue
		}

		notifications = append(notifications, &notificationModel.TbNotification{
			Uid:            uid,
			AnnouncementID: announcement.ID,
			Title:          announcement.Title,
			Content:        announcement.Content,
			Type:           announcement.Type,
			IsRead:         0,
		})
		pushUids = append(pushUids, uid)
	}

	if len(notifications) > 0 {
		_ = notificationModel.TbNotificationModel.CreateBatch(notifications)
	}

	if len(pushUids) > 0 {
		pushEvent := &notificationCs.PushEvent{
			Type: "new_notification",
			Data: map[string]interface{}{
				"announcement_id": announcement.ID,
				"title":           announcement.Title,
				"type":            announcement.Type,
				"priority":        announcement.Priority,
			},
		}
		NotificationPusher.Push(pushUids, pushEvent)
	}
}

func (b *announcementBusiness) resolveTargetUids(announcement *notificationModel.TbAnnouncement) []int64 {
	switch announcement.TargetAudience {
	case notificationModel.AnnouncementTargetAll:
		users, err := userModel.TbUserModel.List(10000, 0)
		if err != nil {
			return nil
		}
		uids := make([]int64, 0, len(users))
		for _, u := range users {
			if u.Status == 1 {
				uids = append(uids, u.Uid)
			}
		}
		return uids

	case notificationModel.AnnouncementTargetAdmin:
		admins, err := userModel.TbAdminUserModel.ListAll()
		if err != nil {
			return nil
		}
		uids := make([]int64, 0, len(admins))
		for _, a := range admins {
			uids = append(uids, a.Uid)
		}
		return uids

	case notificationModel.AnnouncementTargetSpecific:
		if announcement.TargetUids == "" {
			return nil
		}
		parts := strings.Split(announcement.TargetUids, ",")
		uids := make([]int64, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if uid, err := strconv.ParseInt(p, 10, 64); err == nil {
				uids = append(uids, uid)
			}
		}
		return uids

	default:
		return nil
	}
}

func (b *announcementBusiness) toItem(a *notificationModel.TbAnnouncement) *notificationCs.AnnouncementItem {
	return &notificationCs.AnnouncementItem{
		ID:             a.ID,
		Title:          a.Title,
		Content:        a.Content,
		Type:           a.Type,
		Priority:       a.Priority,
		TargetAudience: a.TargetAudience,
		TargetUids:     a.TargetUids,
		Status:         a.Status,
		CreatedBy:      a.CreatedBy,
		PublishedAt:    a.PublishedAt,
		CreatedAt:      a.CreatedAt,
		UpdatedAt:      a.UpdatedAt,
	}
}

func (b *announcementBusiness) toDetailResponse(a *notificationModel.TbAnnouncement) *notificationCs.AnnouncementDetailResponse {
	return &notificationCs.AnnouncementDetailResponse{
		AnnouncementItem: b.toItem(a),
	}
}
