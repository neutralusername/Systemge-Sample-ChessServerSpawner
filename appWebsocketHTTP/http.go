package appWebsocketHTTP

import (
	"Systemge/Application"
	"Systemge/HTTPServer"
)

func (app *AppWebsocketHTTP) GetHTTPRequestHandlers() map[string]Application.HTTPRequestHandler {
	return map[string]Application.HTTPRequestHandler{
		"/": HTTPServer.SendDirectory("../frontend"),
	}
}
