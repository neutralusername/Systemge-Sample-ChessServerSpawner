package appSpawner

import (
	"Systemge/Application"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/topics"
)

func (app *App) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{
		topics.END: app.End,
	}
}

func (app *App) End(message *Message.Message) error {
	id := message.GetPayload()
	err := app.EndClient(id)
	if err != nil {
		return Utilities.NewError("Error ending client "+id, err)
	}
	return nil
}
