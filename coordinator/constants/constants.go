package constants

type MessageType string

type Role string

const (
	Provider Role = "provider"
	Player   Role = "player"

	JoinMessage  MessageType = "join"
	StatsMessage MessageType = "stats"
)
