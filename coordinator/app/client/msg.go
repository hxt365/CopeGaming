package client

import "encoding/json"

type MessageType string

type Role string

const (
	Provider Role = "provider"
	Player   Role = "player"

	JoinMessage MessageType = "join"
)

type Message struct {
	SenderID   string      `json:"senderID"`
	ReceiverID string      `json:"receiverID"`
	Type       MessageType `json:"type"`
	Data       string      `json:"data"`
}

type JoinData struct {
	Role Role `json:"role"`
}

func parseMsg(raw []byte) (*Message, error) {
	var msg Message

	if err := json.Unmarshal(raw, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

func parseJoinData(raw string) (*JoinData, error) {
	var join JoinData

	if err := json.Unmarshal([]byte(raw), &join); err != nil {
		return nil, err
	}

	return &join, nil
}
