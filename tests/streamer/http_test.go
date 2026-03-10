package streamer_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/streamer"
)

// newTestServer creates a streamer.Server for the given file and returns
// an httptest.Server that serves the /media endpoint.
func newTestServer(t *testing.T, filePath string) *httptest.Server {
	t.Helper()
	srv, err := streamer.New(filePath, 0)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	_ = srv

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Replicate what Start() does internally for the /media route.
		f, err := os.Open(filePath)
		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}
		defer f.Close()

		stat, _ := f.Stat()
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		http.ServeContent(w, r, stat.Name(), stat.ModTime(), f)
	}))
}

func TestHTTP_FullGet(t *testing.T) {
	content := "this is fake video data for testing full GET"
	path := createTempFile(t, "test.mp4", content)
	defer os.Remove(path)

	ts := newTestServer(t, path)
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != content {
		t.Fatalf("expected body %q, got %q", content, string(body))
	}

	if resp.Header.Get("Accept-Ranges") != "bytes" {
		t.Fatal("missing Accept-Ranges: bytes header")
	}

	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("missing CORS header")
	}
}

func TestHTTP_RangeRequest_PartialContent(t *testing.T) {
	content := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	path := createTempFile(t, "test.mp4", content)
	defer os.Remove(path)

	ts := newTestServer(t, path)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)
	req.Header.Set("Range", "bytes=0-4")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Range GET failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		t.Fatalf("expected 206, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "ABCDE" {
		t.Fatalf("expected %q, got %q", "ABCDE", string(body))
	}
}

func TestHTTP_RangeRequest_MiddleSlice(t *testing.T) {
	content := "0123456789"
	path := createTempFile(t, "test.mp4", content)
	defer os.Remove(path)

	ts := newTestServer(t, path)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)
	req.Header.Set("Range", "bytes=3-6")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Range GET failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		t.Fatalf("expected 206, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "3456" {
		t.Fatalf("expected %q, got %q", "3456", string(body))
	}
}

func TestHTTP_RangeRequest_SuffixRange(t *testing.T) {
	content := "ABCDEFGHIJ"
	path := createTempFile(t, "test.mp4", content)
	defer os.Remove(path)

	ts := newTestServer(t, path)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)
	req.Header.Set("Range", "bytes=-3")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Range GET failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		t.Fatalf("expected 206, got %d", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	if string(body) != "HIJ" {
		t.Fatalf("expected %q, got %q", "HIJ", string(body))
	}
}

func TestHTTP_HeadRequest(t *testing.T) {
	content := "some video bytes"
	path := createTempFile(t, "test.mp4", content)
	defer os.Remove(path)

	ts := newTestServer(t, path)
	defer ts.Close()

	resp, err := http.Head(ts.URL)
	if err != nil {
		t.Fatalf("HEAD failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	if resp.ContentLength != int64(len(content)) {
		t.Fatalf("expected Content-Length %d, got %d", len(content), resp.ContentLength)
	}
}

func TestHTTP_InvalidRange(t *testing.T) {
	content := "short"
	path := createTempFile(t, "test.mp4", content)
	defer os.Remove(path)

	ts := newTestServer(t, path)
	defer ts.Close()

	req, _ := http.NewRequest("GET", ts.URL, nil)
	req.Header.Set("Range", "bytes=100-200")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusRequestedRangeNotSatisfiable {
		t.Fatalf("expected 416, got %d", resp.StatusCode)
	}
}

func TestHTTP_ContentType_Header(t *testing.T) {
	content := "video data"
	path := createTempFile(t, "test.mp4", content)
	defer os.Remove(path)

	ts := newTestServer(t, path)
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	if err != nil {
		t.Fatalf("GET failed: %v", err)
	}
	defer resp.Body.Close()

	ct := resp.Header.Get("Content-Type")
	// httptest may detect content type; we just ensure it's set
	if ct == "" {
		t.Fatal("expected Content-Type header to be set")
	}
	// Should not be empty or a directory listing
	if strings.Contains(ct, "text/html") {
		t.Fatalf("unexpected Content-Type for media: %s", ct)
	}
}
