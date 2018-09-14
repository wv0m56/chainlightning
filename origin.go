package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/wv0m56/fury/engine"
)

var errNotFound = errors.New("404 from origin")

func createOrigin(c *config) (engine.Origin, error) {
	port := strconv.Itoa(c.Origin.Port)

	var prefix string
	if p := c.Origin.Prefix; p != "" {
		prefix = p
	}
	if prefix != "" {
		prefix = "/" + prefix + "/"
	}

	urlPrefix := c.Origin.Scheme + "://" + c.Origin.Host + ":" + port + "/" + prefix
	_, err := url.ParseRequestURI(urlPrefix)
	if err != nil {
		return nil, err
	}

	return &backend{urlPrefix}, nil
}

type backend struct {
	urlPrefix string
}

func (b *backend) Fetch(key string, timeout time.Duration) (
	io.ReadCloser, *time.Time, error) {

	// TODO: timeout context, expiry

	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	req, err := http.NewRequest("GET", b.urlPrefix+key, nil)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	req = req.WithContext(ctx)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		cancel()
		return nil, nil, err
	}

	return &response{resp, cancel}, nil, nil // TODO: expiry
}

type response struct {
	resp   *http.Response
	cancel context.CancelFunc
}

func (r *response) Close() error {
	r.cancel()
	return r.resp.Body.Close()
}

func (r *response) Read(p []byte) (int, error) {
	return r.resp.Body.Read(p)
}
