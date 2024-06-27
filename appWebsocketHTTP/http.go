package appWebsocketHTTP

import (
	"Systemge/Node"
)

func (app *AppWebsocketHTTP) GetHTTPRequestHandlers() map[string]Node.HTTPRequestHandler {
	return map[string]Node.HTTPRequestHandler{
		"/": Node.SendDirectory("../frontend"),
	}
}

func (app *AppWebsocketHTTP) GetHTTPComponentConfig() Node.HTTPComponentConfig {
	return Node.HTTPComponentConfig{
		Port:        ":8080",
		TlsCertPath: "",
		TlsKeyPath:  "",
	}
}
