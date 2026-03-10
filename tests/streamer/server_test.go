package streamer_test

import (
	"os"
	"testing"

	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/streamer"
)

func TestNew_ValidFile(t *testing.T) {
	f := createTempFile(t, "test.mp4", "fake video content")
	defer os.Remove(f)

	srv, err := streamer.New(f, 0)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if srv.FilePath != f {
		t.Fatalf("expected FilePath %q, got %q", f, srv.FilePath)
	}
}

func TestNew_MissingFile(t *testing.T) {
	_, err := streamer.New("/tmp/nonexistent_synccast_test_file.mp4", 0)
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestStreamURL(t *testing.T) {
	f := createTempFile(t, "test.mp4", "data")
	defer os.Remove(f)

	srv, _ := streamer.New(f, 6969)

	url := srv.StreamURL("192.168.1.10")
	expected := "http://192.168.1.10:6969/media"
	if url != expected {
		t.Fatalf("expected %q, got %q", expected, url)
	}
}

func TestStreamURL_DifferentPort(t *testing.T) {
	f := createTempFile(t, "test.mp4", "data")
	defer os.Remove(f)

	srv, _ := streamer.New(f, 8080)

	url := srv.StreamURL("10.0.0.1")
	expected := "http://10.0.0.1:8080/media"
	if url != expected {
		t.Fatalf("expected %q, got %q", expected, url)
	}
}

func TestStop_WithoutStart(t *testing.T) {
	f := createTempFile(t, "test.mp4", "data")
	defer os.Remove(f)

	srv, _ := streamer.New(f, 0)
	err := srv.Stop()
	if err != nil {
		t.Fatalf("Stop on unstarted server should return nil, got: %v", err)
	}
}

// createTempFile creates a temporary file with the given content and returns its path.
func createTempFile(t *testing.T, name, content string) string {
	t.Helper()
	f, err := os.CreateTemp("", name)
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	f.WriteString(content)
	f.Close()
	return f.Name()
}
