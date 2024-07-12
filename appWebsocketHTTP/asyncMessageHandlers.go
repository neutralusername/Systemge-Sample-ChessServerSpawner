package appWebsocketHTTP

import (
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Node"
	"SystemgeSampleChessServer/topics"
	"strings"
)

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
