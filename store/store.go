package store

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"

	"github.com/gopherguides/jsonstore/storage"
	"github.com/gopherguides/jsonstore/store/collection"
)

// Store represents the in memory JSON store
type Store struct {
	collections map[string]*collection.Collection
	mu          sync.RWMutex
	storer      storage.Storer
	path        string
}

func New(path string, st storage.Storer) *Store {
	return &Store{
		path:        path,
		storer:      st,
		collections: make(map[string]*collection.Collection),
	}
}

func (s *Store) Open() error {
	// Load the collections
	fi, err := ioutil.ReadDir(s.path)
	if err != nil {
		return err
	}
	for _, f := range fi {
		fn := f.Name()
		if strings.HasSuffix(fn, ".db.json") {
			name := strings.TrimSuffix(fn, ".db.json")
			c, err := collection.New(s.path, name, s.storer)
			if err != nil {
				return fmt.Errorf("failed to load collection %q: %s", name, err)
			}
			s.collections[name] = c
		}
	}
	return nil
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// channel to capture errors
	errCh := make(chan error, len(s.collections))

	// close the collections
	for _, v := range s.collections {
		go func(c *collection.Collection) {
			errCh <- c.Close()
		}(v)
	}
	var err error
	for i := 0; i <= len(errCh); i++ {
		e := <-errCh
		if e != nil {
			err = e
		}
	}
	return err
}

// CreateCollection will create a new collection and return it.
// If the collection already exists, it will return an error
func (s *Store) CreateCollection(name string) (*collection.Collection, error) {
	// need only a read lock to check existence
	s.mu.RLock()
	if _, ok := s.collections[name]; ok {
		s.mu.RUnlock()
		return nil, fmt.Errorf("collection for %q already exists", name)
	}
	s.mu.RUnlock()

	// get a full write lock
	s.mu.Lock()
	defer s.mu.Unlock()
	c, err := collection.New(s.path, name, s.storer)
	if err != nil {
		return nil, err
	}
	s.collections[name] = c
	return c, nil
}

func (s *Store) Collection(name string) *collection.Collection {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.collections[name]
	if !ok {
		return nil
	}
	return c
}
