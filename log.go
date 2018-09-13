package main

import (
	"fmt"
	"os"
)

func logErr(err error) {
	fmt.Fprintln(os.Stderr, err)
}

func logInfo(info string) {
	fmt.Fprintln(os.Stdout, info)
}
