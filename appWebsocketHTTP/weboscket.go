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
		"startGame": func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			app.mutex.Lock()
			defer app.mutex.Unlock()
			whiteId := websocketClient.GetId()
			blackId := message.GetPayload()
			if !node.WebsocketClientExists(blackId) {
				return Error.New("Opponent does not exist", nil)
			}
			if blackId == whiteId {
				return Error.New("You cannot play against yourself", nil)
			}
			if app.nodeIds[whiteId] != "" {
				return Error.New("You are already in a game", nil)
			}
			if app.nodeIds[blackId] != "" {
				return Error.New("Opponent is already in a game", nil)
			}
			gameId := whiteId + "-" + blackId
			_, err := node.SyncMessage(topics.NEW, node.GetName(), gameId)
			if err != nil {
				return Error.New("Error spawning new game client", err)
			}
			app.nodeIds[whiteId] = gameId
			app.nodeIds[blackId] = gameId
			return nil
		},
		"endGame": func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			app.mutex.Lock()
			gameId := app.nodeIds[websocketClient.GetId()]
			app.mutex.Unlock()
			if gameId == "" {
				return Error.New("You are not in a game", nil)
			}
			err := node.AsyncMessage(topics.END, node.GetName(), gameId)
			if err != nil {
				node.GetLogger().Log(Error.New("Error sending end message for game: "+gameId, err).Error())
			}
			node.RemoveTopicResolution(gameId)
			return nil
		},
		"move": func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			app.mutex.Lock()
			defer app.mutex.Unlock()
			gameId := app.nodeIds[websocketClient.GetId()]
			if gameId == "" {
				return Error.New("You are not in a game", nil)
			}
			moveSegments := strings.Split(message.GetPayload(), " ")
			if len(moveSegments) != 4 {
				return Error.New("Invalid move format", nil)
			}
			err := app.handleMove(node, gameId, websocketClient.GetId(), message.GetPayload())
			if err != nil {
				return Error.New("Error handling move", err)
			}
			return nil
		},
	}
}

func (app *AppWebsocketHTTP) handleMove(node *Node.Node, gameId, playerId, move string) error {
	segments := strings.Split(move, " ")
	if len(segments) != 4 {
		return Error.New("Invalid message format", nil)
	}
	responseMessage, err := node.SyncMessage(gameId, playerId, move)
	if err != nil {
		return Error.New("Error sending move message", err)
	}
	node.WebsocketGroupcast(gameId, Message.NewAsync("propagate_move", responseMessage.GetOrigin(), responseMessage.GetPayload()))
	return nil

}

func (app *AppWebsocketHTTP) OnConnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	err := websocketClient.Send(Message.NewAsync("connected", node.GetName(), websocketClient.GetId()).Serialize())
	if err != nil {
		websocketClient.Disconnect()
		node.GetLogger().Log(Error.New("Error sending connected message", err).Error())
	}
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	app.mutex.Lock()
	gameId := app.nodeIds[websocketClient.GetId()]
	app.mutex.Unlock()
	if gameId == "" {
		return
	}
	err := node.AsyncMessage(topics.END, node.GetName(), gameId)
	if err != nil {
		node.GetLogger().Log(Error.New("Error sending end message for game: "+gameId, err).Error())
	}
	node.RemoveTopicResolution(gameId)
}
