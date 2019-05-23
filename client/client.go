package client

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

func New(cfg Config) (*Client, error) {
	c := &Client{
		config: cfg,
	}

	u, err := url.Parse(c.config.Host)
	if err != nil {
		return nil, fmt.Errorf("invalid host: %s", c.config.Host)
	}
	c.host = u
	return c, nil
}

type Client struct {
	config Config
	host   *url.URL
}

func (c *Client) Get(p string) (*http.Response, error) {
	// copy the host
	u := *c.host
	u.Path = path.Join(u.Path, "collections", p)
	return http.Get(u.String())
}

func (c *Client) Post(p string, body io.Reader) (*http.Response, error) {
	// copy the host
	u := *c.host
	u.Path = path.Join(u.Path, "collections", p)
	return http.Post(u.String(), "application/json", body)
}

func (c *Client) Put(p string, body io.Reader) (*http.Response, error) {
	// copy the host
	u := *c.host
	u.Path = path.Join(u.Path, "collections", p)
	req, err := http.NewRequest("PUT", u.String(), body)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}

func (c *Client) Delete(p string) (*http.Response, error) {
	u := *c.host
	u.Path = path.Join(u.Path, "collections", p)

	req, err := http.NewRequest("DELETE", u.String(), nil)
	if err != nil {
		return nil, err
	}
	return http.DefaultClient.Do(req)
}
