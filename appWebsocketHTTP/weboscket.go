package appWebsocketHTTP

import (
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/topics"
	"strings"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Client.WebsocketMessageHandler {
	return map[string]Client.WebsocketMessageHandler{
		"startGame": func(client *Client.Client, websocketClient *Client.WebsocketClient, message *Message.Message) error {
			app.mutex.Lock()
			defer app.mutex.Unlock()
			whiteId := websocketClient.GetId()
			blackId := message.GetPayload()
			if !client.ClientExists(blackId) {
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
			_, err := client.SyncMessage(topics.NEW, client.GetName(), gameId)
			if err != nil {
				return Utilities.NewError("Error spawning new game client", err)
			}
			app.clientGameIds[whiteId] = gameId
			app.clientGameIds[blackId] = gameId
			return nil
		},
		"endGame": func(client *Client.Client, websocketClient *Client.WebsocketClient, message *Message.Message) error {
			app.mutex.Lock()
			gameId := app.clientGameIds[websocketClient.GetId()]
			app.mutex.Unlock()
			if gameId == "" {
				return Utilities.NewError("You are not in a game", nil)
			}
			err := client.AsyncMessage(topics.END, client.GetName(), gameId)
			if err != nil {
				client.GetLogger().Log(Utilities.NewError("Error sending end message for game: "+gameId, err).Error())
			}
			client.RemoveTopicResolution(gameId)
			return nil
		},
		"move": func(client *Client.Client, websocketClient *Client.WebsocketClient, message *Message.Message) error {
			app.mutex.Lock()
			defer app.mutex.Unlock()
			gameId := app.clientGameIds[websocketClient.GetId()]
			if gameId == "" {
				return Utilities.NewError("You are not in a game", nil)
			}
			moveSegments := strings.Split(message.GetPayload(), " ")
			if len(moveSegments) != 4 {
				return Utilities.NewError("Invalid move format", nil)
			}
			err := app.handleMove(client, gameId, websocketClient.GetId(), message.GetPayload())
			if err != nil {
				return Utilities.NewError("Error handling move", err)
			}
			return nil
		},
	}
}

func (app *AppWebsocketHTTP) handleMove(client *Client.Client, gameId, playerId, move string) error {
	segments := strings.Split(move, " ")
	if len(segments) != 4 {
		return Utilities.NewError("Invalid message format", nil)
	}
	responseMessage, err := client.SyncMessage(gameId, playerId, move)
	if err != nil {
		return Utilities.NewError("Error sending move message", err)
	}
	client.Groupcast(gameId, Message.NewAsync("propagate_move", responseMessage.GetOrigin(), responseMessage.GetPayload()))
	return nil

}

func (app *AppWebsocketHTTP) OnConnectHandler(client *Client.Client, websocketClient *Client.WebsocketClient) {
	err := websocketClient.Send(Message.NewAsync("connected", client.GetName(), websocketClient.GetId()).Serialize())
	if err != nil {
		websocketClient.Disconnect()
		client.GetLogger().Log(Utilities.NewError("Error sending connected message", err).Error())
	}
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(client *Client.Client, websocketClient *Client.WebsocketClient) {
	app.mutex.Lock()
	gameId := app.clientGameIds[websocketClient.GetId()]
	app.mutex.Unlock()
	if gameId == "" {
		return
	}
	err := client.AsyncMessage(topics.END, client.GetName(), gameId)
	if err != nil {
		client.GetLogger().Log(Utilities.NewError("Error sending end message for game: "+gameId, err).Error())
	}
	client.RemoveTopicResolution(gameId)
}
