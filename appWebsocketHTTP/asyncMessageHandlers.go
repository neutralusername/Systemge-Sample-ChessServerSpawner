package appWebsocketHTTP

import (
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/topics"
	"strings"
)

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Client.AsyncMessageHandler {
	return map[string]Client.AsyncMessageHandler{
		topics.PROPAGATE_GAMEEND: func(client *Client.Client, message *Message.Message) error {
			gameId := message.GetOrigin()
			ids := strings.Split(gameId, "-")
			client.WebsocketGroupcast(message.GetOrigin(), message)
			err := client.RemoveFromGroup(gameId, ids[0])
			if err != nil {
				client.GetLogger().Log(Utilities.NewError("Error removing \""+ids[0]+"\" from group \""+gameId+"\"", err).Error())
			}
			err = client.RemoveFromGroup(gameId, ids[1])
			if err != nil {
				client.GetLogger().Log(Utilities.NewError("Error removing \""+ids[1]+"\" from group \""+gameId+"\"", err).Error())
			}
			app.mutex.Lock()
			delete(app.clientGameIds, ids[0])
			delete(app.clientGameIds, ids[1])
			app.mutex.Unlock()
			return nil
		},
	}
}
