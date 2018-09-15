package main

import (
	"fmt"
	"os"
	"sync"
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

func (l *logger) logRequestErr(errStr string) {
	l.errMu.Lock()
	fmt.Fprintln(os.Stderr, errStr)
	l.errMu.Unlock()
}

func (l *logger) logInfo(info string) {
	l.outMu.Lock()
	fmt.Fprintln(os.Stdout, info)
	l.outMu.Unlock()
}
