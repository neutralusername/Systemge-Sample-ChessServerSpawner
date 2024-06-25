package appSpawner

import (
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/topics"
)

func (app *App) GetAsyncMessageHandlers() map[string]Client.AsyncMessageHandler {
	return map[string]Client.AsyncMessageHandler{
		topics.END: app.End,
	}
}

func (app *App) End(client *Client.Client, message *Message.Message) error {
	id := message.GetPayload()
	err := app.EndClient(client, id)
	if err != nil {
		return Utilities.NewError("Error ending client "+id, err)
	}
	return nil
}
