package appChess

import (
	"Systemge/Error"
	"Systemge/Node"
	"SystemgeSampleChessServer/topics"
)

func (app *App) OnStart(node *Node.Node) error {
	_, err := node.SyncMessage(topics.PROPAGATE_GAMESTART, node.GetName(), app.marshalBoard())
	if err != nil {
		node.GetLogger().Warning(Error.New("Error sending sync message", err).Error())
		err := node.AsyncMessage(topics.END_NODE_ASYNC, node.GetName(), node.GetName())
		if err != nil {
			node.GetLogger().Error(Error.New("Error sending async message", err).Error())
		}
	}
	return nil
}

func (app *App) OnStop(node *Node.Node) error {
	err := node.AsyncMessage(topics.PROPAGATE_GAMEEND, node.GetName(), "")
	if err != nil {
		node.GetLogger().Error(Error.New("Error sending async message", err).Error())
	}
	return nil
}
