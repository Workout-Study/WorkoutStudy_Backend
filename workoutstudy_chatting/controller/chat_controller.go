package controller

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 채팅방별 클라이언트 관리를 위한 맵과 락
var rooms = make(map[string]map[*websocket.Conn]bool)
var roomLock = sync.Mutex{}

func ChatHandler(c *gin.Context) {
	// 기존 wshandler 함수 내용
	fitGroup := c.Param("fit_group")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("웹 소켓 업그레이드 실패:", err)
		return
	}
	defer conn.Close()

	roomLock.Lock()
	if rooms[fitGroup] == nil {
		rooms[fitGroup] = make(map[*websocket.Conn]bool)
	}
	rooms[fitGroup][conn] = true
	roomLock.Unlock()

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("메시지 읽기 실패:", err)
			break
		}
		// 메시지를 동일한 fitGroup의 모든 클라이언트에게 브로드캐스트
		for client := range rooms[fitGroup] {
			if err := client.WriteMessage(mt, message); err != nil {
				fmt.Println("메시지 전송 실패:", err)
				client.Close()
				delete(rooms[fitGroup], client)
			}
		}
	}
}
