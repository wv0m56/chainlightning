package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/wv0m56/fury/engine"
)

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

			if err != nil {

				go l.logRequestErr(fmt.Sprintf("error: %v. Request info: [%v] [%v] [%v] in %v\n",
					err, status, addr, key, time.Since(start)))
				w.WriteHeader(status)

			} else {

				if c.Log.Level == "always" {
					status = http.StatusOK
					go l.logInfo(fmt.Sprintf("[%v] [%v] [%v] in %v\n", status, addr, key, time.Since(start)))
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
			w.Header().Set("Cache-Control", "max-age="+strconv.FormatFloat(ttl[0], 'f', 0, 64))
		}

		_, err = io.Copy(w, data)
		if err != nil {
			status = http.StatusInternalServerError
			return
		}
	})
}
