package entity

type JsonMessage struct {
	Type int `json:"type"`
	Data interface{} `json:"data,omitempty"`
	TimeStamp int64 `json:"timestamp,omitempty"`
}

const (
	// 在线人数
	MessageTypeOnlineTotal = 1
	// 点赞人数
	MessageTypeLikedTotal = 2
	// 获取在线用户
	MessageTypeGetOnlineTotal = 3
	// 获取点赞人数
	MessageTypeGetLikedTotal = 4
	// 文本信息
	MessageTypeText = 5
	CacheKeyOnlineTotal = "online-total"
	CacheKeyLikedCount = "liked-count"
)