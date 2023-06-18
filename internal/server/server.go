package server

import (
	"embed"
	"fmt"
	"net/http"
	"strings"
)

var staticDir = "web"

func rootPath(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")

		if r.URL.Path == "/" {
			r.URL.Path = fmt.Sprintf("/%s/", staticDir)
		} else {
			b := strings.Split(r.URL.Path, "/")[0]
			if b != staticDir {
				r.URL.Path = fmt.Sprintf("/%s%s", staticDir, r.URL.Path)
			}
		}
		h.ServeHTTP(w, r)
	})
}

func RunServer(staticFiles embed.FS, port int, json string) {
	var staticFS = http.FS(staticFiles)
	fs := rootPath(http.FileServer(staticFS))

	http.Handle("/", fs)

	http.HandleFunc("/graph", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(json))
	})

	fmt.Printf("Server is running on http://127.0.0.1:%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		panic(err)
	}
}
