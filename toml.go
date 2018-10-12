package main

import "time"

type config struct {
	Listen    listen
	Limit     limit
	Cert      cert
	Origin    origin
	TTL       ttl
	Cachefill cachefill
	Stats     stats
	Capacity  capacity
	Log       log
}

type listen struct {
	Scheme string
	Host   string
	Port   int
	Prefix string
}

type limit struct {
	MaxKeyLength int
}

type cert struct {
	KeyPath  string
	CertPath string
}

type origin struct {
	Scheme string
	Host   string
	Port   int
	Prefix string
}

type ttl struct {
	TickDelta               duration
	SetCacheControl         bool
	SetChainlightningExpiry bool
}

type cachefill struct {
	Timeout duration
}

type stats struct {
	TickDelta       duration
	RelevanceWindow duration
}

type capacity struct {
	MB int64
}

type log struct {
	Level         string
	RemoteAddress string
}

type duration time.Duration

func (d *duration) UnmarshalText(b []byte) error {
	dur, err := time.ParseDuration(string(b))
	*d = duration(dur)
	return err
}
