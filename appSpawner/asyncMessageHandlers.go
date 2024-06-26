package appSpawner

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Node"
	"SystemgeSampleChessServer/topics"
)

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		topics.END: app.End,
	}
}

func (app *App) End(client *Node.Node, message *Message.Message) error {
	id := message.GetPayload()
	err := app.EndClient(client, id)
	if err != nil {
		return Error.New("Error ending client "+id, err)
	}
	return nil
}
