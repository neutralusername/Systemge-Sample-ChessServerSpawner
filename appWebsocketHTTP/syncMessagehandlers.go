package appWebsocketHTTP

import (
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/topics"
	"strings"
)

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Client.SyncMessageHandler {
	return map[string]Client.SyncMessageHandler{
		topics.PROPAGATE_GAMESTART: func(client *Client.Client, message *Message.Message) (string, error) {
			gameId := message.GetOrigin()
			ids := strings.Split(gameId, "-")
			err := client.AddToGroup(gameId, ids[0])
			if err != nil {
				return "", Utilities.NewError("Error adding \""+ids[0]+"\" to group \""+gameId+"\"", err)
			}
			err = client.AddToGroup(gameId, ids[1])
			if err != nil {
				err := client.RemoveFromGroup(gameId, ids[0])
				if err != nil {
					client.GetLogger().Log(Utilities.NewError("Error removing \""+ids[0]+"\" from group \""+gameId+"\"", err).Error())
				}
				return "", Utilities.NewError("Error adding \""+ids[1]+"\" to group \""+gameId+"\"", err)
			}
			client.Groupcast(gameId, message)
			return "", nil
		},
	}
}
