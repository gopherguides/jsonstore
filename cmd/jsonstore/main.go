package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gopherguides/jsonstore/api"
)

func main() {
	// get working directory
	dir, _ := os.Getwd()
	if dir == "" {
		dir = "/"
	}
	// load any arguments
	var path string = filepath.Join(dir, ".jsonstore")
	var addr string = "localhost:9090"
	flag.StringVar(&path, "path", path, "database path")
	flag.StringVar(&addr, "addr", addr, "address to start up api service")
	flag.Parse()

	// Creat our API
	a, err := api.New(addr, path)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := a.Open(); err != nil {
			log.Fatal(err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)
	log.Println("starting up")

	// Block until one of the signals above is received
	<-signalCh
	log.Println("shutdown Signal received, initializing clean shutdown...")

	// Close the store.  Create a channel to capture any error.
	closeCh := make(chan error)
	go func() {
		closeCh <- a.Close()
	}()

	// Block again until another signal is received, a shutdown timeout elapses,
	// or the Command is gracefully closed
	log.Println("waiting for clean shutdown...")
	select {
	case <-signalCh:
		log.Fatal("second signal received, initializing hard shutdown")
	case <-time.After(time.Second * 10):
		log.Fatal("time limit reached, initializing hard shutdown")
	case err := <-closeCh:
		if err == nil {
			log.Println("clean server shutdown completed")
			return
		}
		log.Fatalf("failed to shut down server cleanly: %v", err)
	}
}
