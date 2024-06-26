package appWebsocketHTTP

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Node"
	"SystemgeSampleChessServer/topics"
	"strings"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{
		"startGame": func(client *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			app.mutex.Lock()
			defer app.mutex.Unlock()
			whiteId := websocketClient.GetId()
			blackId := message.GetPayload()
			if !client.ClientExists(blackId) {
				return Error.New("Opponent does not exist", nil)
			}
			if blackId == whiteId {
				return Error.New("You cannot play against yourself", nil)
			}
			if app.clientGameIds[whiteId] != "" {
				return Error.New("You are already in a game", nil)
			}
			if app.clientGameIds[blackId] != "" {
				return Error.New("Opponent is already in a game", nil)
			}
			gameId := whiteId + "-" + blackId
			_, err := client.SyncMessage(topics.NEW, client.GetName(), gameId)
			if err != nil {
				return Error.New("Error spawning new game client", err)
			}
			app.clientGameIds[whiteId] = gameId
			app.clientGameIds[blackId] = gameId
			return nil
		},
		"endGame": func(client *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			app.mutex.Lock()
			gameId := app.clientGameIds[websocketClient.GetId()]
			app.mutex.Unlock()
			if gameId == "" {
				return Error.New("You are not in a game", nil)
			}
			err := client.AsyncMessage(topics.END, client.GetName(), gameId)
			if err != nil {
				client.GetLogger().Log(Error.New("Error sending end message for game: "+gameId, err).Error())
			}
			client.RemoveTopicResolution(gameId)
			return nil
		},
		"move": func(client *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			app.mutex.Lock()
			defer app.mutex.Unlock()
			gameId := app.clientGameIds[websocketClient.GetId()]
			if gameId == "" {
				return Error.New("You are not in a game", nil)
			}
			moveSegments := strings.Split(message.GetPayload(), " ")
			if len(moveSegments) != 4 {
				return Error.New("Invalid move format", nil)
			}
			err := app.handleMove(client, gameId, websocketClient.GetId(), message.GetPayload())
			if err != nil {
				return Error.New("Error handling move", err)
			}
			return nil
		},
	}
}

func (app *AppWebsocketHTTP) handleMove(client *Node.Node, gameId, playerId, move string) error {
	segments := strings.Split(move, " ")
	if len(segments) != 4 {
		return Error.New("Invalid message format", nil)
	}
	responseMessage, err := client.SyncMessage(gameId, playerId, move)
	if err != nil {
		return Error.New("Error sending move message", err)
	}
	client.WebsocketGroupcast(gameId, Message.NewAsync("propagate_move", responseMessage.GetOrigin(), responseMessage.GetPayload()))
	return nil

}

func (app *AppWebsocketHTTP) OnConnectHandler(client *Node.Node, websocketClient *Node.WebsocketClient) {
	err := websocketClient.Send(Message.NewAsync("connected", client.GetName(), websocketClient.GetId()).Serialize())
	if err != nil {
		websocketClient.Disconnect()
		client.GetLogger().Log(Error.New("Error sending connected message", err).Error())
	}
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(client *Node.Node, websocketClient *Node.WebsocketClient) {
	app.mutex.Lock()
	gameId := app.clientGameIds[websocketClient.GetId()]
	app.mutex.Unlock()
	if gameId == "" {
		return
	}
	err := client.AsyncMessage(topics.END, client.GetName(), gameId)
	if err != nil {
		client.GetLogger().Log(Error.New("Error sending end message for game: "+gameId, err).Error())
	}
	client.RemoveTopicResolution(gameId)
}
