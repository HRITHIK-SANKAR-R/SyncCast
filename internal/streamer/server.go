package streamer

import (
	"bytes"
	"fmt"
	"html/template"
	"mime"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Server struct {
	FilePath  string
	Port      int
	WSHandler func(http.ResponseWriter, *http.Request)
	mu        sync.Mutex
	listener  net.Listener
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
	mux.HandleFunc("/player", s.handlePlayer)
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

	fmt.Printf("Streaming %s on port %d\n", filepath.Base(s.FilePath), s.Port)
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

func (s *Server) handlePlayer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	html, err := playerPageHTML(filepath.Base(s.FilePath))
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
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

		function wsURL() {
			const proto = location.protocol === "https:" ? "wss" : "ws";
			return proto + "://" + location.host + "/ws?role=player";
		}

		let ws;
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
		video.addEventListener("pause", () => sendStatus("pause", { position: video.currentTime }));
		video.addEventListener("ended", () => sendStatus("ended", { position: video.currentTime }));
		video.addEventListener("volumechange", () => sendStatus("volume", { volume: video.volume, muted: video.muted }));
		video.addEventListener("timeupdate", () => sendStatus("time", { position: video.currentTime }));
		video.addEventListener("error", () => sendStatus("error", { code: video.error ? video.error.code : 0 }));

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
