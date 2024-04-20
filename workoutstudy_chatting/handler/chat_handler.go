package handler

import (
	"encoding/json"
	"errors"
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
	ChatService     *service.ChatService
	FitMateService  service.FitMateService // FitMateService 추가
	FitGroupService *service.FitGroupService
}

func NewChatHandler(chatService *service.ChatService, fitMateService service.FitMateService, fitGroupService *service.FitGroupService) *ChatHandler {
	return &ChatHandler{
		ChatService:     chatService,
		FitMateService:  fitMateService,
		FitGroupService: fitGroupService, // FitMateService 초기화
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Room struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan model.ChatMessage
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
}

func NewRoom() *Room {
	return &Room{
		broadcast:  make(chan model.ChatMessage),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		clients:    make(map[*websocket.Conn]bool),
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.register:
			// 채팅방 연결시 사용자를 클라이언트로 등록
			r.clients[client] = true
		case client := <-r.unregister:
			// 클라이언트 등록 해제
			if _, ok := r.clients[client]; ok {
				delete(r.clients, client)
				client.Close() // 여기서 웹소켓 연결 종료
			}
		case message := <-r.broadcast:
			// 채팅방에 연결(웹 소켓으로 통신하고 있는)되어 있는 모든 사용자들에게 메시지 브로드캐스트
			for client := range r.clients {
				err := client.WriteJSON(message)
				if err != nil {
					// 에러 발생 시 클라이언트 해제 처리
					log.Printf("error: %v", err)
					client.Close()
					delete(r.clients, client)
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

func (h *ChatHandler) Chat(c *gin.Context) {
	fitGroupIDStr := c.Query("fitGroupId")
	// fitMateIDStr := c.Query("fitMateId")
	/*
		TODO : 위 처럼 파라미터가 들어왔을 때
		1. 해당 fitMate가 존재하는지 검증
		1-1. 없으면 에러 메시지와 함께 웹소켓 연결 거부
		2. fitGroupId 로 ftiGroup 존재하는지 검증
		2-1. fitGroup 이 존재하지 않는다면 에러 메시지와 함께 웹 소켓 연결 거부
		3. fitMate가 존재한다면 해당 fitMate가 fitGroup에 속해있는지 검증
		3-1. fitMate가 fitGroup에 속해 있지 않다면 에러 메시지와 함께 웹소켓 연결 거부
		4. 올바른 사용자라면 romm을 통해 연결
		5. DB 테이블은 fit_group_mate 사용
	*/

	// FitMate 조회
	// fitMate, err := h.FitMateService.GetFitMateByID(fitMateIDStr)
	// if err != nil {
	// 	// 에러 처리: 조회 중 에러 발생
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "서버 내부 오류"})
	// 	return
	// }
	// if fitMate == nil {
	// 	// 에러 처리: FitMate가 존재하지 않음. 여기서 WebSocket 연결을 거부합니다.
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "해당 FitMate가 존재하지 않습니다"})
	// 	return
	// }

	roomLock.Lock()
	room, ok := rooms[fitGroupIDStr]
	if !ok {
		room = NewRoom()
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

	room.register <- conn

	// 클라이언트로부터 메시지를 읽고 room의 broadcast 채널에 전달하는 로직...
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read error: %v", err)
			break
		}

		var chatMsg model.ChatMessage
		log.Printf(chatMsg.Message)
		if err := json.Unmarshal(message, &chatMsg); err != nil {
			log.Printf("unmarshal error: %v", err)
			continue
		}

		// 채팅방에 채팅 브로드 캐스팅
		room.broadcast <- chatMsg

		// 데이터베이스에 메시지 저장 로직은 여기에 포함...
		err = h.ChatService.SaveChatMessage(chatMsg)
		if err != nil {
			log.Printf("메시지 저장 실패: %v", err)
			// 메시지 저장 실패 시 클라이언트에게 실패 메시지 전송
			failMsg := model.ChatMessage{Message: "메시지 저장에 실패했습니다."}
			failMsgJSON, _ := json.Marshal(failMsg)
			if writeErr := conn.WriteMessage(websocket.TextMessage, failMsgJSON); writeErr != nil {
				log.Printf("클라이언트에게 실패 메시지 전송 실패: %v", writeErr)
				conn.Close()
				return
			}
			continue
		}
	}
	room.unregister <- conn
}

func (h *ChatHandler) RetrieveMessages(c *gin.Context) {
	messageID := c.Query("messageId")
	fitGroupIDStr := c.Query("fitGroupId")
	fitMateID := c.Query("fitMateId")
	messageTimeStr := c.Query("messageTime")
	messageType := c.Query("messageType")

	log.Printf("Received messageId: %s", messageID)
	log.Printf("Received fitGroupId: %s", fitGroupIDStr)
	log.Printf("Received fitMateId: %s", fitMateID)
	log.Printf("Received messageTime: %s", messageTimeStr)
	log.Printf("Received messageType: %s", messageType)

	fitGroupID, err := strconv.Atoi(fitGroupIDStr)
	if err != nil {
		// TODO : 에러는 소문자로
		// c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fit-group-id"})
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.New("error: invalid fit-group-id"))
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
