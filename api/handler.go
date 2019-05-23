package api

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gopherguides/jsonstore/store/collection"
	"github.com/gorilla/mux"
)

type Handler struct {
	Store interface {
		Open() error
		CreateCollection(string) (*collection.Collection, error)
		Collection(string) *collection.Collection
	}
	debug bool
	mux   *mux.Router
}

func NewHandler(debug bool) *Handler {
	h := &Handler{
		debug: debug,
	}
	m := mux.NewRouter()
	m.HandleFunc("/collections/{collection}/{id}", h.Read).Methods("GET")
	m.HandleFunc("/collections/{collection}/{id}", h.Create).Methods("PUT") // Update is dumb right now, just overwrites entire object
	m.HandleFunc("/collections/{collection}/{id}", h.Create).Methods("POST")
	m.HandleFunc("/collections", h.CreateCollection).Methods("POST")
	h.mux = m

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logging(h.mux).ServeHTTP(w, r)
}

func (h *Handler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	if h.debug {
		log.Println("create collection")
	}

	data := struct {
		Name string `json:"name"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&data)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	_, err = h.Store.CreateCollection(data.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	if h.debug {
		log.Println("write document")
	}

	vars := mux.Vars(r)
	name := vars["collection"]
	id := vars["id"]
	c := h.Store.Collection(name)
	if c == nil {
		// collection doesn't exist, lets create it
		if h.debug {
			log.Printf("creating collection %q", name)
		}
		cc, err := h.Store.CreateCollection(name)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		c = cc
	}
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		if h.debug {
			log.Println(err)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if h.debug {
		log.Printf("received json: %s", string(b))
	}
	c.Update(id, string(b))
	w.WriteHeader(http.StatusAccepted)

	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (h *Handler) Read(w http.ResponseWriter, r *http.Request) {
	if h.debug {
		log.Println("read document")
	}
	vars := mux.Vars(r)
	name := vars["collection"]
	id := vars["id"]
	c := h.Store.Collection(name)
	if c == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	doc, ok := c.Query(id)
	if !ok {
		w.WriteHeader(http.StatusNotFound)
	}
	w.WriteHeader(http.StatusOK)
	if h.debug {
		log.Printf("writing out json: %s", string(doc))
	}
	w.Write([]byte(doc))
}
