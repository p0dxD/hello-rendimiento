// hello-rendimiento is a tiny web app used to test rendimiento.ai end to end.
package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

var (
	version = "dev" // set at build time from GIT_SHA
	started = time.Now()
	visits  atomic.Int64
)

var page = template.Must(template.New("page").Parse(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<title>hello from rendimiento</title>
<style>
  body { font-family: system-ui, sans-serif; display: grid; place-items: center; min-height: 100vh; margin: 0; background: #0d0e12; color: #ececf1; }
  main { text-align: center; } h1 { font-size: 2.4rem; margin: 0 0 .5rem; } p { color: #9a9aab; }
  code { background: #1c1d25; padding: 2px 8px; border-radius: 6px; color: #818cf8; }
</style></head>
<body><main>
  <h1>👋👋 {{.Greeting}}</h1>
  <p>Deployed by <b>rendimiento.ai</b> · version <code>{{.Version}}</code></p>
  <p>host <code>{{.Host}}</code> · up {{.Uptime}} · visit #{{.Visits}}</p>
</main></body></html>`))

func greeting() string {
	if g := os.Getenv("GREETING"); g != "" {
		return g
	}
	return "Hello, world"
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		host, _ := os.Hostname()
		_ = page.Execute(w, map[string]any{
			"Greeting": greeting(),
			"Version":  version,
			"Host":     host,
			"Uptime":   time.Since(started).Round(time.Second),
			"Visits":   visits.Add(1),
		})
	})
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	mux.HandleFunc("GET /api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"greeting": greeting(), "version": version})
	})
	return mux
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("hello-rendimiento %s listening on :%s", version, port)
	srv := &http.Server{Addr: ":" + port, Handler: handler(), ReadHeaderTimeout: 5 * time.Second}
	log.Fatal(srv.ListenAndServe())
}
