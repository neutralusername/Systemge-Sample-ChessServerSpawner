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

func (app *App) New(node *Node.Node, message *Message.Message) (string, error) {
	id := message.GetPayload()
	err := app.StartNode(node, id)
	if err != nil {
		return "", Error.New("Error starting node "+id, err)
	}
	return id, nil
}
