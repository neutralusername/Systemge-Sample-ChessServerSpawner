package appWebsocketHTTP

import (
	"Systemge/Node"
)

func (app *AppWebsocketHTTP) GetHTTPRequestHandlers() map[string]Node.HTTPRequestHandler {
	return map[string]Node.HTTPRequestHandler{
		"/": Node.SendDirectory("../frontend"),
	}
}
