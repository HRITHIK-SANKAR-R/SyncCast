package control_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/control"
	"github.com/gorilla/websocket"
)

func startHubTestServer(t *testing.T) (*control.Hub, *httptest.Server, string) {
	t.Helper()

	hub := control.NewHub()
	go hub.Run()

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", hub.HandleWS)
	ts := httptest.NewServer(mux)

	wsBase := "ws" + strings.TrimPrefix(ts.URL, "http")
	return hub, ts, wsBase
}

func dialRole(t *testing.T, wsBase, role string) *websocket.Conn {
	t.Helper()

	conn, _, err := websocket.DefaultDialer.Dial(wsBase+"/ws?role="+role, nil)
	if err != nil {
		t.Fatalf("failed to dial role %q: %v", role, err)
	}
	return conn
}

func waitForClientCount(t *testing.T, hub *control.Hub, expected int) {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if hub.ClientCount() == expected {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("expected client count %d, got %d", expected, hub.ClientCount())
}

func TestHandleWS_InvalidRoleReturnsBadRequest(t *testing.T) {
	hub, ts, _ := startHubTestServer(t)
	defer ts.Close()
	defer hub.Stop()

	resp, err := http.Get(ts.URL + "/ws?role=invalid")
	if err != nil {
		t.Fatalf("http get failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestHubRoutesRemoteToPlayer(t *testing.T) {
	hub, ts, wsBase := startHubTestServer(t)
	defer ts.Close()
	defer hub.Stop()

	player := dialRole(t, wsBase, string(control.RolePlayer))
	defer player.Close()

	remote := dialRole(t, wsBase, string(control.RoleRemote))
	defer remote.Close()

	waitForClientCount(t, hub, 2)

	want := `{"type":"command","action":"play"}`
	if err := remote.WriteMessage(websocket.TextMessage, []byte(want)); err != nil {
		t.Fatalf("remote write failed: %v", err)
	}

	_ = player.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, got, err := player.ReadMessage()
	if err != nil {
		t.Fatalf("player read failed: %v", err)
	}
	if string(got) != want {
		t.Fatalf("expected %q, got %q", want, string(got))
	}
}

func TestHubRoutesPlayerToRemote(t *testing.T) {
	hub, ts, wsBase := startHubTestServer(t)
	defer ts.Close()
	defer hub.Stop()

	remote := dialRole(t, wsBase, string(control.RoleRemote))
	defer remote.Close()

	player := dialRole(t, wsBase, string(control.RolePlayer))
	defer player.Close()

	waitForClientCount(t, hub, 2)

	want := `{"type":"status","data":{"position":42}}`
	if err := player.WriteMessage(websocket.TextMessage, []byte(want)); err != nil {
		t.Fatalf("player write failed: %v", err)
	}

	_ = remote.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, got, err := remote.ReadMessage()
	if err != nil {
		t.Fatalf("remote read failed: %v", err)
	}
	if string(got) != want {
		t.Fatalf("expected %q, got %q", want, string(got))
	}
}

func TestHubClientCountDecrementsOnDisconnect(t *testing.T) {
	hub, ts, wsBase := startHubTestServer(t)
	defer ts.Close()
	defer hub.Stop()

	remote := dialRole(t, wsBase, string(control.RoleRemote))
	waitForClientCount(t, hub, 1)

	if err := remote.Close(); err != nil {
		t.Fatalf("remote close failed: %v", err)
	}

	waitForClientCount(t, hub, 0)
}

func TestHubLifecycleHooksCalled(t *testing.T) {
	hub, ts, wsBase := startHubTestServer(t)
	defer ts.Close()
	defer hub.Stop()

	var connected int32
	var disconnected int32
	hub.SetLifecycleHooks(
		func(role control.Role) {
			atomic.AddInt32(&connected, 1)
		},
		func(role control.Role) {
			atomic.AddInt32(&disconnected, 1)
		},
	)

	remote := dialRole(t, wsBase, string(control.RoleRemote))
	waitForClientCount(t, hub, 1)

	if err := remote.Close(); err != nil {
		t.Fatalf("remote close failed: %v", err)
	}
	waitForClientCount(t, hub, 0)

	deadline := time.Now().Add(1 * time.Second)
	for time.Now().Before(deadline) {
		if atomic.LoadInt32(&connected) == 1 && atomic.LoadInt32(&disconnected) == 1 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("expected hooks called once each, got connected=%d disconnected=%d", connected, disconnected)
}
