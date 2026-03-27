package streamer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/HRITHIK-SANKAR-R/SyncCast/internal/discovery"
)

// Supported media file extensions for the file browser.
var mediaExtensions = map[string]bool{
	".mp4": true, ".mkv": true, ".avi": true, ".mov": true,
	".webm": true, ".flv": true, ".m4v": true, ".ts": true,
	".wmv": true, ".mpg": true, ".mpeg": true, ".3gp": true,
	".mp3": true, ".flac": true, ".wav": true, ".aac": true,
	".ogg": true, ".m4a": true,
}

type Server struct {
	Port          int
	WSHandler     func(http.ResponseWriter, *http.Request)
	StateSnapshot func() interface{}
	DiscoverFunc  func(ips []string) []discovery.Device
	HostIP        string
	mu            sync.RWMutex
	filePath      string
	devices       []discovery.Device
	listener      net.Listener
}

// NewServer creates a server that starts immediately — no file required.
func NewServer(port int) *Server {
	return &Server{Port: port}
}

// New creates a server with a pre-set file (kept for test compatibility).
func New(filePath string, port int) (*Server, error) {
	if _, err := os.Stat(filePath); err != nil {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}
	return &Server{filePath: filePath, Port: port}, nil
}

// SetFile sets or changes the active streaming file at runtime.
func (s *Server) SetFile(path string) error {
	abs, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	if _, err := os.Stat(abs); err != nil {
		return fmt.Errorf("file not found: %s", abs)
	}
	s.mu.Lock()
	s.filePath = abs
	s.mu.Unlock()
	return nil
}

// GetFile returns the current file path (thread-safe).
func (s *Server) GetFile() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.filePath
}

// SetDevices updates the discovered device list (thread-safe).
func (s *Server) SetDevices(devices []discovery.Device) {
	s.mu.Lock()
	s.devices = devices
	s.mu.Unlock()
}

// GetDevices returns the current device list (thread-safe).
func (s *Server) GetDevices() []discovery.Device {
	s.mu.RLock()
	defer s.mu.RUnlock()
	devs := make([]discovery.Device, len(s.devices))
	copy(devs, s.devices)
	return devs
}

func (s *Server) StreamURL(hostIP string) string {
	return fmt.Sprintf("http://%s:%d/media", hostIP, s.Port)
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/media", s.handleMedia)
	mux.HandleFunc("/player", s.handlePlayer)
	mux.HandleFunc("/remote", serveSPA)
	mux.HandleFunc("/remote/", serveSPA)
	mux.HandleFunc("/dashboard", serveSPA)
	mux.HandleFunc("/dashboard/", serveSPA)
	mux.HandleFunc("/assets/", serveSPA)
	mux.HandleFunc("/sw.js", serveSPA)
	mux.HandleFunc("/manifest.json", serveSPA)
	mux.HandleFunc("/api/info", s.handleAPIInfo)
	mux.HandleFunc("/api/state", s.handleAPIState)
	mux.HandleFunc("/api/devices", s.handleAPIDevices)
	mux.HandleFunc("/api/files", s.handleAPIFiles)
	mux.HandleFunc("/api/stream", s.handleAPIStream)
	mux.HandleFunc("/api/discover", s.handleAPIDiscover)
	if s.WSHandler != nil {
		mux.HandleFunc("/ws", s.WSHandler)
	}

	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		return fmt.Errorf("failed to bind port %d: %w", s.Port, err)
	}

	s.mu.Lock()
	s.listener = ln
	s.mu.Unlock()

	file := s.GetFile()
	if file != "" {
		fmt.Printf("Streaming %s on port %d\n", filepath.Base(file), s.Port)
	}
	return http.Serve(ln, mux)
}

func (s *Server) Stop() error {
	s.mu.Lock()
	ln := s.listener
	s.mu.Unlock()
	if ln != nil {
		return ln.Close()
	}
	return nil
}


func (s *Server) handleMedia(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath := s.GetFile()
	if filePath == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"error": "No file selected — pick one from the dashboard",
		})
		return
	}

	f, err := os.Open(filePath)
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

	contentType := detectContentType(filePath)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	http.ServeContent(w, r, filepath.Base(filePath), stat.ModTime(), f)
}

func (s *Server) handlePlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	filePath := s.GetFile()
	fileName := "No file selected"
	if filePath != "" {
		fileName = filepath.Base(filePath)
	}

	html, err := playerPageHTML(fileName)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}


func (s *Server) handleAPIInfo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	filePath := s.GetFile()
	fileName := ""
	if filePath != "" {
		fileName = filepath.Base(filePath)
	}

	home, _ := os.UserHomeDir()

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"host":      s.HostIP,
		"port":      s.Port,
		"file":      fileName,
		"file_path": filePath,
		"streaming": filePath != "",
		"home":      home,
	})
}

func (s *Server) handleAPIState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if s.StateSnapshot != nil {
		json.NewEncoder(w).Encode(s.StateSnapshot())
	} else {
		json.NewEncoder(w).Encode(map[string]string{"state": "unknown"})
	}
}

func (s *Server) handleAPIDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	devs := s.GetDevices()
	if devs == nil {
		devs = []discovery.Device{}
	}
	json.NewEncoder(w).Encode(devs)
}

func (s *Server) handleAPIFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dir := r.URL.Query().Get("dir")
	if dir == "" {
		dir, _ = os.UserHomeDir()
	}

	dir, _ = filepath.Abs(dir)

	entries, err := os.ReadDir(dir)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Cannot read directory: " + dir})
		return
	}

	type fileEntry struct {
		Name  string `json:"name"`
		Path  string `json:"path"`
		IsDir bool   `json:"is_dir"`
		Size  int64  `json:"size"`
	}

	var files []fileEntry

	parent := filepath.Dir(dir)
	if parent != dir {
		files = append(files, fileEntry{Name: "..", Path: parent, IsDir: true})
	}

	var dirs, media []fileEntry
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		fullPath := filepath.Join(dir, name)
		if entry.IsDir() {
			dirs = append(dirs, fileEntry{Name: name, Path: fullPath, IsDir: true})
		} else {
			ext := strings.ToLower(filepath.Ext(name))
			if mediaExtensions[ext] {
				info, _ := entry.Info()
				size := int64(0)
				if info != nil {
					size = info.Size()
				}
				media = append(media, fileEntry{Name: name, Path: fullPath, Size: size})
			}
		}
	}

	sort.Slice(dirs, func(i, j int) bool { return dirs[i].Name < dirs[j].Name })
	sort.Slice(media, func(i, j int) bool { return media[i].Name < media[j].Name })
	files = append(files, dirs...)
	files = append(files, media...)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"dir":   dir,
		"files": files,
	})
}


func (s *Server) handleAPIStream(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body struct {
		File string `json:"file"`
	}
	data, _ := io.ReadAll(io.LimitReader(r.Body, 4096))
	if err := json.Unmarshal(data, &body); err != nil || body.File == "" {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Provide a valid file path"})
		return
	}

	if err := s.SetFile(body.File); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":   true,
		"file": filepath.Base(body.File),
		"path": body.File,
	})
}


func (s *Server) handleAPIDiscover(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if s.DiscoverFunc == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Discovery not available"})
		return
	}

	var body struct {
		IPs []string `json:"ips"`
	}
	data, _ := io.ReadAll(io.LimitReader(r.Body, 4096))
	if len(data) > 0 {
		json.Unmarshal(data, &body)
	}

	devices := s.DiscoverFunc(body.IPs)
	s.SetDevices(devices)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":      true,
		"count":   len(devices),
		"devices": devices,
	})
}


func playerPageHTML(fileName string) (string, error) {
	const page = `<!doctype html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width,initial-scale=1">
	<title>SyncCast Player</title>
	<style>
		:root {
			--bg: #0f172a;
			--bg2: #111827;
			--panel: rgba(17, 24, 39, 0.78);
			--text: #f8fafc;
			--muted: #cbd5e1;
			--accent: #14b8a6;
			--accent2: #22d3ee;
			--ring: rgba(34, 211, 238, 0.45);
		}
		* { box-sizing: border-box; }
		body {
			margin: 0;
			min-height: 100vh;
			display: grid;
			place-items: center;
			font-family: "Trebuchet MS", "Segoe UI", sans-serif;
			color: var(--text);
			background:
				radial-gradient(1100px 500px at 85% 8%, #0e7490 0%, transparent 60%),
				radial-gradient(900px 500px at 10% 90%, #065f46 0%, transparent 55%),
				linear-gradient(160deg, var(--bg), var(--bg2));
			overflow: hidden;
		}
		.shell {
			width: min(96vw, 1320px);
			padding: 18px;
			border: 1px solid rgba(255, 255, 255, 0.18);
			border-radius: 18px;
			background: var(--panel);
			backdrop-filter: blur(8px);
			box-shadow: 0 18px 45px rgba(0, 0, 0, 0.35);
			animation: rise 420ms ease-out;
		}
		.head {
			display: flex;
			justify-content: space-between;
			align-items: baseline;
			gap: 12px;
			margin-bottom: 12px;
		}
		.title {
			font-size: clamp(17px, 2.1vw, 28px);
			letter-spacing: 0.02em;
		}
		.badge {
			padding: 4px 10px;
			border-radius: 999px;
			font-size: 12px;
			font-weight: 700;
			text-transform: uppercase;
			letter-spacing: 0.08em;
			background: rgba(20, 184, 166, 0.2);
			color: #99f6e4;
			border: 1px solid rgba(20, 184, 166, 0.45);
		}
		.meta {
			color: var(--muted);
			font-size: 14px;
			margin-bottom: 12px;
			white-space: nowrap;
			overflow: hidden;
			text-overflow: ellipsis;
		}
		video {
			width: 100%;
			max-height: 78vh;
			border-radius: 12px;
			background: #000;
			outline: none;
			border: 1px solid rgba(255, 255, 255, 0.2);
			box-shadow: 0 0 0 0 var(--ring);
			transition: box-shadow 180ms ease;
		}
		video:focus { box-shadow: 0 0 0 4px var(--ring); }
		@keyframes rise {
			from { opacity: 0; transform: translateY(14px) scale(0.98); }
			to { opacity: 1; transform: translateY(0) scale(1); }
		}
	</style>
</head>
<body>
	<main class="shell">
		<header class="head">
			<div class="title">SyncCast TV Viewer</div>
			<div class="badge" id="status">connecting</div>
		</header>
		<div class="meta">Now playing: {{.FileName}}</div>
		<video id="video" src="/media" controls autoplay playsinline preload="metadata"></video>
	</main>

	<script>
		const video = document.getElementById("video");
		const statusEl = document.getElementById("status");
		let wakeLockSentinel = null;

		function wsURL() {
			const proto = location.protocol === "https:" ? "wss" : "ws";
			return proto + "://" + location.host + "/ws?role=player";
		}

		let ws;

		async function requestWakeLock() {
			if (!("wakeLock" in navigator)) {
				return;
			}
			if (wakeLockSentinel && !wakeLockSentinel.released) {
				return;
			}
			try {
				wakeLockSentinel = await navigator.wakeLock.request("screen");
				wakeLockSentinel.addEventListener("release", () => {
					wakeLockSentinel = null;
				});
				sendStatus("wakelock", { active: true });
			} catch (_) {
				sendStatus("wakelock", { active: false });
			}
		}

		async function releaseWakeLock() {
			if (!wakeLockSentinel) {
				return;
			}
			try {
				await wakeLockSentinel.release();
			} catch (_) {
				// Ignore release errors from browser lifecycle races.
			} finally {
				wakeLockSentinel = null;
				sendStatus("wakelock", { active: false });
			}
		}

		function connect() {
			ws = new WebSocket(wsURL());

			ws.onopen = () => {
				statusEl.textContent = "connected";
				sendStatus("ready", {
					duration: Number.isFinite(video.duration) ? video.duration : 0,
					muted: video.muted,
					volume: video.volume
				});
			};

			ws.onclose = () => {
				statusEl.textContent = "reconnecting";
				setTimeout(connect, 1000);
			};

			ws.onerror = () => {
				statusEl.textContent = "error";
			};

			ws.onmessage = (event) => {
				let msg;
				try {
					msg = JSON.parse(event.data);
				} catch (_) {
					return;
				}

				if (msg.type !== "command") {
					return;
				}

				switch (msg.action) {
					case "play":
						video.play().catch(() => {});
						break;
					case "pause":
						video.pause();
						break;
					case "seek":
						if (msg.data && typeof msg.data.position === "number") {
							video.currentTime = Math.max(0, msg.data.position);
						}
						break;
					case "back10":
						video.currentTime = Math.max(0, video.currentTime - 10);
						break;
					case "forward10":
						if (Number.isFinite(video.duration) && video.duration > 0) {
							video.currentTime = Math.min(video.duration, video.currentTime + 10);
						} else {
							video.currentTime = Math.max(0, video.currentTime + 10);
						}
						break;
					case "volume":
						if (msg.data && typeof msg.data.value === "number") {
							video.volume = Math.min(1, Math.max(0, msg.data.value));
						}
						break;
					case "stop":
						video.pause();
						video.currentTime = 0;
						break;
				}
			};
		}

		function sendStatus(action, data) {
			if (!ws || ws.readyState !== WebSocket.OPEN) {
				return;
			}
			ws.send(JSON.stringify({ type: "status", action, data }));
		}

		video.addEventListener("play", () => sendStatus("play", { position: video.currentTime }));
		video.addEventListener("play", () => requestWakeLock());
		video.addEventListener("pause", () => {
			sendStatus("pause", { position: video.currentTime });
			releaseWakeLock();
		});
		video.addEventListener("ended", () => {
			sendStatus("ended", { position: video.currentTime });
			releaseWakeLock();
		});
		video.addEventListener("volumechange", () => sendStatus("volume", { volume: video.volume, muted: video.muted }));
		video.addEventListener("timeupdate", () => sendStatus("time", { position: video.currentTime }));
		video.addEventListener("error", () => sendStatus("error", { code: video.error ? video.error.code : 0 }));

		document.addEventListener("visibilitychange", () => {
			if (document.visibilityState === "visible" && !video.paused) {
				requestWakeLock();
			}
		});

		connect();
	</script>
</body>
</html>`

	tpl, err := template.New("player").Parse(page)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	if err := tpl.Execute(&out, struct{ FileName string }{FileName: fileName}); err != nil {
		return "", err
	}

	return out.String(), nil
}

func detectContentType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))

	if ct := mime.TypeByExtension(ext); ct != "" {
		return ct
	}

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
