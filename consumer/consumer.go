package main

import (
	"go_confluence/kafka"
	"go_confluence/websocket"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	go websocket.StartWebSocketServer()

	router.GET("/ws", websocket.HandleWebSocket)

	go router.Run(":8000")

	kafka.StartKafkaConsumer()
}
