package ws

import (
	"encoding/json"
	"log"
	"net/url"
	"sync"

	"provider/constants"

	"github.com/gorilla/websocket"
)

type Message struct {
	SenderID   string                `json:"senderID"`
	ReceiverID string                `json:"receiverID"`
	Type       constants.MessageType `json:"type"`
	Data       string                `json:"data"`
}

type Connection struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

func Connect(addr string) (*Connection, error) {
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	return &Connection{conn: c}, nil
}

func (c *Connection) Send(v interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.conn.WriteJSON(v)
}

func (c *Connection) ReadMsg() (*Message, error) {
	for {
		msgType, rawMsg, err := c.conn.ReadMessage()
		if err != nil {
			return nil, err
		}

		if msgType != websocket.TextMessage {
			continue
		}

		var msg Message
		if err := json.Unmarshal(rawMsg, &msg); err != nil {
			log.Println("Couldn't unmarshal WS message", err)
			continue
		}

		return &msg, nil
	}
}

func (c *Connection) Close() error {
	return c.conn.Close()
}
