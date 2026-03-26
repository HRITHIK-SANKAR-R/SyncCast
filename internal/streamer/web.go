package streamer

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"
)

var webDist embed.FS

func webDistFS() http.FileSystem {
	sub, err := fs.Sub(webDist, "web_dist")
	if err != nil {
		panic("embedded web_dist not found: " + err.Error())
	}
	return http.FS(sub)
}


func serveSPA(w http.ResponseWriter, r *http.Request) {
	fsys := webDistFS()

	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" {
		path = "index.html"
	}

	f, err := fsys.Open(path)
	if err == nil {
		f.Close()
		http.FileServer(fsys).ServeHTTP(w, r)
		return
	}

	r.URL.Path = "/"
	http.FileServer(fsys).ServeHTTP(w, r)
}
