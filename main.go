/*
A SimpleStatic File Server
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type config struct {
	path    string
	port    string
	verbose bool
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
	}

	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

func (cfg *config) serveFile(w http.ResponseWriter, r *http.Request) {
	if cfg.verbose {
		log.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
	}
	http.ServeFile(w, r, cfg.path)
}

func (cfg *config) serveDir(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cfg.verbose {
			log.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	var cfg config

	flag.Usage = func() {
		fmt.Println("sfserver - A simple static file server by Anoop")
		flag.PrintDefaults()
	}

	flag.StringVar(&cfg.port, "p", "8000", "port to serve on")
	flag.StringVar(&cfg.path, "f", ".", "file or directory to serve")
	flag.BoolVar(&cfg.verbose, "v", false, "verbose")

	flag.Parse()

	absolutePath, _ := filepath.Abs(cfg.path)
	localIP := getLocalIP()

	fileInfo, err := os.Stat(cfg.path)

	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		http.Handle("/", cfg.serveDir(http.FileServer(http.Dir(absolutePath))))
		if cfg.verbose {
			log.Printf("Serving %s on %s:%s", absolutePath, localIP, cfg.port)
		} else {
			fmt.Printf("%s:%s", localIP, cfg.port)
		}
	} else {
		http.HandleFunc("/"+filepath.Base(cfg.path), cfg.serveFile)
		encodedFileName := url.PathEscape(filepath.Base(cfg.path))

		if cfg.verbose {
			log.Printf("Serving %s:%s/%s", localIP, cfg.port, encodedFileName)
		} else {
			fmt.Printf("%s:%s/%s", localIP, cfg.port, encodedFileName)
		}
	}

	log.Fatal(http.ListenAndServe(":"+cfg.port, nil))
}
