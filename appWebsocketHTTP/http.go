package appWebsocketHTTP

import (
	"Systemge/Client"
)

func (app *AppWebsocketHTTP) GetHTTPRequestHandlers() map[string]Client.HTTPRequestHandler {
	return map[string]Client.HTTPRequestHandler{
		"/": Client.SendDirectory("../frontend"),
	}
}
