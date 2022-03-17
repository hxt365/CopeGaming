package ws

import (
	"log"
	"net/http"

	"coordinator/app/client"
	"coordinator/settings"
	"coordinator/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		for _, origin := range settings.AllowedWSOrigins {
			if r.Header.Get("Origin") == origin || origin == "*" {
				return true
			}
		}

		return false
	},
}

func ServeWs(hub *client.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade connection from", r.RemoteAddr, err)
		return
	}

	randID := utils.RandString(6)
	c := client.NewClient(randID, conn, hub)

	hub.AddClient(c)
}
