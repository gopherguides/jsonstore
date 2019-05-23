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
	s := fmt.Sprintf("%s", reflect.TypeOf(v))
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

func (c *Client) Read(id string, v interface{}) error {
	collection := collectionName(v)
	resp, err := c.Get(path.Join(collection, id))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
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

	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}
