package appWebsocketHTTP

import (
	"SystemgeSampleChessServer/dto"
	"SystemgeSampleChessServer/topics"
	"strings"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		topics.PROPAGATE_GAMEEND: func(node *Node.Node, message *Message.Message) error {
			gameId := message.GetOrigin()
			ids := strings.Split(gameId, "-")
			node.WebsocketGroupcast(gameId, message)
			node.RemoveFromWebsocketGroup(gameId, ids...)
			tcpEndpointConfig := Config.UnmarshalTcpEndpoint(message.GetPayload())
			if err := node.DisconnectFromNode(tcpEndpointConfig.Address); err != nil {
				panic(Error.New("Error disconnecting from \""+tcpEndpointConfig.Address+"\"", err))
			}
			app.mutex.Lock()
			defer app.mutex.Unlock()
			delete(app.gameIds, ids[0])
			delete(app.gameIds, ids[1])
			return nil
		},
		topics.PROPAGATE_GAMESTART: func(node *Node.Node, message *Message.Message) error {
			gameId := message.GetOrigin()
			gameStart := dto.UnmarshalGameStart(message.GetPayload())
			ids := strings.Split(gameId, "-")
			err := node.AddToWebsocketGroup(gameId, ids...)
			if err != nil {
				panic(Error.New("Error adding \""+ids[0]+"\" to group \""+gameId+"\"", err))
			}
			if err := node.ConnectToNode(gameStart.TcpEndpointConfig, false); err != nil {
				panic(Error.New("Error connecting to \""+gameStart.TcpEndpointConfig.Address+"\"", err))
			}
			app.mutex.Lock()
			app.gameIds[ids[0]] = gameId
			app.gameIds[ids[1]] = gameId
			app.mutex.Unlock()
			node.WebsocketGroupcast(gameId, Message.NewAsync("propagate_gameStart", gameStart.Board))
			return nil
		},
	}
}

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{}
}
