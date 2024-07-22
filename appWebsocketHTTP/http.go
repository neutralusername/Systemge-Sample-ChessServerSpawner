package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Http"
	"net/http"
)

func (app *AppWebsocketHTTP) GetHTTPRequestHandlers() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/": Http.SendDirectory("../frontend"),
	}
}

func (app *AppWebsocketHTTP) GetHTTPComponentConfig() *Config.Http {
	return &Config.Http{
		Server: &Config.TcpServer{
			Port: 8080,
		},
	}
}
