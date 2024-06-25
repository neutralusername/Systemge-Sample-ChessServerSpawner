package appWebsocketHTTP

import (
	"Systemge/Client"
	"Systemge/Error"
	"Systemge/Message"
	"SystemgeSampleChessServer/topics"
	"strings"
)

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Client.AsyncMessageHandler {
	return map[string]Client.AsyncMessageHandler{
		topics.PROPAGATE_GAMEEND: func(client *Client.Client, message *Message.Message) error {
			gameId := message.GetOrigin()
			ids := strings.Split(gameId, "-")
			client.WebsocketGroupcast(message.GetOrigin(), message)
			err := client.RemoveFromWebsocketGroup(gameId, ids[0])
			if err != nil {
				client.GetLogger().Log(Error.New("Error removing \""+ids[0]+"\" from group \""+gameId+"\"", err).Error())
			}
			err = client.RemoveFromWebsocketGroup(gameId, ids[1])
			if err != nil {
				client.GetLogger().Log(Error.New("Error removing \""+ids[1]+"\" from group \""+gameId+"\"", err).Error())
			}
			app.mutex.Lock()
			delete(app.clientGameIds, ids[0])
			delete(app.clientGameIds, ids[1])
			app.mutex.Unlock()
			return nil
		},
	}
}
