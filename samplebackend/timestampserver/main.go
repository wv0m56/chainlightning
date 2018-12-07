package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

func main() {
	r := chi.NewRouter()
	r.Get("/timeout/{ms}", func(w http.ResponseWriter, r *http.Request) {
		msParam := chi.URLParam(r, "ms")
		ms, err := strconv.ParseInt(msParam, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("Chainlighting-Expiry", msParam)
		sec := (time.Duration(ms) * time.Millisecond).Round(time.Second)
		secStr := strconv.FormatFloat(sec.Seconds(), 'f', 0, 64)
		w.Header().Set("Cache-Control", "max-age="+secStr)

		respBodyStr := "Time now is " + time.Now().String()
		respBodyStr += "\nCheck header for expiry values"
		_, err = w.Write([]byte(respBodyStr))
		if err != nil {
			panic(err)
		}
	})

	http.ListenAndServe(":12000", r)
	fmt.Println("server stopped")
}
