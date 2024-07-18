package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Http"
)

func (app *AppWebsocketHTTP) GetHTTPRequestHandlers() map[string]Http.RequestHandler {
	return map[string]Http.RequestHandler{
		"/": Http.SendDirectory("../frontend"),
	}
}

func (app *AppWebsocketHTTP) GetHTTPComponentConfig() Config.HTTP {
	return Config.HTTP{
		Server: Config.TcpServer{
			Port: 8080,
		},
	}
}
