package appSpawner

import (
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/topics"
)

func (app *App) GetSyncMessageHandlers() map[string]Client.SyncMessageHandler {
	return map[string]Client.SyncMessageHandler{
		topics.NEW: app.New,
	}
}

func (app *App) New(client *Client.Client, message *Message.Message) (string, error) {
	id := message.GetPayload()
	err := app.StartClient(client, id)
	if err != nil {
		return "", Utilities.NewError("Error starting client "+id, err)
	}
	return id, nil
}
