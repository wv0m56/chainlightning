package main

import (
	"flag"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
)

func init() {
	flag.StringVar(&confPath, "c", "config.toml", "path to .toml config file")
}

func main() {
	flag.Parse()

	var conf config
	_, err := toml.DecodeFile(confPath, &conf)
	if err != nil {
		panic(err)
	}
}
