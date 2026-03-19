package state_test

import (
	"testing"
	"time"

	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/state"
)

func TestManager_ConnectFlowToConnected(t *testing.T) {
	m := state.NewManagerWithTimeout(2 * time.Second)

	m.OnClientConnected("remote")
	s := m.Snapshot()
	if s.State != state.StateWaitingPlayer {
		t.Fatalf("expected state %q, got %q", state.StateWaitingPlayer, s.State)
	}

	m.OnClientConnected("player")
	s = m.Snapshot()
	if s.State != state.StateConnected {
		t.Fatalf("expected state %q, got %q", state.StateConnected, s.State)
	}
	if s.RemoteClients != 1 || s.PlayerClients != 1 {
		t.Fatalf("unexpected counts remote=%d player=%d", s.RemoteClients, s.PlayerClients)
	}
}

func TestManager_EntersReconnectingOnDrop(t *testing.T) {
	m := state.NewManagerWithTimeout(1 * time.Second)

	m.OnClientConnected("remote")
	m.OnClientConnected("player")
	m.OnClientDisconnected("player")

	s := m.Snapshot()
	if s.State != state.StateReconnecting {
		t.Fatalf("expected state %q, got %q", state.StateReconnecting, s.State)
	}
	if s.ReconnectDeadline.IsZero() {
		t.Fatal("expected reconnect deadline to be set")
	}
}

func TestManager_ReconnectBeforeDeadlineReturnsConnected(t *testing.T) {
	m := state.NewManagerWithTimeout(2 * time.Second)

	m.OnClientConnected("remote")
	m.OnClientConnected("player")
	m.OnClientDisconnected("player")

	m.OnClientConnected("player")
	s := m.Snapshot()
	if s.State != state.StateConnected {
		t.Fatalf("expected state %q, got %q", state.StateConnected, s.State)
	}
}

func TestManager_ReconnectTimeoutFallsBackWaiting(t *testing.T) {
	m := state.NewManagerWithTimeout(40 * time.Millisecond)

	m.OnClientConnected("remote")
	m.OnClientConnected("player")
	m.OnClientDisconnected("player")

	time.Sleep(60 * time.Millisecond)
	m.Tick(time.Now())

	s := m.Snapshot()
	if s.State != state.StateWaitingPlayer {
		t.Fatalf("expected state %q, got %q", state.StateWaitingPlayer, s.State)
	}
}

func TestManager_IdleWhenAllClientsLeave(t *testing.T) {
	m := state.NewManager()

	m.OnClientConnected("remote")
	m.OnClientDisconnected("remote")

	s := m.Snapshot()
	if s.State != state.StateIdle {
		t.Fatalf("expected state %q, got %q", state.StateIdle, s.State)
	}
}
