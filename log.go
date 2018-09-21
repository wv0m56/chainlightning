package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type logger struct {
	outMu sync.Mutex // guards stdout
	errMu sync.Mutex // guards stderr
}

func (l *logger) logSimpleErr(err error) {
	l.errMu.Lock()
	fmt.Fprintln(os.Stderr, err)
	l.errMu.Unlock()
}

func (l *logger) logRequestErr(err error, status int, addr string, key string, dur time.Duration) {
	l.errMu.Lock()
	fmt.Fprintf(os.Stderr, "error: %v. request info: [%v] [%v] [%v] in %v\n",
		err, status, addr, key, dur)
	l.errMu.Unlock()
}

func (l *logger) logPanic(addr, key string, rec interface{}) {
	l.errMu.Lock()
	fmt.Fprintf(os.Stderr, "panicked. request info: [%v] [%v]\nstacktrace:\n%v", addr, key, rec)
	l.errMu.Unlock()
}

func (l *logger) logInfo(status int, addr string, key string, dur time.Duration) {
	l.outMu.Lock()
	fmt.Fprintf(os.Stdout, "[%v] [%v] [%v] in %v\n", status, addr, key, dur)
	l.outMu.Unlock()
}
