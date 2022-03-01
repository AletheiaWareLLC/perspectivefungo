package main

import (
	"aletheiaware.com/netgo"
	"aletheiaware.com/netgo/handler"
	"crypto/tls"
	"embed"
	"errors"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

//go:embed assets
var embeddedFS embed.FS

func main() {
	// Configure Logging
	logFile, err := netgo.SetupLogging()
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()
	log.Println("Log File:", logFile.Name())

	// Create Multiplexer
	mux := http.NewServeMux()

	handler.AttachHealthHandler(mux)

	// Handle Static Assets
	staticFS, err := fs.Sub(embeddedFS, path.Join("assets", "static"))
	if err != nil {
		log.Fatal(err)
	}
	//handler.AttachStaticFSHandler(mux, staticFS, false, fmt.Sprintf("public, max-age=%d", 60*60*24*7*52)) // 52 week max-age
	handler.AttachStaticFSHandler(mux, staticFS, false, "no-cache")

	// Parse Templates
	templateFS, err := fs.Sub(embeddedFS, path.Join("assets", "template"))
	if err != nil {
		log.Fatal(err)
	}
	templates, err := template.ParseFS(templateFS, "*.go.html")
	if err != nil {
		log.Fatal(err)
	}

	AttachAssetHandlers(mux)

	puzzles, ok := os.LookupEnv("PUZZLE_DIRECTORY")
	if !ok {
		puzzles = "puzzles"
	}
	if err := os.MkdirAll(puzzles, os.ModePerm); err != nil {
		log.Fatal(err)
	}
	log.Println("Puzzles Directory:", puzzles)

	mux.Handle("/daily.json", handler.Log(handler.Compress(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join(puzzles, time.Now().UTC().Format("2006-01-02")+".json"))
	}))))

	mux.Handle("/daily", handler.Log(handler.Compress(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Live bool
		}{
			Live: netgo.IsLive(),
		}
		if err := templates.ExecuteTemplate(w, "daily.go.html", data); err != nil {
			log.Println(err)
			return
		}
	}))))

	// Handle favicon.ico
	mux.Handle("/favicon.ico", handler.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/daily.gif", http.StatusFound)
	})))

	// Handle robots.txt
	mux.Handle("/robots.txt", handler.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/robots.txt", http.StatusFound)
	})))

	// Handle sitemap.txt
	mux.Handle("/sitemap.txt", handler.Log(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/sitemap.txt", http.StatusFound)
	})))

	mux.Handle("/", handler.Compress(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if p := strings.TrimSuffix(r.URL.Path, "index.html"); p != "/" {
			log.Println(r.RemoteAddr, r.Proto, r.Method, r.Host, r.URL, r.Header, "not found")
			http.NotFound(w, r)
			return
		}
		netgo.LogRequest(r)
		data := struct {
			Live bool
			Date string
		}{
			Live: netgo.IsLive(),
			Date: time.Now().UTC().Format("2006-01-02"),
		}
		if err := templates.ExecuteTemplate(w, "index.go.html", data); err != nil {
			log.Println(err)
			return
		}
	})))

	if netgo.IsSecure() {
		host, ok := os.LookupEnv("HOST")
		if !ok {
			log.Fatal(errors.New("Missing HOST environment variable"))
		}

		routeMap := make(map[string]bool)

		routes, ok := os.LookupEnv("ROUTES")
		if ok {
			for _, route := range strings.Split(routes, ",") {
				routeMap[route] = true
			}
		}

		// Redirect HTTP Requests to HTTPS
		go http.ListenAndServe(":80", http.HandlerFunc(netgo.HTTPSRedirect(host, routeMap)))

		certificates, ok := os.LookupEnv("CERTIFICATE_DIRECTORY")
		if !ok {
			certificates = "certificates"
		}
		log.Println("Certificate Directory:", certificates)

		// Serve HTTPS Requests
		config := &tls.Config{MinVersion: tls.VersionTLS12}
		server := &http.Server{
			Addr:              ":443",
			Handler:           mux,
			TLSConfig:         config,
			ReadTimeout:       time.Hour,
			ReadHeaderTimeout: time.Hour,
			WriteTimeout:      time.Hour,
			IdleTimeout:       time.Hour,
		}
		if err := server.ListenAndServeTLS(path.Join(certificates, "fullchain.pem"), path.Join(certificates, "privkey.pem")); err != nil {
			log.Fatal(err)
		}
	} else {
		// Serve HTTP Requests
		log.Println("HTTP Server Listening on :80")
		if err := http.ListenAndServe(":80", mux); err != nil {
			log.Fatal(err)
		}
	}
}
