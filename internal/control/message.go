package control

import "encoding/json"

// Message defines the WebSocket protocol frame exchanged between remotes and players.
type Message struct {
	Type   string          `json:"type"`
	Action string          `json:"action,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}

// Message types.
const (
	TypeCommand = "command"
	TypeStatus  = "status"
)

// Command actions.
const (
	ActionPlay   = "play"
	ActionPause  = "pause"
	ActionSeek   = "seek"
	ActionVolume = "volume"
	ActionStop   = "stop"
)
