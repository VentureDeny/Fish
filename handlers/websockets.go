// handlers/websocket.go
package handlers

import (
	"fmt"
	"log"
	"net/http"

	"fish/db"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 根据需要限制跨域
	},
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("升级到WebSocket失败: %v", err)
		return
	}
	defer conn.Close()

	log.Println("客户端已连接")

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket连接异常关闭: %v", err)
			}
			break
		}

		log.Printf("接收到消息: %s", message)

		// 将消息存储到数据库
		err = storeMessage(string(message))
		if err != nil {
			log.Printf("存储消息到数据库失败: %v", err)
			errorMsg := fmt.Sprintf("存储消息失败: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte(errorMsg))
			continue
		}

		// 回发消息给客户端（回声）
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			log.Printf("发送消息失败: %v", err)
			break
		}
	}

	log.Println("客户端已断开连接")
}

func storeMessage(content string) error {
	query := `INSERT INTO messages (content) VALUES ($1)`
	_, err := db.DB.Exec(query, content)
	return err
}
