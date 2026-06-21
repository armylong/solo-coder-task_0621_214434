package notification

import (
	commonWs "github.com/armylong/armylong-go/internal/common/websocket"
	notificationCs "github.com/armylong/armylong-go/internal/cs/notification"
	libWs "github.com/armylong/go-library/service/websocket"
)

func PushNotification(uids []int64, event *notificationCs.PushEvent) {
	if commonWs.Manager == nil {
		return
	}
	msg := libWs.NewMessage(event.Type, event.Data)
	commonWs.Manager.PushToUsers(uids, msg)
}
