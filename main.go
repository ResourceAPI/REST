package main

import (
	"github.com/ResourceAPI/Interface/plugins"
	"github.com/ResourceAPI/REST/server"
)

type RESTPlugin string

func (RESTPlugin) Name() string {
	return "REST"
}

func (RESTPlugin) Entrypoint() {
	plugins.GetRegistry().RegisterFacade("REST", &server.RESTFacade{})
}

var CorePlugin RESTPlugin
