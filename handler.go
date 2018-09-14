package main

import (
	"fmt"
	"io"
	"net/http"

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
		defer func() {
			if c.Log.Level == "verbose" {
				if err == nil {
					status = http.StatusOK
				}

				var addr string
				if c.Log.RemoteAddress == "RemoteAddr" {
					addr = r.RemoteAddr
				} else if c.Log.RemoteAddress == "X-Forwarded-For" {
					addr = r.Header.Get("X-Forwarded-For")
				}

				logInfo(fmt.Sprintf("[%v] [%v] [%v]\n", status, addr, key))
			}
			if err != nil {
				logErr(err)
				w.WriteHeader(status)
			}
		}()

		if len(key) > c.Limit.MaxKeyLength {
			status = http.StatusBadRequest
			return
		}

		data, err := e.Get(key)
		if err != nil {
			if err == errNotFound {
				status = http.StatusNotFound
				return
			}
			status = http.StatusInternalServerError
			return
		}

		_, err = io.Copy(w, data)
		if err != nil {
			status = http.StatusInternalServerError
			return
		}
	})
}
