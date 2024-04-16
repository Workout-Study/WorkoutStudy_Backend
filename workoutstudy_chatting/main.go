package main

import (
	"context"
	"workoutstudy_chatting/config"
	"workoutstudy_chatting/handler"
	"workoutstudy_chatting/persistence"
	"workoutstudy_chatting/service"

	"github.com/gin-gonic/gin"
)

func main() {
	// persistence 패키지에서 DB 인스턴스를 초기화하고 반환받습니다.
	DB := persistence.InitializeDB()

	// 서비스 인스턴스 생성
	chatService := service.NewChatService(persistence.NewChatRepository(DB))
	fitMateService := service.NewFitMateService(persistence.NewPostgresFitMateRepository(DB))
	fitGroupService := service.NewFitGroupService(persistence.NewFitGroupRepository(DB))

	// Handler 인스턴스 생성
	chatHandler := handler.NewChatHandler(chatService, fitMateService, fitGroupService)
	fitMateHandler := handler.NewFitMateHandler(fitMateService)

	r := gin.Default()
	r.GET("/chat", chatHandler.Chat)
	r.GET("/retrieve/fit-group", fitMateHandler.RetrieveFitGroupByMateID)
	r.GET("/retrieve/message", chatHandler.RetrieveMessages)

	// Kafka Consumer 설정 및 실행
	kafkaConsumer := config.NewKafkaConsumer("kafka-1:9092", "chatting-server-consumer", []string{"fit-mate", "fit-group"})
	ctx := context.Background()
	go func() {
		kafkaConsumer.Consume(ctx) // 토픽 구독, 토픽 목록을 여기서는 명시하지 않습니다.
	}()
	r.Run(":8888")
}
