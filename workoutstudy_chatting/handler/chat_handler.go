package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"
	"workoutstudy_chatting/model"
	"workoutstudy_chatting/service"
	"workoutstudy_chatting/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	ChatService     service.ChatUseCase     // 인터페이스 사용
	FitMateService  service.FitMateUseCase  // 인터페이스 사용
	FitGroupService service.FitGroupUseCase // 인터페이스 사용
}

func NewChatHandler(chatService service.ChatUseCase, fitMateService service.FitMateUseCase, fitGroupService service.FitGroupUseCase) *ChatHandler {
	return &ChatHandler{
		ChatService:     chatService,
		FitMateService:  fitMateService,
		FitGroupService: fitGroupService,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn   *websocket.Conn
	userID int
}

type Room struct {
	clients       map[*websocket.Conn]int // 사용자 ID를 저장
	broadcast     chan model.ChatMessage
	register      chan Client
	unregister    chan Client
	fitGroupIDStr string
	activeUsers   map[int]bool // 현재 채팅방에 접속한 사용자 ID를 저장
}

func NewRoom(fitGroupIDStr string) *Room {
	return &Room{
		broadcast:     make(chan model.ChatMessage),
		register:      make(chan Client),
		unregister:    make(chan Client),
		clients:       make(map[*websocket.Conn]int),
		fitGroupIDStr: fitGroupIDStr,
		activeUsers:   make(map[int]bool),
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			r.clients[client.conn] = client.userID
			r.activeUsers[client.userID] = true
		case client := <-r.unregister:
			if userID, ok := r.clients[client.conn]; ok {
				delete(r.clients, client.conn)
				delete(r.activeUsers, userID)
				client.conn.Close()
				if len(r.clients) == 0 {
					roomLock.Lock()
					delete(rooms, r.fitGroupIDStr)
					roomLock.Unlock()
					return
				}
			}
		case message := <-r.broadcast:
			for conn, userID := range r.clients {
				if userID != message.UserID {
					err := conn.WriteJSON(message)
					if err != nil {
						log.Printf("error: %v", err)
						conn.Close()
						delete(r.clients, conn)
						delete(r.activeUsers, userID)
					}
				}
			}
		}
	}
}

// 채팅방별 클라이언트 관리를 위한 맵과 락
var (
	roomLock sync.Mutex
	rooms    = make(map[string]*Room)
)

// @Summary websocket chat
// @Description 실시간 채팅 초기 연결 요청입니다.
// @Description 첫 연결 요청 시 웹소켓 연결이 설정되며, 이후 채팅 메시지는 웹소켓을 통해 전송됩니다.
// @Tags chat
// @Accept json
// @Produce json
// @Param fitGroupId query int true "채팅방 연결을 위한 피트그룹 ID"
// @Param userId query int true "사용자 ID"
// @Success 101 {string} string "WebSocket 연결이 성공적으로 설정되었습니다."
// @Router /chat [get]
func (h *ChatHandler) Chat(c *gin.Context) {
	fitGroupIDStr := c.Query("fitGroupId")

	roomLock.Lock()
	room, ok := rooms[fitGroupIDStr]
	if !ok {
		room = NewRoom(fitGroupIDStr)
		rooms[fitGroupIDStr] = room
		go room.run()
	}
	roomLock.Unlock()

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Websocket upgrade failed:", err)
		return
	}
	defer conn.Close()

	// 클라이언트로부터 사용자 ID를 얻어오는 로직
	userIDStr := c.Query("userId")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		log.Printf("Invalid user ID: %v", err)
		return
	}

	client := Client{conn: conn, userID: userID}
	room.register <- client

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read error: %v", err)
			break
		}

		var chatMsg model.ChatMessage
		if err := json.Unmarshal(message, &chatMsg); err != nil {
			log.Printf("unmarshal error: %v", err)
			continue
		}

		room.broadcast <- chatMsg

		err = h.ChatService.SaveChatMessage(chatMsg)
		if err != nil {
			log.Printf("메시지 저장 실패: %v", err)
			failMsg := model.ChatMessage{Message: "메시지 저장에 실패했습니다."}
			failMsgJSON, _ := json.Marshal(failMsg)
			if writeErr := conn.WriteMessage(websocket.TextMessage, failMsgJSON); writeErr != nil {
				log.Printf("클라이언트에게 실패 메시지 전송 실패: %v", writeErr)
				conn.Close()
				return
			}
			continue
		}

		// 현재 접속해 있지 않은 사용자에게 푸시 알림을 보냅니다.
		roomLock.Lock()
		for id := range room.activeUsers {
			if id != chatMsg.UserID {
				go sendWebhook(chatMsg, id)
			}
		}
		roomLock.Unlock()
	}
	room.unregister <- client
}

func sendWebhook(chatMsg model.ChatMessage, userID int) {
	// TODO : 알림 바디 로그 출력 해야함
	webhookURL := "http://auth-service:8080/chat/real-time-chat"
	jsonData, err := json.Marshal(chatMsg)
	if err != nil {
		log.Printf("웹훅 요청을 위한 JSON 변환 실패: %v", err)
		return
	}

	req, err := http.NewRequest("POST", webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("웹훅 요청 생성 실패: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	_, err = client.Do(req)
	if err != nil {
		log.Printf("웹훅 요청 실패: %v", err)
	}
	// 응답을 무시하고 로그만 남김
}

// @Summary 최신 채팅 내역을 확인하고 동기화 하기 위한 API
// @Description messageId 로 서버측 최신 채팅과 앱의 최신 채팅을 비교
// @Tags message
// @Accept  json
// @Produce  json
// @Param messageId query int true "안드로이드 앱에서 생성된 message UUID"
// @Param fitGroupId query int true "피트그룹 채팅방 ID"
// @Param userId query int true "사용자 ID, auth-server 의 userId, fit-mate-service 의 fitMateUserId"
// @Param messageTime query int true "클라이언트측 메시지 생성 시간"
// @Param messageType query int true "메시지 타입 (CHATTING or TICKET)""
// @Success 200 {array} model.ChatMessage
// @Router /retrieve/message [get]
func (h *ChatHandler) RetrieveMessages(c *gin.Context) {
	messageID := c.Query("messageId")
	fitGroupIDStr := c.Query("fitGroupId")
	userId := c.Query("userId")
	messageTimeStr := c.Query("messageTime")
	messageType := c.Query("messageType")

	// URL 디코딩을 시도합니다.
	decodedMessageTimeStr, err := url.QueryUnescape(messageTimeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message time format"})
		return
	}

	log.Printf("Received messageId: %s", messageID)
	log.Printf("Received fitGroupId: %s", fitGroupIDStr)
	log.Printf("Received userId: %s", userId)
	log.Printf("Received messageTime: %s", messageTimeStr)
	log.Printf("Received messageType: %s", messageType)

	fitGroupID, err := strconv.Atoi(fitGroupIDStr)
	if err != nil {
		// TODO : 에러는 소문자로
		// c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fit-group-id"})
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "잘못된 fit-group-id"})
		return
	}

	// 디코딩된 시간을 파싱합니다.
	messageTime, err := util.ParseMessageTime(decodedMessageTimeStr)
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
