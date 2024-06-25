package appWebsocketHTTP

import (
	"Systemge/Client"
	"Systemge/Error"
	"Systemge/Message"
	"SystemgeSampleChessServer/topics"
	"strings"
)

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Client.SyncMessageHandler {
	return map[string]Client.SyncMessageHandler{
		topics.PROPAGATE_GAMESTART: func(client *Client.Client, message *Message.Message) (string, error) {
			gameId := message.GetOrigin()
			ids := strings.Split(gameId, "-")
			err := client.AddToWebsocketGroup(gameId, ids[0])
			if err != nil {
				return "", Error.New("Error adding \""+ids[0]+"\" to group \""+gameId+"\"", err)
			}
			err = client.AddToWebsocketGroup(gameId, ids[1])
			if err != nil {
				err := client.RemoveFromWebsocketGroup(gameId, ids[0])
				if err != nil {
					client.GetLogger().Log(Error.New("Error removing \""+ids[0]+"\" from group \""+gameId+"\"", err).Error())
				}
				return "", Error.New("Error adding \""+ids[1]+"\" to group \""+gameId+"\"", err)
			}
			client.WebsocketGroupcast(gameId, message)
			return "", nil
		},
	}
}
