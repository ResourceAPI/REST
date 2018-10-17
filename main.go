package main

import (
	"github.com/StratoAPI/Interface/plugins"
	"github.com/StratoAPI/REST/server"
)

type RESTPlugin string

func (RESTPlugin) Name() string {
	return "REST"
}

func (RESTPlugin) Entrypoint() {
	plugins.GetRegistry().RegisterFacade("REST", &server.RESTFacade{})
}

var CorePlugin RESTPlugin
