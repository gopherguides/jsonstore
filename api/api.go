package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gopherguides/jsonstore/storage"
	"github.com/gopherguides/jsonstore/store"
)

type API struct {
	handler *Handler
	store   *store.Store
	srv     *http.Server
	debug   bool
}

func New(addr string, path string, debug bool) (*API, error) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return nil, err
	}
	fs := storage.NewFileSystem()
	// Create our store instance
	s := store.New(path, fs)

	h := NewHandler(debug)
	h.Store = s

	// create our server
	srv := &http.Server{
		Addr:           addr,
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return &API{
		debug:   debug,
		handler: h,
		srv:     srv,
		store:   s,
	}, nil
}

func (a *API) Open() error {
	// Open the store
	if err := a.handler.Store.Open(); err != nil {
		log.Fatal(err)
	}

	log.Printf("starting api service: %s", a.srv.Addr)
	if err := a.srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		return err
	}
	return nil
}

func (a *API) Close() error {
	if err := a.srv.Shutdown(context.Background()); err != nil {
		// Error from closing listeners, or context timeout:
		return err
	}
	if err := a.store.Close(); err != nil {
		return err
	}

	log.Println("successfully shut down api")
	return nil
}
