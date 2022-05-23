package main

import (
	"encoding/json"
	"flag"
	"log"
	"time"

	"provider/app/session"
	"provider/app/stats"
	"provider/app/ws"
	"provider/constants"
	"provider/settings"

	"github.com/gorilla/websocket"
)

type JoinData struct {
	Role       string  `json:"role"`
	OwnerID    string  `json:"ownerID"`
	HostName   string  `json:"hostName"`
	Platform   string  `json:"platform"`
	CpuName    string  `json:"cpuName"`
	CpuNum     int     `json:"cpuNum"`
	MemSize    float64 `json:"memSize"`
	CpuPercent float64 `json:"cpuPercent"`
	MemPercent float64 `json:"memPercent"`
}

func joinAsProvider(ownerID string, conn *ws.Connection) error {
	sysInfo, err := stats.GetSysInfo()
	if err != nil {
		return err
	}
	sysStats, err := stats.GetSysStats(5 * time.Second)
	if err != nil {
		return err
	}

	joinData, err := json.Marshal(JoinData{
		Role:       "provider",
		OwnerID:    ownerID,
		HostName:   sysInfo.HostName,
		Platform:   sysInfo.Platform,
		CpuName:    sysInfo.CpuName,
		CpuNum:     sysInfo.CpuNum,
		MemSize:    sysInfo.MemSize,
		CpuPercent: sysStats.CpuPercent,
		MemPercent: sysStats.MemPercent,
	})
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
func tryConnect(ownerID, addr string, maxTries int) *ws.Connection {
	count := 0
	ticker := time.NewTicker(5 * time.Second)
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
		if err = joinAsProvider(ownerID, conn); err != nil {
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

type StatsData struct {
	CpuPercent float64 `json:"cpuPercent"`
	MemPercent float64 `json:"memPercent"`
}

func updateStats(conn *ws.Connection, interval time.Duration) {
	for {
		sysStats, err := stats.GetSysStats(interval)
		if err != nil {
			log.Println("Couldn't get system stats", err)
			continue
		}

		statsData, err := json.Marshal(StatsData{
			CpuPercent: sysStats.CpuPercent,
			MemPercent: sysStats.MemPercent,
		})
		if err != nil {
			log.Println("Couldn't marshal stats data", err)
			continue
		}

		msg := ws.Message{
			Type: constants.StatsMessage,
			Data: string(statsData),
		}

		conn.Send(msg)
	}
}

var ownerID = flag.String("owner", "", "ID of this computer's owner")

func main() {
	flag.Parse()

	hub := session.NewHub()

	conn := tryConnect(*ownerID, settings.CoordinatorAddr, 1)
	if conn == nil {
		log.Fatalln("Couldn't connect to coordinator service")
	}
	log.Println("Connected to Coordinator service as a Provider")

	go updateStats(conn, 5*time.Second)

	for {
		msg, err := conn.ReadMsg()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				log.Println("Reconnecting to Coordinator service..")
				conn = tryConnect(*ownerID, settings.CoordinatorAddr, -1)
				log.Println("Connected to Coordinator service")
			} else {
				log.Println("Error when reading WS message", err)
			}
			continue
		}

		var s *session.Session
		if msg.Type == constants.JoinAcceptedMessage {
			log.Printf("Owner's ID: %s", msg.Data)
			continue
		} else if msg.Type == constants.StartMessage {
			s = session.NewSession(msg.SenderID, conn, hub)
			hub.AddSession(s)
		} else {
			s = hub.GetSession(msg.SenderID)
		}

		s.ReceiveMsg(msg)
	}
}
