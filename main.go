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
	"path/filepath"
)

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

func main() {
	flag.Usage = func() {
		fmt.Println("sfserver - A simple static file server by Anoop")
		flag.PrintDefaults()
	}

	port := flag.String("p", "8000", "port to serve on")
	directory := flag.String("d", ".", "the directory to serve")
	flag.Parse()

	absolutePath, _ := filepath.Abs(*directory)
	localIP := getLocalIP()

	http.Handle("/", http.FileServer(http.Dir(absolutePath)))
	log.Printf("Serving %s on %s:%s\n", absolutePath, localIP, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
