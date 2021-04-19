package service

import (
	"encoding/json"
	"fmt"
	"go-websocket-cluster/entity"
	"strconv"
	"sync"
)

type HubService struct {
	clients map[*entity.Client]bool
	broadcast chan []byte
	connected chan *entity.Client
	disconnected chan *entity.Client
	redisService *RedisService
}

func NewHubService(redisService *RedisService) *HubService {
	return &HubService{
		broadcast: make(chan []byte),
		connected: make(chan *entity.Client),
		disconnected: make(chan *entity.Client),
		clients: make(map[*entity.Client]bool),
		redisService: redisService,
	}
}


func (hs *HubService) Run()  {
	mutex := &sync.Mutex{}
	for {
		select {
		case client := <-hs.connected:
			mutex.Lock()
			hs.clients[client] = true
			mutex.Unlock()
		case client := <-hs.disconnected:
			if _, ok :=hs.clients[client]; ok {
				mutex.Lock()
				delete(hs.clients, client)
				mutex.Unlock()
				close(client.Send)

			}
		case message := <-hs.broadcast:
			for client := range hs.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					mutex.Lock()
					delete(hs.clients, client)
					mutex.Unlock()
				}
			}
		}
	}
}

// 订阅消息
func (hs *HubService) SubscribeMessage(channel string) {
	sub := hs.redisService.client.Subscribe(hs.redisService.ctx, channel)
	ch := sub.Channel()
	for msg := range ch {
		hs.broadcast <- []byte(msg.Payload)
		fmt.Println(msg.Channel, msg.Payload)
	}
}

// 发布消息
func (hs *HubService) PublishMessage(channel string, message *entity.JsonMessage) {
	str, _ := json.Marshal(message)
	err := hs.redisService.client.Publish(hs.redisService.ctx, channel, str).Err()
	if err != nil {
		panic(err)
	}
}

// 增加在线人数
func (hs *HubService) AddOnlineTotal()  {
	hs.redisService.client.Incr(hs.redisService.ctx, entity.CacheKeyOnlineTotal)
}

// 减少在线人数
func (hs *HubService) SubOnlineTotal()  {
	hs.redisService.client.Decr(hs.redisService.ctx, entity.CacheKeyOnlineTotal)
}

// 获取在线人数
func (hs *HubService) GetOnlineTotal() int64 {
	val, err := hs.redisService.client.Get(hs.redisService.ctx, entity.CacheKeyOnlineTotal).Result()
	if err != nil {
		return 0
	}
	count, _:= strconv.ParseInt(val,10, 64)
	return count
}

// 增加点赞人数
func (hs *HubService) AddLikedCount()  {
	hs.redisService.client.Incr(hs.redisService.ctx, entity.CacheKeyLikedCount)
}

// 获取点赞人数
func (hs *HubService) GetLikedCount() int64 {
	val, err := hs.redisService.client.Get(hs.redisService.ctx, entity.CacheKeyLikedCount).Result()
	if err != nil {
		return 0
	}
	count, _:= strconv.ParseInt(val,10, 64)
	return count
}