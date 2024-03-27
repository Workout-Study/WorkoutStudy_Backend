package main

import (
	"net/http"
	"workoutstudy_chatting/controller"
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

func main() {
	persistence.InitializeDB()

	r := gin.Default()
	r.GET("/ws/:fit_group", controller.ChatHandler)
	r.Run(":8080")
}
