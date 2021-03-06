package main

import (
	"errors"
	"flag"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/go-chi/chi"
	"github.com/tylerb/graceful"
	"github.com/wv0m56/fury/engine"
)

var (
	confPath string
	l        *logger
)

func init() {
	flag.StringVar(&confPath, "c", "config.toml", "path to toml config file")
	l = &logger{sync.Mutex{}, sync.Mutex{}}
}

func main() {
	flag.Parse()

	var conf config
	_, err := toml.DecodeFile(confPath, &conf)
	if err != nil {
		l.logSimpleErr(err)
		return
	}

	err = validateConfig(&conf)
	if err != nil {
		l.logSimpleErr(err)
		return
	}

	opts := engine.Options{}
	opts.AccessStatsRelevanceWindow = time.Duration(conf.Stats.RelevanceWindow)
	opts.AccessStatsTickStep = time.Duration(conf.Stats.TickDelta)
	opts.CacheFillTimeout = time.Duration(conf.Cachefill.Timeout)
	opts.ExpectedLen = 100 * 1000 * 1000
	opts.MaxPayloadTotalBytes = conf.Capacity.MB * 1000 * 1000
	opts.TTLTickStep = time.Duration(conf.TTL.TickDelta)
	opts.O, err = createOrigin(&conf)
	if err != nil {
		l.logSimpleErr(err)
		return
	}

	e, err := engine.NewEngine(&opts)
	if err != nil {
		l.logSimpleErr(err)
		return
	}

	r := chi.NewRouter()
	routeHttp(e, &conf, r)

	if conf.Listen.Scheme == "http" {

		graceful.Run(conf.Listen.Host+":"+strconv.Itoa(conf.Listen.Port), 1*time.Second, r)

	} else if conf.Listen.Scheme == "https" {

		gracefulServer := &graceful.Server{
			Timeout: 1 * time.Second,
		}
		err = gracefulServer.ListenAndServeTLS(
			conf.Cert.CertPath,
			conf.Cert.KeyPath,
		)
		if err != nil {
			panic(err)
		}

	} else {

		panic("unknown listen scheme")
	}
}

func validateConfig(c *config) error {

	if sch := c.Listen.Scheme; sch != "http" && sch != "https" {
		return errors.New("[listen]Scheme must be http or https")
	}

	if sch := c.Origin.Scheme; sch != "http" && sch != "https" {
		return errors.New("[origin]Scheme must be http or https")
	}

	var host string
	if c.Listen.Host == "*" {
		host = ""
	} else {
		host = c.Listen.Host
	}
	if _, err := url.ParseRequestURI(
		c.Listen.Scheme +
			"://" +
			host +
			":" +
			strconv.Itoa(c.Listen.Port) +
			"/" +
			c.Listen.Prefix); err != nil {

		return errors.New("[listen] parameters do not form valid url")
	}

	if c.Limit.MaxKeyLength < 8 {
		return errors.New("[limit]MaxKeyLength must be 8 or more")
	}

	if _, err := url.ParseRequestURI(
		c.Origin.Scheme +
			"://" +
			c.Origin.Host +
			":" +
			strconv.Itoa(c.Origin.Port) +
			"/" +
			c.Origin.Prefix); err != nil {

		return errors.New("[origin] parameters do not form valid url")
	}

	if level := c.Log.Level; level != "error" && level != "always" {
		return errors.New("[log]Level must be always or error")
	}

	if addr := c.Log.RemoteAddress; addr != "RemoteAddr" && addr != "X-Forwarded-For" {
		return errors.New(`[log]RemoteAddress must be "RemoteAddr" or "X-Forwarded-For"`)
	}

	return nil
}
