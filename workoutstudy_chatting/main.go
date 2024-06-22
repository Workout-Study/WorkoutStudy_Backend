package main

import (
	"context"
	"workoutstudy_chatting/config"
	"workoutstudy_chatting/handler"
	"workoutstudy_chatting/persistence"
	"workoutstudy_chatting/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// persistence 패키지에서 DB 인스턴스를 초기화하고 반환받습니다.
	DB := persistence.InitializeDB()

	// 서비스 인스턴스 생성
	chatService := service.NewChatService(persistence.NewChatRepository(DB))
	fitMateService := service.NewFitMateService(persistence.NewPostgresFitMateRepository(DB), make(chan int))
	fitGroupService := service.NewFitGroupService(persistence.NewFitGroupRepository(DB), make(chan int))
	userService := service.NewUserService(persistence.NewUserRepository(DB)) // 사용자 서비스 추가

	// Handler 인스턴스 생성
	chatHandler := handler.NewChatHandler(chatService, fitMateService, fitGroupService)
	fitMateHandler := handler.NewFitMateHandler(fitMateService)

	r := gin.Default()

	// 정적 파일 제공 설정
	r.Static("/docs", "./docs")
	// Swagger 라우트 설정
	r.GET("/swagger/index.html", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/doc.json")))

	r.GET("/chat", chatHandler.Chat)
	r.GET("/retrieve/fit-group", fitMateHandler.RetrieveFitGroupByUserID)
	r.GET("/retrieve/message", chatHandler.RetrieveMessages)

	// Kafka Consumer 설정 및 실행
	kafkaConsumer := config.NewKafkaConsumer("kafka-1:9092", "chatting-server-consumer", []string{"fit-mate", "fit-group", "user-create-event", "user-info"})
	// context 생성
	// Go 에서 요청 간의 데이터, 취소, 신호, 데드라인 등을 전달하는 방법을 제공함
	// 주로 네트워크 요청, 서버 핸들러, 백그라운드 작업을 제어하고 취소하는 데 사용
	ctx := context.Background()

	msgChan := make(chan handler.MessageEvent)
	go kafkaConsumer.Consume(ctx, msgChan)

	go handler.HandleMessage(msgChan, fitMateService, fitGroupService, userService)

	r.Run(":8888")
}
