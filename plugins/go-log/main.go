package main

import (
	"time"
	"os"
	"log"
	"fmt"
  	"path/filepath"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
)

var Version = "0.0.1"

var Priority = 1000

type Config struct {
}

func New() interface{} {
	return &Config{}
}

func (cfg Config) Access(kong *pdk.PDK) {
	now := time.Now()
	host, _ := kong.Request.GetHost()
	method, _ := kong.Request.GetMethod()
	path, _ := kong.Request.GetPath()

	filePath := "/go-log/log.txt"

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		kong.Log.Err(fmt.Sprintf("failed to create directory: %s", err))
		log.Fatalf("failed to create directory: %s", err)
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		kong.Log.Err(fmt.Sprintf("failed opening file: %s", err))
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	if _, err := file.WriteString(fmt.Sprintf("[%s] %s - %s~%s.\n", now, host, method, path)); err != nil {
		kong.Log.Err(fmt.Sprintf("failed writing to file: %s", err))
		log.Fatalf("failed writing to file: %s", err)
	}
}

func main() {
	server.StartServer(New, Version, Priority)
}