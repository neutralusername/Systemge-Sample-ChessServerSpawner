package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Node"
)

func (app *AppWebsocketHTTP) GetHTTPRequestHandlers() map[string]Node.HTTPRequestHandler {
	return map[string]Node.HTTPRequestHandler{
		"/": Node.SendDirectory("../frontend"),
	}
}

func (app *AppWebsocketHTTP) GetHTTPComponentConfig() Config.HTTP {
	return Config.HTTP{
		Port:        ":8080",
		TlsCertPath: "",
		TlsKeyPath:  "",
	}
}
