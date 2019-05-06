package collection

import (
	"encoding/json"
	"log"
	"path/filepath"
	"sync"

	"github.com/gopherguides/jsonstore/storage"
)

// Collection represents a table of json data
type Collection struct {
	Name      string `json:"name"`
	mu        sync.RWMutex
	Documents Documents `json:"-"`
	storer    storage.Storer
	path      string
}

func New(path string, name string, s storage.Storer) (*Collection, error) {
	c := &Collection{
		Name:      name,
		Documents: make(Documents),
		path:      filepath.Join(path, name+".db.json"),
		storer:    s,
	}

	if err := c.load(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *Collection) Update(id string, document string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Documents[id] = document
}

func (c *Collection) Query(id string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	d, ok := c.Documents[id]
	return d, ok
}

func (c *Collection) load() error {
	r, err := c.storer.Reader(c.path)
	if err != nil {
		return err
	}
	if r == nil {
		// no file to load
		return nil
	}
	log.Printf("loading collection %q", c.Name)
	return json.NewDecoder(r).Decode(&c.Documents)
}

func (c *Collection) Close() error {
	log.Printf("writing collection %q to storage", c.Name)
	c.mu.Lock()
	defer c.mu.Unlock()
	w, err := c.storer.Writer(c.path)
	if err != nil {
		return err
	}
	defer w.Close()
	return json.NewEncoder(w).Encode(c.Documents)
}
