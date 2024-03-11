package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var user = make(map[*websocket.Conn]bool)
var channel = make(chan Message)

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {
	r := gin.Default()
	wsupgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	// TOPページ
	r.GET("/", func(ctx *gin.Context) {
		http.ServeFile(ctx.Writer, ctx.Request, "index.html")
	})

	r.GET("/ws", func(ctx *gin.Context) {
		conn, err := wsupgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			log.Printf("error: %+v\n", err)
			return
		}
		user[conn] = true
	})

	go handleMessages()

	port := "8080"
	log.Printf("Server listening on port %s", port)
	r.Run(":" + port)
}

func handleMessages() {
	for {
		message := <-channel
		for client := range user {
			err := client.WriteJSON(message)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(user, client)
			}
		}
	}
}

//参考
//https://qiita.com/TetsuyaFukunaga/items/4c83a8dedd34e65ffbdc
//GoでWebSocketを使いチャットサーバー構築
//https://zenn.dev/show_yeah/articles/bece10823d182c
//【入門】Goでwebsocketを使ったリアルタイムなチャットを実装してみた
