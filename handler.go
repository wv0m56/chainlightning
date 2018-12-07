package main

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/wv0m56/fury/engine"
)

/*
(1) /route/to/http/resource - ok
(2) /route/to/http/resource/myimage.jpg - ok
(3) /route/to/http/resource/myimage.jpg?param1=val1&param2=val2 - ok but query ignored

(3) is handled in an identical way as (2)
*/
func routeHttp(e *engine.Engine, c *config, r chi.Router) {
	var prefix string
	if p := c.Listen.Prefix; p != "" {
		prefix = p
	}
	if prefix != "" {
		prefix = "/" + prefix + "/"
	} else {
		prefix = "/"
	}

	r.Get(prefix+"*", func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "*")
		var status int
		var err error
		var addr string
		start := time.Now()
		if c.Log.RemoteAddress == "RemoteAddr" {
			addr = r.RemoteAddr
		} else if c.Log.RemoteAddress == "X-Forwarded-For" {
			addr = r.Header.Get("X-Forwarded-For")
		}
		defer func() {

			if rec := recover(); rec != nil {
				go l.logPanic(addr, key, rec)
				w.WriteHeader(http.StatusInternalServerError)
			}

			if err != nil {

				go l.logRequestErr(err, status, addr, key, time.Since(start))
				w.WriteHeader(status)

			} else {

				if c.Log.Level == "always" {
					status = http.StatusOK
					go l.logInfo(status, addr, key, time.Since(start))
				}
			}
		}()

		if len(key) > c.Limit.MaxKeyLength {
			status = http.StatusBadRequest
			return
		}

		data, err := e.Get(key)
		if err != nil {
			if err == err404 {
				status = http.StatusNotFound
				return
			}
			status = http.StatusInternalServerError
			return
		}

		if ttl := e.GetTTL(key); ttl[0] > 0 {
			if c.TTL.SetCacheControlResponseHeader {
				w.Header().Set("Cache-Control", "max-age="+strconv.FormatFloat(ttl[0], 'f', 0, 64))
			}
			if c.TTL.SetChainlightningExpiryResponseHeader {
				w.Header().Set("Chainlightning-Expiry", strconv.FormatFloat(ttl[0], 'f', 3, 64))
			}
		}

		_, err = io.Copy(w, data)
		if err != nil {
			status = http.StatusInternalServerError
			return
		}
	})
}
