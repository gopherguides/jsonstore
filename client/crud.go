package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"reflect"
	"strings"
)

func collectionName(v interface{}) string {
	s := strings.ToLower(fmt.Sprintf("%s", reflect.TypeOf(v)))
	f := strings.Split(s, ".")
	if len(f) == 0 {
		return pluralize(s)
	}
	return pluralize(f[len(f)-1])
}

// pluraize is very simplistic
func pluralize(s string) string {
	if strings.HasSuffix(s, "s") {
		return s
	}
	return s + "s"
}

func (c *Client) Create(id string, v interface{}) error {
	collection := collectionName(v)

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	resp, err := c.Post(path.Join(collection, id), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &statusCodeError{Code: resp.StatusCode}
	}

	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

func (c *Client) Update(id string, v interface{}) error {
	collection := collectionName(v)

	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	resp, err := c.Put(path.Join(collection, id), bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &statusCodeError{Code: resp.StatusCode}
	}

	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

func (c *Client) Read(id string, v interface{}) error {
	collection := collectionName(v)
	resp, err := c.Get(path.Join(collection, id))
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &statusCodeError{Code: resp.StatusCode}
	}

	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

func (c *Client) Remove(id string, v interface{}) error {
	collection := collectionName(v)
	resp, err := c.Delete(path.Join(collection, id))
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &statusCodeError{Code: resp.StatusCode}
	}

	return err
}
