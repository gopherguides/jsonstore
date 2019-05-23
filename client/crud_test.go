package client_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gopherguides/jsonstore/client"
)

type Foo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestRead(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"id":1,"name":"Bar"}`)
	}))
	defer ts.Close()
	c, err := client.New(client.Config{Host: ts.URL})
	if err != nil {
		t.Fatal(err)
	}

	f := &Foo{}

	err = c.Read("1", f)
	if err != nil {
		t.Fatal(err)
	}
	if got, exp := f.ID, 1; got != exp {
		t.Errorf("unexpected ID: got %d, exp %d", got, exp)
	}

	if got, exp := f.Name, "Bar"; got != exp {
		t.Errorf("unexpected name: got %s, exp %s", got, exp)
	}
}

func TestCreate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"id":2,"name":"Baz"}`)
	}))
	defer ts.Close()
	c, err := client.New(client.Config{Host: ts.URL})
	if err != nil {
		t.Fatal(err)
	}

	f := &Foo{ID: 1, Name: "Bar"}

	err = c.Create("1", f)
	if err != nil {
		t.Fatal(err)
	}
	if got, exp := f.ID, 2; got != exp {
		t.Errorf("unexpected ID: got %d, exp %d", got, exp)
	}

	if got, exp := f.Name, "Baz"; got != exp {
		t.Errorf("unexpected name: got %s, exp %s", got, exp)
	}
}
