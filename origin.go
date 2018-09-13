package main

import (
	"errors"

	"github.com/wv0m56/fury/engine"
)

var errNotFound = errors.New("404 from origin")

func createOrigin(c *config) engine.Origin {
	return nil
}
