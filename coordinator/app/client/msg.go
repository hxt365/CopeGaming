package client

import (
	"encoding/json"

	"coordinator/constants"
)

type Message struct {
	SenderID   string                `json:"senderID"`
	ReceiverID string                `json:"receiverID"`
	Type       constants.MessageType `json:"type"`
	Data       string                `json:"data"`
}

func parseMsg(raw []byte) (*Message, error) {
	var msg Message

	if err := json.Unmarshal(raw, &msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

type JoinData struct {
	Role       constants.Role `json:"role"`
	OwnerID    string         `json:"ownerID"`
	HostName   string         `json:"hostName"`
	Platform   string         `json:"platform"`
	CpuName    string         `json:"cpuName"`
	CpuNum     int            `json:"cpuNum"`
	MemSize    float64        `json:"memSize"`
	CpuPercent float64        `json:"cpuPercent"`
	MemPercent float64        `json:"memPercent"`
}

func parseJoinData(raw string) (*JoinData, error) {
	var join JoinData

	if err := json.Unmarshal([]byte(raw), &join); err != nil {
		return nil, err
	}

	return &join, nil
}

type StatsData struct {
	CpuPercent float64 `json:"cpuPercent"`
	MemPercent float64 `json:"memPercent"`
}

func parseStatsData(raw string) (*StatsData, error) {
	var stats StatsData

	if err := json.Unmarshal([]byte(raw), &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}
