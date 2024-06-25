package appChess

import (
	"Systemge/Client"
)

func (app *App) GetAsyncMessageHandlers() map[string]Client.AsyncMessageHandler {
	return map[string]Client.AsyncMessageHandler{}
}
