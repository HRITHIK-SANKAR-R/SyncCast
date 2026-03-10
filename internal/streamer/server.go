package streamer

import (
	"fmt"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Server struct {
	FilePath string
	Port     int
	listener net.Listener
}

func New(filePath string, port int) (*Server, error) {
	if _, err := os.Stat(filePath); err != nil {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}
	return &Server{FilePath: filePath, Port: port}, nil
}

func (s *Server) StreamURL(hostIP string) string {
	return fmt.Sprintf("http://%s:%d/media", hostIP, s.Port)
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/media", s.handleMedia)

	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		return fmt.Errorf("failed to bind port %d: %w", s.Port, err)
	}

	fmt.Printf("Streaming %s on port %d\n", filepath.Base(s.FilePath), s.Port)
	return http.Serve(s.listener, mux)
}

func (s *Server) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) handleMedia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	f, err := os.Open(s.FilePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	contentType := detectContentType(s.FilePath)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// http.ServeContent handles Range headers, If-Modified-Since,
	// and responds with 206 Partial Content when appropriate
	http.ServeContent(w, r, filepath.Base(s.FilePath), stat.ModTime(), f)
}

func detectContentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))

	// Prefer mime.TypeByExtension for well-known media types
	if ct := mime.TypeByExtension(ext); ct != "" {
		return ct
	}

	// Fallback for types not always in the system MIME database
	switch ext {
	case ".mkv":
		return "video/x-matroska"
	case ".webm":
		return "video/webm"
	case ".flv":
		return "video/x-flv"
	case ".ts":
		return "video/mp2t"
	case ".m4v":
		return "video/x-m4v"
	default:
		return "application/octet-stream"
	}
}