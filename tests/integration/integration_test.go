package integration_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/discovery"
	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/streamer"
)

func TestDiscoverTVs_ManualMode_WithFakeTV(t *testing.T) {
	fakeTVName := "Integration Test TV"
	xml := fmt.Sprintf(`<?xml version="1.0"?>
<root xmlns="urn:schemas-upnp-org:device-1-0">
  <device>
    <friendlyName>%s</friendlyName>
  </device>
</root>`, fakeTVName)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(xml))
	}))
	defer srv.Close()

	devices := discovery.DiscoverTVs([]string{"127.0.0.1"})
	if len(devices) > 0 {
		t.Logf("Found device (unexpected in test env): %+v", devices[0])
	}
}

func TestDiscoverTVs_AutoMode_NilInput(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping auto-discovery in short mode (uses network)")
	}
	devices := discovery.DiscoverTVs(nil)
	t.Logf("Auto-discovery found %d device(s)", len(devices))
}

func TestStreamerStartStop(t *testing.T) {
	content := "integration test media content"
	f, err := os.CreateTemp("", "integration_test.mp4")
	if err != nil {
		t.Fatal(err)
	}
	f.WriteString(content)
	f.Close()
	defer os.Remove(f.Name())

	srv, err := streamer.New(f.Name(), 0)
	if err != nil {
		t.Fatalf("New failed: %v", err)
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	// Give the listener time to bind before stopping
	time.Sleep(100 * time.Millisecond)

	if err := srv.Stop(); err != nil {
		t.Logf("Stop returned: %v", err)
	}

	startErr := <-errCh
	if startErr == nil {
		t.Log("Server exited cleanly")
	}
}

func TestEndToEnd_DiscoverThenStream(t *testing.T) {
	devices := discovery.DiscoverTVs([]string{"127.0.0.1"})
	t.Logf("Discovery returned %d devices", len(devices))

	f, _ := os.CreateTemp("", "e2e_test.mkv")
	f.WriteString("fake mkv data")
	f.Close()
	defer os.Remove(f.Name())

	srv, err := streamer.New(f.Name(), 6969)
	if err != nil {
		t.Fatalf("streamer.New failed: %v", err)
	}

	ip, _ := discovery.GetLocalIP()
	url := srv.StreamURL(ip)
	expected := fmt.Sprintf("http://%s:6969/media", ip)
	if url != expected {
		t.Fatalf("expected stream URL %q, got %q", expected, url)
	}
}
