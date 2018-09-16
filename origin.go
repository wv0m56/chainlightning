package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/wv0m56/fury/engine"
)

var (
	err404     = errors.New("404 from origin")
	err500     = errors.New("500 from origin")
	errUnknown = errors.New("unknown error")
)

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

	if status := resp.StatusCode; status != http.StatusOK {
		cancel()
		if status == http.StatusNotFound {
			return nil, nil, err404
		}

		if status == http.StatusInternalServerError {
			return nil, nil, err500
		} else {
			return nil, nil, errUnknown
		}
	}

	var expiry time.Time
	if maxAge := resp.Header.Get("Cache-Control"); maxAge != "" {
		i, err := strconv.Atoi(strings.TrimPrefix(maxAge, "max-age="))
		expiry = time.Now().Add(time.Duration(i) * time.Second)
		if err != nil {
			cancel()
			return nil, nil, err500
		}
	}

	// TODO: propagate max-age to client

	return &response{resp, cancel}, &expiry, nil
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
