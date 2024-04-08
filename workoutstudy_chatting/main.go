package main

import (
	"workoutstudy_chatting/handler"
	"workoutstudy_chatting/persistence"
	"workoutstudy_chatting/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// persistence 패키지에서 DB 인스턴스를 초기화하고 반환받습니다.
	DB := persistence.InitializeDB()

	// ChatService와 FitMateService 인스턴스 생성
	// ChatService, FitMateService, FitGroupService 인스턴스 생성
	chatService := service.NewChatService(persistence.NewChatRepository(DB))
	fitMateService := service.NewFitMateService(persistence.NewPostgresFitMateRepository(DB))
	fitGroupService := service.NewFitGroupService(persistence.NewFitGroupRepository(DB)) // 예시로 추가

	// ChatHandler와 FitMateHandler 인스턴스 생성, 서비스 인터페이스 주입
	chatHandler := handler.NewChatHandler(chatService, fitMateService, fitGroupService)
	fitMateHandler := handler.NewFitMateHandler(fitMateService)

	r := gin.Default()

	r.GET("/chat", chatHandler.Chat)
	r.GET("/retrieve/fit-group", fitMateHandler.RetrieveFitGroupByMateID)
	r.GET("/retrieve/message", chatHandler.RetrieveMessages)
	r.Run(":8888")
}
