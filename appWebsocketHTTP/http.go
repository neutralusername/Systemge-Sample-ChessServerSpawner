package appWebsocketHTTP

import (
	"Systemge/Application"
	"Systemge/HTTP"
)

func (app *AppWebsocketHTTP) GetHTTPRequestHandlers() map[string]Application.HTTPRequestHandler {
	return map[string]Application.HTTPRequestHandler{
		"/": HTTP.SendDirectory("../frontend"),
	}
}
