package streamer

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestHandlePlayer_GetServesHTML(t *testing.T) {
	f, err := os.CreateTemp("", "synccast_player_*.mp4")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	_, _ = f.WriteString("fake-media")
	_ = f.Close()

	srv, err := New(f.Name(), 0)
	if err != nil {
		t.Fatalf("new server: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/player", nil)
	rr := httptest.NewRecorder()
	srv.handlePlayer(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if got := rr.Header().Get("Content-Type"); !strings.Contains(got, "text/html") {
		t.Fatalf("expected text/html content type, got %q", got)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "SyncCast TV Viewer") {
		t.Fatalf("expected player title in body")
	}
	if !strings.Contains(body, "src=\"/media\"") {
		t.Fatalf("expected /media source in player page")
	}
	if !strings.Contains(body, "role=player") {
		t.Fatalf("expected player websocket role in script")
	}
}

func TestHandlePlayer_MethodNotAllowed(t *testing.T) {
	f, err := os.CreateTemp("", "synccast_player_*.mp4")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	defer os.Remove(f.Name())
	_, _ = f.WriteString("fake-media")
	_ = f.Close()

	srv, err := New(f.Name(), 0)
	if err != nil {
		t.Fatalf("new server: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/player", nil)
	rr := httptest.NewRecorder()
	srv.handlePlayer(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rr.Code)
	}
}
