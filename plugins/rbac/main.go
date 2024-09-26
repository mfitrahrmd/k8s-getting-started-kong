package main

import (
	"log"

	"github.com/Kong/go-pdk"
	"github.com/Kong/go-pdk/server"
	"github.com/casbin/casbin/v2"
)

var serverErrorMessage = "Internal Server Error"

const PluginName = "RBAC"
const Version = "1.0.0"
const Priority = 1

type Config struct {
}

func main() {
	err := server.StartServer(New, Version, Priority)
	if err != nil {
		log.Fatalf("Failed start %s plugin", PluginName)
	}
}

func New() interface{} {
	return &Config{}
}

func (conf *Config) Access(kong *pdk.PDK) {
	e, err := casbin.NewEnforcer("/etc/kong/casbin/model.conf", "/etc/kong/casbin/policy.csv")
	if err != nil {
		kong.Log.Err(err.Error())
		kong.Response.Exit(500, []byte(serverErrorMessage), nil)

		return
	}

	sub, _ := kong.Request.GetHeader("username")
	obj, _ := kong.Request.GetPath()
	act, _ := kong.Request.GetMethod()

	allowed, err := e.Enforce(sub, obj, act)
	if err != nil {
		kong.Log.Err(err.Error())
		kong.Response.Exit(500, []byte(serverErrorMessage), nil)

		return
	}

	if !allowed {
		kong.Response.Exit(403, []byte("Forbidden"), nil)

		return
	}
}
