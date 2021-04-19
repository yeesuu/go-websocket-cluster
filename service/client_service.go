package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"go-websocket-cluster/entity"
	"log"
	"net/http"
	"time"
)

const (
	writeWait = 10 * time.Second
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
	maxMessageSize = 512
)


var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func readPump(c *entity.Client, hubService *HubService, messageChannel string) {
	defer func() {
		hubService.disconnected <- c
		c.Conn.Close()
		go handleSubOnlineTotal(hubService, messageChannel)
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
				go handleSubOnlineTotal(hubService, messageChannel)
			}
			break
		}
		log.Println(string(message))
		jsonMessage := &entity.JsonMessage{}
		err = json.Unmarshal(message, &jsonMessage)
		if err != nil {
			jsonMessage = &entity.JsonMessage{
				Type: entity.MessageTypeText,
				Data: string(message),
				TimeStamp: time.Now().Unix(),
			}
		}
		go handleMessage(hubService, jsonMessage, messageChannel)
	}
}

func writePump(c *entity.Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(<-c.Send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func handleAddOnlineTotal(hubService *HubService, messageChannel string)  {
	hubService.AddOnlineTotal()
	total := hubService.GetOnlineTotal()
	jsonMessage := &entity.JsonMessage{
		Type: entity.MessageTypeOnlineTotal,
		Data: total,
		TimeStamp: time.Now().Unix(),
	}
	hubService.PublishMessage(messageChannel, jsonMessage)
}

func handleSubOnlineTotal(hubService *HubService, messageChannel string)  {
	hubService.SubOnlineTotal()
	total := hubService.GetOnlineTotal()
	jsonMessage := &entity.JsonMessage{
		Type: entity.MessageTypeOnlineTotal,
		Data: total,
		TimeStamp: time.Now().Unix(),
	}
	hubService.PublishMessage(messageChannel, jsonMessage)
}

func handleMessage(hubService *HubService, message *entity.JsonMessage, messageChannel string)  {
	switch message.Type {
	case entity.MessageTypeGetOnlineTotal:
		total := hubService.GetOnlineTotal()
		jsonMessage := &entity.JsonMessage{
			Type: entity.MessageTypeOnlineTotal,
			Data: total,
			TimeStamp: time.Now().Unix(),
		}
		hubService.PublishMessage(messageChannel, jsonMessage)
	case entity.MessageTypeGetLikedTotal:
		total := hubService.GetLikedCount()
		jsonMessage := &entity.JsonMessage{
			Type: entity.MessageTypeLikedTotal,
			Data: total,
			TimeStamp: time.Now().Unix(),
		}
		hubService.PublishMessage(messageChannel, jsonMessage)
	default:
		hubService.PublishMessage(messageChannel, message)
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(hubService *HubService, w http.ResponseWriter, r *http.Request, messageChannel string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &entity.Client{ Conn: conn, Send: make(chan []byte, 256)}
	hubService.connected <- client
	go handleAddOnlineTotal(hubService, messageChannel)
	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go writePump(client)
	go readPump(client, hubService, messageChannel)
}