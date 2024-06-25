package appSpawner

import (
	"Systemge/Client"
	"Systemge/Error"
	"Systemge/Message"
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
		return "", Error.New("Error starting client "+id, err)
	}
	return id, nil
}
