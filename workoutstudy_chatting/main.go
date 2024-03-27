package main

import (
	"fmt"
	"net/http"
	"sync"
	"workoutstudy_chatting/persistence"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 모든 원점(Origin)으로부터의 연결을 허용
		return true
	},
}

// 채팅방별 클라이언트 관리를 위한 맵과 락
var rooms = make(map[string]map[*websocket.Conn]bool)
var roomLock = sync.Mutex{}

func main() {

	// DB 초기화
	persistence.InitializeDB()

	r := gin.Default()

	r.GET("/ws", func(c *gin.Context) {
		wshandler(c.Writer, c.Request)
	})

	r.Run(":8080")
}

func wshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("웹 소켓 업그레이드 실패:", err)
		return
	}
	defer conn.Close()

	for {
		// 클라이언트로부터 메시지 읽기
		mt, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("메시지 읽기 실패:", err)
			break
		}
		fmt.Printf("메시지 받음: %s\n", message)

		// 받은 메시지를 클라이언트에게 다시 전송
		if err := conn.WriteMessage(mt, message); err != nil {
			fmt.Println("메시지 전송 실패:", err)
			break
		}
	}
}
