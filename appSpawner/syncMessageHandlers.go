package appSpawner

import (
	"Systemge/Application"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/topics"
)

func (app *App) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{
		topics.NEW: app.New,
	}
}

func (app *App) New(message *Message.Message) (string, error) {
	id := message.GetPayload()
	err := app.StartClient(id)
	if err != nil {
		return "", Utilities.NewError("Error starting client "+id, err)
	}
	return id, nil
}
