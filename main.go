package main

import (
	"mafia/route"
	"mafia/services"
	"mafia/storage"
	"mafia/websocket/chat"
	"mafia/websocket/timer"

	"github.com/gin-gonic/gin"
)

func main() {
	// Инициализация конфигурации
	services.InitConfig()


	// Инициализация базы данных
  storage.Start()

	// Создание записи в таблице UserAccount
	db := storage.GetDatabase()
	defer db.Close()


	// Создание роутов
	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Next()
	})
	route.RunAuthRoutes(r)
	route.RunVerifyRoutes(r)
	route.RunUserRoutes(r)
	
	
	chat.RunWSRoutes(r)
	timer.RunGameRoutes(r)
	r.Run(":8082")
	
}