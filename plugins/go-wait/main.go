package main

import (
	"fmt"
	"time"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

var Version = "0.0.1"

var Priority = 1000

type Config struct {
	WaitTime int
}

func New() interface{} {
	return &Config{}
}

var requests = make(map[string]time.Time)

func (cfg Config) Access(kong *pdk.PDK) {
	kong.Response.SetHeader("x-wait-time", fmt.Sprintf("%d seconds", cfg.WaitTime))
	host, _ := kong.Request.GetHost()
	lastRequest, exists := requests[host]
	if exists && time.Now().Sub(lastRequest) < time.Duration(cfg.WaitTime)*time.Second {
		kong.Response.Exit(400, []byte("Maximum requests reached"), make(map[string][]string))
	} else {
		requests[host] = time.Now()
	}
}

func main() {
	server.StartServer(New, Version, Priority)
}
