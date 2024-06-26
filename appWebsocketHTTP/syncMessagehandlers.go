package appWebsocketHTTP

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Node"
	"SystemgeSampleChessServer/topics"
	"strings"
)

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		topics.PROPAGATE_GAMESTART: func(client *Node.Node, message *Message.Message) (string, error) {
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
