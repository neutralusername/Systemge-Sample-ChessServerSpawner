package appWebsocketHTTP

import (
	"Systemge/Application"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/topics"
	"strings"
)

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{
		topics.PROPAGATE_GAMESTART: func(message *Message.Message) (string, error) {
			gameId := message.GetOrigin()
			ids := strings.Split(gameId, "-")
			err := app.client.GetWebsocketServer().AddToGroup(gameId, ids[0])
			if err != nil {
				return "", Utilities.NewError("Error adding \""+ids[0]+"\" to group \""+gameId+"\"", err)
			}
			err = app.client.GetWebsocketServer().AddToGroup(gameId, ids[1])
			if err != nil {
				err := app.client.GetWebsocketServer().RemoveFromGroup(gameId, ids[0])
				if err != nil {
					app.client.GetLogger().Log(Utilities.NewError("Error removing \""+ids[0]+"\" from group \""+gameId+"\"", err).Error())
				}
				return "", Utilities.NewError("Error adding \""+ids[1]+"\" to group \""+gameId+"\"", err)
			}
			app.client.GetWebsocketServer().Groupcast(gameId, message)
			return "", nil
		},
	}
}
