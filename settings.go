package htmlparsing

import (
	"net/http"

	"time"
)

type Settings struct {
	Transport                http.RoundTripper
	Timeout                  time.Duration
	MaxHttpRetries           int
	MaxServerErrorRetries    int
	HttpRetryInterval        time.Duration
	ServerErrorRetryInterval time.Duration
}

func SensibleSettings() *Settings {
	return &Settings{
		Transport:                http.DefaultTransport,
		Timeout:                  60 * time.Second,
		MaxHttpRetries:           3,
		MaxServerErrorRetries:    2,
		HttpRetryInterval:        5 * time.Second,
		ServerErrorRetryInterval: 10 * time.Second,
	}
}
