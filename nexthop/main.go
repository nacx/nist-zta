// Copyright (c) Tetrate, Inc 2021 All Rights Reserved.

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	nexthop string
	address string
)

func main() {
	host := os.Getenv("HOSTNAME")

	flag.StringVar(&nexthop, "next-hop", "", "The next hop")
	flag.StringVar(&address, "listen-address", ":8000", "Address where the HTTP server listens")
	flag.Parse()

	handler := func(w http.ResponseWriter, req *http.Request) {
		_, err := w.Write([]byte(host + "\n"))
		if err != nil {
			sendError(w, err)
			return
		}

		if nexthop != "" {
			res, err := http.Get(nexthop)
			if err != nil {
				sendError(w, err)
				return
			}

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				sendError(w, err)
				return
			}

			_, err = w.Write([]byte("  -> "))
			if err != nil {
				sendError(w, err)
				return
			}

			_, err = w.Write(b)
			if err != nil {
				sendError(w, err)
				return
			}
		}
	}

	log.Printf("starting server at: %s", address)
	log.Printf("hostname: %s", host)
	log.Printf("next hop: %s ", nexthop)

	go func() { log.Fatal(http.ListenAndServe(address, http.HandlerFunc(handler))) }()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done

	log.Print("shutting down")
}

func sendError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusServiceUnavailable)
	_, _ = w.Write([]byte(err.Error()))
}
