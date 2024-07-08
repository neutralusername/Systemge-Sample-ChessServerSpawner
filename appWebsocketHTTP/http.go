package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Node"
	"Systemge/TcpServer"
)

func (app *AppWebsocketHTTP) GetHTTPRequestHandlers() map[string]Node.HTTPRequestHandler {
	return map[string]Node.HTTPRequestHandler{
		"/": Node.SendDirectory("../frontend"),
	}
}

func (app *AppWebsocketHTTP) GetHTTPComponentConfig() Config.HTTP {
	return Config.HTTP{
		Server: TcpServer.New(8080, "", ""),
	}
}
