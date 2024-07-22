package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/HTTP"
	"net/http"
)

func (app *AppWebsocketHTTP) GetHTTPMessageHandlers() map[string]http.HandlerFunc {
	return map[string]http.HandlerFunc{
		"/": HTTP.SendDirectory("../frontend"),
	}
}

func (app *AppWebsocketHTTP) GetHTTPComponentConfig() *Config.HTTP {
	return &Config.HTTP{
		Server: &Config.TcpServer{
			Port: 8080,
		},
	}
}
