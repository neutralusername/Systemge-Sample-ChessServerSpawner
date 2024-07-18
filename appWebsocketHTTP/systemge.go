package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Error"
	"Systemge/Helpers"
	"Systemge/Message"
	"Systemge/Node"
	"Systemge/Tcp"
	"SystemgeSampleChessServer/topics"
	"strings"
)

func (app *AppWebsocketHTTP) GetSystemgeComponentConfig() Config.Systemge {
	return Config.Systemge{
		HandleMessagesSequentially: false,

		BrokerSubscribeDelayMs:    1000,
		TopicResolutionLifetimeMs: 10000,
		SyncResponseTimeoutMs:     10000,
		TcpTimeoutMs:              5000,

		ResolverEndpoint: Tcp.NewEndpoint("127.0.0.1:60000", "example.com", Helpers.GetFileContent("MyCertificate.crt")),
	}
}

func (app *AppWebsocketHTTP) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		topics.PROPAGATE_GAMEEND: func(node *Node.Node, message *Message.Message) error {
			gameId := message.GetOrigin()
			ids := strings.Split(gameId, "-")
			node.WebsocketGroupcast(message.GetOrigin(), message)
			err := node.RemoveFromWebsocketGroup(gameId, ids[0])
			if err != nil {
				node.GetLogger().Error(Error.New("Error removing \""+ids[0]+"\" from group \""+gameId+"\"", err).Error())
			}
			err = node.RemoveFromWebsocketGroup(gameId, ids[1])
			if err != nil {
				node.GetLogger().Error(Error.New("Error removing \""+ids[1]+"\" from group \""+gameId+"\"", err).Error())
			}
			app.mutex.Lock()
			delete(app.nodeIds, ids[0])
			delete(app.nodeIds, ids[1])
			app.mutex.Unlock()
			return nil
		},
	}
}

func (app *AppWebsocketHTTP) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		topics.PROPAGATE_GAMESTART: func(node *Node.Node, message *Message.Message) (string, error) {
			gameId := message.GetOrigin()
			ids := strings.Split(gameId, "-")
			err := node.AddToWebsocketGroup(gameId, ids[0])
			if err != nil {
				return "", Error.New("Error adding \""+ids[0]+"\" to group \""+gameId+"\"", err)
			}
			err = node.AddToWebsocketGroup(gameId, ids[1])
			if err != nil {
				err := node.RemoveFromWebsocketGroup(gameId, ids[0])
				if err != nil {
					node.GetLogger().Warning(Error.New("Error removing \""+ids[0]+"\" from group \""+gameId+"\"", err).Error())
				}
				return "", Error.New("Error adding \""+ids[1]+"\" to group \""+gameId+"\"", err)
			}
			node.WebsocketGroupcast(gameId, message)
			return "", nil
		},
	}
}
