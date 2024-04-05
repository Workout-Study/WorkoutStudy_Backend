package main

import (
	"net/http"
	"workoutstudy_chatting/handler"
	"workoutstudy_chatting/persistence"
	"workoutstudy_chatting/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 모든 원점(Origin)으로부터의 연결을 허용
	},
}

func main() {
	// persistence 패키지에서 DB 인스턴스를 초기화하고 반환받습니다.
	DB := persistence.InitializeDB()

	// ChatService와 FitMateService 인스턴스 생성
	chatService := service.NewChatService(persistence.NewChatRepository(DB))
	fitMateService := service.NewFitMateService(persistence.NewPostgresFitMateRepository(DB))

	// ChatHandler와 FitMateHandler 인스턴스 생성, 서비스 인터페이스 주입
	chatHandler := handler.NewChatHandler(chatService)
	fitMateHandler := handler.NewFitMateHandler(fitMateService)

	r := gin.Default()

	r.GET("/ws/:fit-group", chatHandler.Chat)
	r.GET("/retrieve/fit-group/:fit-mate-id", fitMateHandler.RetrieveFitGroupByMateID)
	r.GET("/retrieve/message", chatHandler.RetrieveMessages)
	r.Run(":8080")
}
