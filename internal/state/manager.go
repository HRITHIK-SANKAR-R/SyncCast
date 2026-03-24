package state

import (
	"sync"
	"time"
)

// ConnectionState describes the live relationship between remote and player clients.
type ConnectionState string

const (
	StateIdle          ConnectionState = "idle"
	StateWaitingRemote ConnectionState = "waiting_remote"
	StateWaitingPlayer ConnectionState = "waiting_player"
	StateConnected     ConnectionState = "connected"
	StateReconnecting  ConnectionState = "reconnecting"
)

// Snapshot is a read-only view of the current connection state.
type Snapshot struct {
	State             ConnectionState `json:"state"`
	RemoteClients     int             `json:"remote_clients"`
	PlayerClients     int             `json:"player_clients"`
	ReconnectDeadline time.Time       `json:"reconnect_deadline"`
}

// Manager tracks connection lifecycle transitions for remote/player clients.
type Manager struct {
	mu                sync.RWMutex
	state             ConnectionState
	remotes           int
	players           int
	reconnectTimeout  time.Duration
	reconnectDeadline time.Time
}

// NewManager constructs a state manager with a default reconnect timeout.
func NewManager() *Manager {
	return NewManagerWithTimeout(10 * time.Second)
}

// NewManagerWithTimeout allows overriding reconnect timeout for tests/configuration.
func NewManagerWithTimeout(reconnectTimeout time.Duration) *Manager {
	if reconnectTimeout <= 0 {
		reconnectTimeout = 10 * time.Second
	}
	return &Manager{
		state:            StateIdle,
		reconnectTimeout: reconnectTimeout,
	}
}

// OnClientConnected updates state when a WebSocket client joins.
func (m *Manager) OnClientConnected(role string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch role {
	case "remote":
		m.remotes++
	case "player":
		m.players++
	default:
		return
	}

	m.recomputeLocked(false)
}

// OnClientDisconnected updates state when a WebSocket client leaves.
func (m *Manager) OnClientDisconnected(role string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	switch role {
	case "remote":
		if m.remotes > 0 {
			m.remotes--
		}
	case "player":
		if m.players > 0 {
			m.players--
		}
	default:
		return
	}

	m.recomputeLocked(true)
}

// Tick advances time-dependent transitions (reconnect timeout expiry).
func (m *Manager) Tick(now time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.state != StateReconnecting {
		return
	}
	if m.reconnectDeadline.IsZero() || now.Before(m.reconnectDeadline) {
		return
	}

	m.reconnectDeadline = time.Time{}
	switch {
	case m.remotes > 0 && m.players > 0:
		m.state = StateConnected
	case m.remotes > 0:
		m.state = StateWaitingPlayer
	case m.players > 0:
		m.state = StateWaitingRemote
	default:
		m.state = StateIdle
	}
}

// Snapshot returns the latest state and participant counts.
func (m *Manager) Snapshot() Snapshot {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return Snapshot{
		State:             m.state,
		RemoteClients:     m.remotes,
		PlayerClients:     m.players,
		ReconnectDeadline: m.reconnectDeadline,
	}
}

func (m *Manager) recomputeLocked(onDisconnect bool) {
	prev := m.state

	switch {
	case m.remotes > 0 && m.players > 0:
		m.state = StateConnected
		m.reconnectDeadline = time.Time{}
	case m.remotes == 0 && m.players == 0:
		m.state = StateIdle
		m.reconnectDeadline = time.Time{}
	default:
		if onDisconnect && prev == StateConnected {
			m.state = StateReconnecting
			m.reconnectDeadline = time.Now().Add(m.reconnectTimeout)
			return
		}

		// While reconnect window is open, keep reconnecting state stable.
		if prev == StateReconnecting && !m.reconnectDeadline.IsZero() && time.Now().Before(m.reconnectDeadline) {
			m.state = StateReconnecting
			return
		}

		m.reconnectDeadline = time.Time{}
		if m.remotes > 0 {
			m.state = StateWaitingPlayer
		} else {
			m.state = StateWaitingRemote
		}
	}
}
