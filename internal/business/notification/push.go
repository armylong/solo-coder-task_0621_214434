package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/armylong/armylong-go/internal/common/webcache"
	notificationCs "github.com/armylong/armylong-go/internal/cs/notification"
)

type notificationPusher struct {
	mu   sync.RWMutex
	hubs map[int64]chan *notificationCs.PushEvent
}

var NotificationPusher = &notificationPusher{
	hubs: make(map[int64]chan *notificationCs.PushEvent),
}

func (p *notificationPusher) Subscribe(uid int64) chan *notificationCs.PushEvent {
	p.mu.Lock()
	defer p.mu.Unlock()

	if ch, exists := p.hubs[uid]; exists {
		return ch
	}

	ch := make(chan *notificationCs.PushEvent, 64)
	p.hubs[uid] = ch
	return ch
}

func (p *notificationPusher) Unsubscribe(uid int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if ch, exists := p.hubs[uid]; exists {
		close(ch)
		delete(p.hubs, uid)
	}
}

func (p *notificationPusher) Push(uids []int64, event *notificationCs.PushEvent) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for _, uid := range uids {
		if ch, exists := p.hubs[uid]; exists {
			select {
			case ch <- event:
			default:
			}
		}
	}

	p.pushViaRedis(uids, event)
}

func (p *notificationPusher) pushViaRedis(uids []int64, event *notificationCs.PushEvent) {
	if webcache.RedisClient == nil {
		return
	}

	data, err := json.Marshal(event)
	if err != nil {
		return
	}

	for _, uid := range uids {
		channel := fmt.Sprintf("notification:%d", uid)
		webcache.RedisClient.Publish(context.Background(), channel, string(data))
	}
}
