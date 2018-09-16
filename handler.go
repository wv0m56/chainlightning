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
		var addr string
		if c.Log.RemoteAddress == "RemoteAddr" {
			addr = r.RemoteAddr
		} else if c.Log.RemoteAddress == "X-Forwarded-For" {
			addr = r.Header.Get("X-Forwarded-For")
		}
		defer func() {

			if err != nil {

				go l.logRequestErr(fmt.Sprintf("error: %v. Request info: [%v] [%v] [%v]\n",
					err, status, addr, key))
				w.WriteHeader(status)

			} else {

				if c.Log.Level == "verbose" {
					status = http.StatusOK
					go l.logInfo(fmt.Sprintf("[%v] [%v] [%v]\n", status, addr, key))
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

		_, err = io.Copy(w, data)
		if err != nil {
			status = http.StatusInternalServerError
			return
		}
	})
}
