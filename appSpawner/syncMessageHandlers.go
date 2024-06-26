package appSpawner

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Node"
	"SystemgeSampleChessServer/topics"
)

func (app *App) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		topics.NEW: app.New,
	}
}

func (app *App) New(client *Node.Node, message *Message.Message) (string, error) {
	id := message.GetPayload()
	err := app.StartClient(client, id)
	if err != nil {
		return "", Error.New("Error starting client "+id, err)
	}
	return id, nil
}
