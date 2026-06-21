package notification

import (
	"errors"

	notificationCs "github.com/armylong/armylong-go/internal/cs/notification"
	notificationModel "github.com/armylong/armylong-go/internal/model/notification"
)

type notificationBusiness struct{}

var NotificationBusiness = &notificationBusiness{}

func (b *notificationBusiness) List(uid int64, req *notificationCs.ListNotificationsRequest) (*notificationCs.ListNotificationsResponse, error) {
	if uid == 0 {
		return nil, errors.New("请先登录")
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	list, total, err := notificationModel.TbNotificationModel.ListByUid(uid, req.Type, req.Page, req.PageSize)
	if err != nil {
		return nil, errors.New("获取通知列表失败: " + err.Error())
	}

	items := make([]*notificationCs.NotificationItem, 0, len(list))
	for _, n := range list {
		items = append(items, b.toItem(n))
	}

	return &notificationCs.ListNotificationsResponse{
		List:     items,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

func (b *notificationBusiness) GetDetail(uid int64, id int64) (*notificationCs.NotificationDetailResponse, error) {
	if uid == 0 {
		return nil, errors.New("请先登录")
	}

	notification, err := notificationModel.TbNotificationModel.GetByID(id)
	if err != nil || notification == nil {
		return nil, errors.New("通知不存在")
	}

	if notification.Uid != uid {
		return nil, errors.New("无权访问此通知")
	}

	if notification.IsRead == 0 {
		_ = notificationModel.TbNotificationModel.MarkAsRead(id)
		notification.IsRead = 1
	}

	return &notificationCs.NotificationDetailResponse{
		NotificationItem: b.toItem(notification),
	}, nil
}

func (b *notificationBusiness) MarkAsRead(uid int64, ids []int64) error {
	if uid == 0 {
		return errors.New("请先登录")
	}
	if len(ids) == 0 {
		return errors.New("请选择要标记的通知")
	}
	return notificationModel.TbNotificationModel.MarkAsReadByUid(uid, ids)
}

func (b *notificationBusiness) MarkAllAsRead(uid int64) error {
	if uid == 0 {
		return errors.New("请先登录")
	}
	return notificationModel.TbNotificationModel.MarkAllAsRead(uid)
}

func (b *notificationBusiness) DeleteRead(uid int64) error {
	if uid == 0 {
		return errors.New("请先登录")
	}
	return notificationModel.TbNotificationModel.DeleteReadByUid(uid)
}

func (b *notificationBusiness) Delete(uid int64, ids []int64) error {
	if uid == 0 {
		return errors.New("请先登录")
	}
	if len(ids) == 0 {
		return errors.New("请选择要删除的通知")
	}
	return notificationModel.TbNotificationModel.DeleteByUid(uid, ids)
}

func (b *notificationBusiness) UnreadCount(uid int64) (*notificationCs.UnreadCountResponse, error) {
	if uid == 0 {
		return nil, errors.New("请先登录")
	}

	total, err := notificationModel.TbNotificationModel.CountUnreadByUid(uid)
	if err != nil {
		return nil, errors.New("获取未读数失败: " + err.Error())
	}

	byTypeRows, err := notificationModel.TbNotificationModel.CountUnreadByUidGroupByType(uid)
	if err != nil {
		return nil, errors.New("获取未读数分组失败: " + err.Error())
	}

	byType := make([]*notificationCs.UnreadCountByType, 0, len(byTypeRows))
	for _, row := range byTypeRows {
		nType, _ := row["type"].(string)
		count, _ := row["count"].(int64)
		byType = append(byType, &notificationCs.UnreadCountByType{
			Type:  nType,
			Count: count,
		})
	}

	return &notificationCs.UnreadCountResponse{
		Total:  total,
		ByType: byType,
	}, nil
}

func (b *notificationBusiness) toItem(n *notificationModel.TbNotification) *notificationCs.NotificationItem {
	return &notificationCs.NotificationItem{
		ID:             n.ID,
		AnnouncementID: n.AnnouncementID,
		Title:          n.Title,
		Content:        n.Content,
		Type:           n.Type,
		IsRead:         n.IsRead,
		ReadAt:         n.ReadAt,
		CreatedAt:      n.CreatedAt,
	}
}
