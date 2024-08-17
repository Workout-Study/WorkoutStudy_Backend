package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"workoutstudy_chatting/config"
	"workoutstudy_chatting/handler"
	"workoutstudy_chatting/persistence"
	"workoutstudy_chatting/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "workoutstudy_chatting/docs" // Swagger docs
)

func main() {
	DB := persistence.InitializeDB()

	chatService := service.NewChatService(persistence.NewChatRepository(DB))
	fitMateService := service.NewFitMateService(persistence.NewPostgresFitMateRepository(DB), make(chan int))
	fitGroupService := service.NewFitGroupService(persistence.NewFitGroupRepository(DB), make(chan int))
	userService := service.NewUserService(persistence.NewUserRepository(DB))

	testHandler := handler.NewTestHandler(userService, fitGroupService, fitMateService)
	chatHandler := handler.NewChatHandler(chatService, fitMateService, fitGroupService)
	fitMateHandler := handler.NewFitMateHandler(fitMateService)

	r := gin.Default()
	r.Static("/docs", "./docs")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/doc.json")))

	r.POST("/test/create/user", testHandler.CreateUser)
	r.POST("/test/create/fit-group", testHandler.CreateFitGroup)
	r.POST("/test/create/fit-mate", testHandler.CreateFitMate)

	r.GET("/chat", chatHandler.Chat)
	r.GET("/retrieve/fit-group", fitMateHandler.RetrieveFitGroupByUserID)
	r.GET("/retrieve/message", chatHandler.RetrieveMessages)

	msgChan := make(chan handler.MessageEvent)

	kafkaConsumer := config.NewKafkaConsumer([]string{"kafka-1:9092" /*, "kafka-2:9093", "kafka-3:9094"*/}, "chatting-service", []string{"fit-mate", "fit-group", "user-create-event", "user-info-event"})

	ctx, cancel := context.WithCancel(context.Background())
	log.Println("Context created for Kafka consumer")

	go kafkaConsumer.Consume(ctx, msgChan)

	go handler.HandleMessage(msgChan, fitMateService, fitGroupService, userService)
	// Graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Printf("Received signal: %s, initiating shutdown", sig)
		cancel()
	}()

	log.Println("Kafka Consumer and Handlers started")
	r.Run(":8888")
}
