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
	}

	r.Get(prefix+"{key}", func(w http.ResponseWriter, r *http.Request) {
		key := chi.URLParam(r, "key")
		var status int
		var err error
		defer func() {
			if c.Log.Level == "verbose" {
				info := ""

				if err == nil {
					info += fmt.Sprintf("[%v] ", http.StatusOK)
				} else {
					info += fmt.Sprintf("[%v] ", status)
				}

				var addr string
				if c.Log.RemoteAddress == "RemoteAddr" {
					addr = r.RemoteAddr
				} else if c.Log.RemoteAddress == "X-Forwarded-For" {
					addr = r.Header.Get("X-Forwarded-For")
				}

				info += fmt.Sprintf("[%v] [%v]\n", addr, key)

				logInfo(info)
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
