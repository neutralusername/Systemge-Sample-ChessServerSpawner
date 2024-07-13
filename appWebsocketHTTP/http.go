package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Http"
	"Systemge/TcpServer"
)

func (app *AppWebsocketHTTP) GetHTTPRequestHandlers() map[string]Http.RequestHandler {
	return map[string]Http.RequestHandler{
		"/": Http.SendDirectory("../frontend"),
	}
}

func (app *AppWebsocketHTTP) GetHTTPComponentConfig() Config.HTTP {
	return Config.HTTP{
		Server: TcpServer.New(8080, "", ""),
	}
}
