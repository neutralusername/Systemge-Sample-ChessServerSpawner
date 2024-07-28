package appWebsocketHTTP

import (
	"SystemgeSampleChessServer/topics"
	"strings"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
	"github.com/neutralusername/Systemge/Spawner"
)

func (app *AppWebsocketHTTP) GetSystemgeComponentConfig() *Config.Systemge {
	return &Config.Systemge{
		HandleMessagesSequentially: false,

		BrokerSubscribeDelayMs:    1000,
		TopicResolutionLifetimeMs: 10000,
		SyncResponseTimeoutMs:     10000,
		TcpTimeoutMs:              5000,

		ResolverEndpoint: &Config.TcpEndpoint{
			Address: "127.0.0.1:60000",
			Domain:  "example.com",
			TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
		},
	}
}

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		topics.PROPAGATE_GAMEEND: func(node *Node.Node, message *Message.Message) error {
			gameId := message.GetOrigin()
			ids := strings.Split(gameId, "-")
			node.WebsocketGroupcast(message.GetOrigin(), message)
			err := node.RemoveFromWebsocketGroup(gameId, ids...)
			if err != nil {
				if errorLogger := node.GetErrorLogger(); errorLogger != nil {
					errorLogger.Log(Error.New("Error removing from websocket group", err).Error())
				}
			}
			app.mutex.Lock()
			delete(app.nodeIds, ids[0])
			delete(app.nodeIds, ids[1])
			app.mutex.Unlock()
			err = node.AsyncMessage(Spawner.DESPAWN_NODE_ASYNC, node.GetName(), gameId)
			if err != nil {
				if errorLogger := node.GetErrorLogger(); errorLogger != nil {
					errorLogger.Log(Error.New("Error despawning node", err).Error())
				}
			}
			return nil
		},
	}
}

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		topics.PROPAGATE_GAMESTART: func(node *Node.Node, message *Message.Message) (string, error) {
			gameId := message.GetOrigin()
			ids := strings.Split(gameId, "-")
			err := node.AddToWebsocketGroup(gameId, ids...)
			if err != nil {
				return "", Error.New("Error adding \""+ids[0]+"\" to group \""+gameId+"\"", err)
			}
			app.mutex.Lock()
			app.nodeIds[ids[0]] = gameId
			app.nodeIds[ids[1]] = gameId
			app.mutex.Unlock()
			node.WebsocketGroupcast(gameId, message)
			return "", nil
		},
	}
}
