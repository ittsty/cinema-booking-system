package ws

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{

	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServeWS(hub *Hub, c *gin.Context) {

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		return
	}

	client := &Client{
		Conn: conn,
		Send: make(chan []byte),
	}

	hub.Register <- client

	go func() {
		for {
			message := <-client.Send
			conn.WriteMessage(websocket.TextMessage, message)
		}
	}()
}
