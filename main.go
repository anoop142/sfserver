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
		log.Println(r.URL.Path)
	}
	if r.URL.Path != "/"+filepath.Base(cfg.path) {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, cfg.path)
}

func (cfg *config) serveDir(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("/tmp"))
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
		http.Handle("/", http.FileServer(http.Dir(absolutePath)))
		if cfg.verbose {
			log.Printf("Serving %s on %s:%s\n", absolutePath, localIP, cfg.port)
		} else {
			fmt.Printf("%s:%s\n", localIP, cfg.port)
		}
	} else {
		http.HandleFunc("/", cfg.serveFile)
		encodedFileName := url.PathEscape(filepath.Base(cfg.path))

		if cfg.verbose {
			log.Printf("Serving %s:%s/%s\n", localIP, cfg.port, encodedFileName)
		} else {
			fmt.Printf("%s:%s/%s\n", localIP, cfg.port, encodedFileName)
		}
	}

	log.Fatal(http.ListenAndServe(":"+cfg.port, nil))
}
