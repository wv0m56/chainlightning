package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestExpiry(t *testing.T) {
	resp1 := &http.Response{
		Header: http.Header(map[string][]string{}),
	}
	resp1.Header.Set("Cache-Control", "private, max-age=555")
	exp1, err := getExpiry(resp1)
	assert.Nil(t, err)
	assert.True(t, roughly(555, time.Until(*exp1).Seconds()))

	// double header
	resp2 := &http.Response{
		Header: http.Header(map[string][]string{}),
	}
	resp2.Header.Set("Cache-Control", "private, max-age=555")

	b, err := time.Now().Add(500 * time.Millisecond).MarshalText()
	assert.Nil(t, err)

	resp2.Header.Set("Chainlightning-Expiry", string(b))
	exp2, err := getExpiry(resp2)
	assert.Nil(t, err)
	assert.True(t, roughly(0.5, time.Until(*exp2).Seconds()))
}

func roughly(a, b float64) bool {
	if (a/b > 0.999 && a/b <= 1.0) || (a/b >= 1.0 && a/b < 1.001) {
		return true
	}
	return false
}
