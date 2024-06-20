package appWebsocketHTTP

import (
	"Systemge/Application"
	"Systemge/Message"
	"Systemge/Utilities"
	"Systemge/WebsocketClient"
	"SystemgeSampleChessServer/topics"
	"strings"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Application.WebsocketMessageHandler {
	return map[string]Application.WebsocketMessageHandler{
		"startGame": func(client *WebsocketClient.Client, message *Message.Message) error {
			app.mutex.Lock()
			defer app.mutex.Unlock()
			whiteId := client.GetId()
			blackId := message.GetPayload()
			if !app.client.GetWebsocketServer().ClientExists(blackId) {
				return Utilities.NewError("Opponent does not exist", nil)
			}
			if blackId == whiteId {
				return Utilities.NewError("You cannot play against yourself", nil)
			}
			if app.clientGameIds[whiteId] != "" {
				return Utilities.NewError("You are already in a game", nil)
			}
			if app.clientGameIds[blackId] != "" {
				return Utilities.NewError("Opponent is already in a game", nil)
			}
			gameId := whiteId + "-" + blackId
			_, err := app.client.SyncMessage(topics.NEW, app.client.GetName(), gameId)
			if err != nil {
				return Utilities.NewError("Error spawning new game client", err)
			}
			app.clientGameIds[whiteId] = gameId
			app.clientGameIds[blackId] = gameId
			return nil
		},
		"endGame": func(client *WebsocketClient.Client, message *Message.Message) error {
			app.mutex.Lock()
			gameId := app.clientGameIds[client.GetId()]
			app.mutex.Unlock()
			if gameId == "" {
				return Utilities.NewError("You are not in a game", nil)
			}
			err := app.client.AsyncMessage(topics.END, app.client.GetName(), gameId)
			if err != nil {
				app.client.GetLogger().Log(Utilities.NewError("Error sending end message for game: "+gameId, err).Error())
			}
			return nil
		},
		"move": func(client *WebsocketClient.Client, message *Message.Message) error {
			app.mutex.Lock()
			defer app.mutex.Unlock()
			gameId := app.clientGameIds[client.GetId()]
			if gameId == "" {
				return Utilities.NewError("You are not in a game", nil)
			}
			moveSegments := strings.Split(message.GetPayload(), " ")
			if len(moveSegments) != 4 {
				return Utilities.NewError("Invalid move format", nil)
			}
			err := app.handleMove(gameId, client.GetId(), message.GetPayload())
			if err != nil {
				return Utilities.NewError("Error handling move", err)
			}
			return nil
		},
	}
}

func (app *AppWebsocketHTTP) handleMove(gameId, playerId, move string) error {
	segments := strings.Split(move, " ")
	if len(segments) != 4 {
		return Utilities.NewError("Invalid message format", nil)
	}
	responseMessage, err := app.client.SyncMessage(gameId, playerId, move)
	if err != nil {
		return Utilities.NewError("Error sending move message", err)
	}
	app.client.GetWebsocketServer().Groupcast(gameId, Message.NewAsync("propagate_move", responseMessage.GetOrigin(), responseMessage.GetPayload()))
	return nil

}

func (app *AppWebsocketHTTP) OnConnectHandler(client *WebsocketClient.Client) {
	err := client.Send(Message.NewAsync("connected", app.client.GetName(), client.GetId()).Serialize())
	if err != nil {
		client.Disconnect()
		app.client.GetLogger().Log(Utilities.NewError("Error sending connected message", err).Error())
	}
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(client *WebsocketClient.Client) {
	app.mutex.Lock()
	gameId := app.clientGameIds[client.GetId()]
	app.mutex.Unlock()
	if gameId == "" {
		return
	}
	err := app.client.AsyncMessage(topics.END, app.client.GetName(), gameId)
	if err != nil {
		app.client.GetLogger().Log(Utilities.NewError("Error sending end message for game: "+gameId, err).Error())
	}
}
