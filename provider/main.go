package main

import (
	"encoding/json"
	"log"
	"time"

	"provider/app/session"
	"provider/app/ws"
	"provider/constants"
	"provider/settings"

	"github.com/gorilla/websocket"
)

type JoinData struct {
	Role string `json:"role"`
}

func joinAsProvider(conn *ws.Connection) error {
	joinData, err := json.Marshal(JoinData{Role: "provider"})
	if err != nil {
		return err
	}

	msg := ws.Message{
		Type: constants.JoinMessage,
		Data: string(joinData),
	}

	return conn.Send(msg)
}

// tryConnect tries to dial and setup a WS connection with Coordinator service
// maxTries = -1 means it will retry forever
func tryConnect(addr string, maxTries int) *ws.Connection {
	count := 0
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		conn, err := ws.Connect(addr)
		if err != nil {
			count++
			log.Println("Failed to connect to Coordinator", count, err)
			if count == maxTries {
				log.Println("Stop trying to connect to Coordinator")
				break
			}
			continue
		}
		if err = joinAsProvider(conn); err != nil {
			conn.Close()
			count++
			log.Println("Failed to join as a provider", count, err)
			if count == maxTries {
				log.Println("Stop trying to connect to Coordinator")
				break
			}
			continue
		}

		return conn
	}

	return nil
}

func main() {
	hub := session.NewHub()

	conn := tryConnect(settings.CoordinatorAddr, 1)
	if conn == nil {
		log.Fatalln("Couldn't connect to coordinator service")
	}
	log.Println("Connected to Coordinator service as a Provider")

	for {
		msg, err := conn.ReadMsg()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				log.Println("Reconnecting to Coordinator service..")
				conn = tryConnect(settings.CoordinatorAddr, -1)
				log.Println("Connected to Coordinator service")
			} else {
				log.Println("Error when reading WS message", err)
			}
			continue
		}

		var s *session.Session
		if msg.Type == constants.StartMessage {
			s = session.NewSession(msg.SenderID, conn, hub)
			hub.AddSession(s)
		} else {
			s = hub.GetSession(msg.SenderID)
		}

		s.ReceiveMsg(msg)
	}
}
