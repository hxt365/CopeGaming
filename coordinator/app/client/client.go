package client

import (
	"log"
	"sync"
	"time"

	"coordinator/constants"
	"coordinator/utils"

	"github.com/gorilla/websocket"
)

const (

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10240
)

type ProviderInfo struct {
	OwnerID    string
	HostName   string
	Platform   string
	CpuName    string
	CpuNum     int
	MemSize    float64
	CpuPercent float64
	MemPercent float64
}

type Client struct {
	ID        string
	role      constants.Role
	hub       *Hub
	conn      *websocket.Conn
	outputBuf chan interface{}
	// Info of provider
	Provider *ProviderInfo
}

func NewClient(id string, conn *websocket.Conn, hub *Hub) *Client {
	c := &Client{
		ID:        id,
		conn:      conn,
		hub:       hub,
		outputBuf: make(chan interface{}),
	}

	go c.readPump()
	go c.writePump()

	return c
}

func (c *Client) close() {
	close(c.outputBuf)
	c.conn.Close()
	c.hub.RemoveClient(c)
}

func (c *Client) sendMsg(receiver *Client, msg interface{}) {
	receiver.outputBuf <- msg
}

func (c *Client) readPump() {
	defer func() {
		c.close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, rawMsg, err := c.conn.ReadMessage()
		if err != nil {
			if _, ok := err.(*websocket.CloseError); ok {
				return
			} else {
				log.Println("Couldn't read WS message", err)
			}
		}
		msg, err := parseMsg(rawMsg)
		if err != nil {
			log.Println("Couldn't parse WS message", err)
			continue
		}

		c.handleMsg(msg)
	}
}

func (c *Client) handleJoinMsg(msg *Message) error {
	joinData, err := parseJoinData(msg.Data)
	if err != nil {
		return err
	}

	if joinData.Role == constants.Provider {
		ownerID := joinData.OwnerID
		if ownerID == "" {
			ownerID = utils.RandString(6)
		}
		c.role = constants.Provider
		c.Provider = &ProviderInfo{
			OwnerID:    ownerID,
			HostName:   joinData.HostName,
			Platform:   joinData.Platform,
			CpuName:    joinData.CpuName,
			CpuNum:     joinData.CpuNum,
			MemSize:    joinData.MemSize,
			CpuPercent: joinData.CpuPercent,
			MemPercent: joinData.MemPercent,
		}

		c.sendMsg(c, Message{
			Type: constants.JoinAcceptedMessage,
			Data: ownerID,
		})
	} else {
		c.role = constants.Player
	}

	return nil
}

func (c *Client) handleStatsMsg(msg *Message) error {
	statsData, err := parseStatsData(msg.Data)
	if err != nil {
		return err
	}

	if c.role == constants.Provider {
		c.Provider.CpuPercent = statsData.CpuPercent
		c.Provider.MemPercent = statsData.MemPercent
	}

	return nil
}

func (c *Client) handleMsg(msg *Message) {
	switch msg.Type {
	case constants.JoinMessage:
		if err := c.handleJoinMsg(msg); err != nil {
			return
		}
	case constants.StatsMessage:
		if err := c.handleStatsMsg(msg); err != nil {
			return
		}
	default:
		receiver := c.hub.GetClient(msg.ReceiverID)
		if receiver == nil {
			return
		}

		msg.SenderID = c.ID
		c.sendMsg(receiver, msg)
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case msg, ok := <-c.outputBuf:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteJSON(msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type Hub struct {
	clients map[string]*Client
	rwMutex sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
		rwMutex: sync.RWMutex{},
	}
}

func (h *Hub) AddClient(c *Client) {
	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()

	h.clients[c.ID] = c
}

func (h *Hub) RemoveClient(c *Client) {
	h.rwMutex.Lock()
	defer h.rwMutex.Unlock()

	if _, ok := h.clients[c.ID]; ok {
		delete(h.clients, c.ID)
	}
}

func (h *Hub) GetClient(id string) *Client {
	h.rwMutex.RLock()
	defer h.rwMutex.RUnlock()

	if c, ok := h.clients[id]; ok {
		return c
	}

	return nil
}

func (h *Hub) GetProviders() []*Client {
	h.rwMutex.RLock()
	defer h.rwMutex.RUnlock()

	var providers []*Client

	for _, c := range h.clients {
		if c.role == constants.Provider {
			providers = append(providers, c)
		}
	}

	return providers
}
