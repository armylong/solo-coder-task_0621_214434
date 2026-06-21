package websocket

import (
	"github.com/armylong/armylong-go/internal/common/webcache"
	libWs "github.com/armylong/go-library/service/websocket"
	"github.com/redis/go-redis/v9"
)

var Manager *libWs.ConnManager

func init() {
	var rdb *redis.Client
	if webcache.RedisClient != nil {
		rdb = webcache.RedisClient.Client
	}
	Manager = libWs.NewConnManager(rdb, "armylong:ws", nil)
}
