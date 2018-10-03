package main

import (
	"github.com/ResourceAPI/Core/plugins"
	"github.com/ResourceAPI/REST/server"
)

type RESTPlugin string

func (RESTPlugin) Name() string {
	return "REST"
}

func (RESTPlugin) Entrypoint() {
	plugins.RegisterFacade("rest", &server.RESTFacade{})
}

var CorePlugin RESTPlugin
