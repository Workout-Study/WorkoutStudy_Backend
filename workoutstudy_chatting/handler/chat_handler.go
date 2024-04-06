package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"workoutstudy_chatting/model" // model 패키지 import 추가
	"workoutstudy_chatting/service"
	"workoutstudy_chatting/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	ChatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{ChatService: chatService}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 채팅방별 클라이언트 관리를 위한 맵과 락
var rooms = make(map[string]map[*websocket.Conn]bool)
var roomLock = sync.Mutex{}

func (h *ChatHandler) Chat(c *gin.Context) {
	fitGroup := c.Param("fit-group")
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

		// 메시지를 model 패키지의 ChatMessage 구조체로 언마샬링
		var chatMsg model.ChatMessage
		if err := json.Unmarshal(message, &chatMsg); err != nil {
			fmt.Println("메시지 파싱 실패:", err)
			continue
		}

		// 메시지를 데이터베이스에 저장
		err = h.ChatService.SaveChatMessage(chatMsg)
		if err != nil {
			log.Printf("메시지 저장 실패: %v", err)
			// 메시지 저장 실패시 로그를 기록하고 계속 진행합니다. (또는 적절한 에러 처리를 수행합니다.)
			continue
		}

		fmt.Println("받은 메시지:", chatMsg.Message)

		// 메시지를 모든 클라이언트에게 브로드캐스트
		msgToSend, err := json.Marshal(chatMsg)
		if err != nil {
			fmt.Println("메시지 JSON 변환 실패:", err)
			continue
		}

		for client := range rooms[fitGroup] {
			if err := client.WriteMessage(mt, msgToSend); err != nil {
				fmt.Println("메시지 전송 실패:", err)
				client.Close()
				delete(rooms[fitGroup], client)
			}
		}
	}
}

func (h *ChatHandler) RetrieveMessages(c *gin.Context) {
	messageID := c.Query("message-id")
	fitGroupIDStr := c.Query("fit-group-id")
	fitMateID := c.Query("fit-mate-id")
	messageTimeStr := c.Query("message-time")
	messageType := c.Query("message-type")

	log.Printf("Received message-id: %s", messageID)
	log.Printf("Received fit-group-id: %s", fitGroupIDStr)
	log.Printf("Received fit-mate-id: %s", fitMateID)
	log.Printf("Received message-time: %s", messageTimeStr)
	log.Printf("Received message-type: %s", messageType)

	fitGroupID, err := strconv.Atoi(fitGroupIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fit-group-id"})
		return
	}

	messageTime, err := util.ParseMessageTime(messageTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "시간 파싱 실패"})
		return
	}

	log.Printf("Retrieving messages for fitGroupID: %d, since: %v, messageID: %s", fitGroupID, messageTime, messageID)
	messages, latestMessageId, err := h.ChatService.RetrieveMessages(fitGroupID, messageTime, messageID)
	if err != nil {
		log.Printf("Error retrieving messages: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "채팅 메시지 조회 실패"})
		return
	}
	log.Printf("Retrieved messages: %d, latestMessageId: %s", len(messages), latestMessageId)

	// 조건에 따라 메시지 반환 로직
	if messageID == latestMessageId {
		c.JSON(http.StatusOK, gin.H{"messages": messages[:1]}) // 가장 최신 메시지만 반환
	} else {
		c.JSON(http.StatusOK, gin.H{"messages": messages})
	}
}
