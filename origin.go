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
	var err error
	defer func() {
		if err != nil {
			cancel()
		}
	}()

	req, err := http.NewRequest("GET", b.urlPrefix+key, nil)
	if err != nil {
		return nil, nil, err
	}

	req = req.WithContext(ctx)
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if err != nil {
			_ = resp.Body.Close()
		}
	}()

	if status := resp.StatusCode; status != http.StatusOK {
		if status == http.StatusNotFound {
			return nil, nil, err404
		}

		if status == http.StatusInternalServerError {
			return nil, nil, err500
		} else {
			return nil, nil, errUnknown
		}
	}

	exp, err := getExpiry(resp)
	if err != nil {
		return nil, nil, err
	}

	return &response{resp, cancel}, exp, nil
}

func getExpiry(resp *http.Response) (*time.Time, error) {
	if ce := resp.Header.Get("Chainlightning-Expiry"); ce != "" {

		exp, err := time.Parse(time.RFC3339, ce)
		if err != nil {
			return nil, err
		}
		return &exp, nil

	} else if cc := resp.Header.Get("Cache-Control"); cc != "" {

		i, err := strconv.Atoi(strings.TrimPrefix(cc, "max-age="))
		if err != nil {
			return nil, err
		}
		exp := time.Now().Add(time.Duration(i) * time.Second)
		return &exp, nil
	}

	return nil, nil
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
