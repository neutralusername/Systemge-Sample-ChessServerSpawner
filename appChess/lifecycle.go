package appChess

import (
	"Systemge/Error"
	"Systemge/Node"
	"SystemgeSampleChessServer/topics"
)

func (app *App) OnStart(node *Node.Node) error {
	_, err := node.SyncMessage(topics.PROPAGATE_GAMESTART, node.GetName(), app.marshalBoard())
	if err != nil {
		if warningLogger := node.GetWarningLogger(); warningLogger != nil {
			warningLogger.Log(Error.New("Error sending sync message", err).Error())
		}
	}
	return nil
}

func (app *App) OnStop(node *Node.Node) error {
	err := node.AsyncMessage(topics.PROPAGATE_GAMEEND, node.GetName(), "")
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log(Error.New("Error sending async message", err).Error())
		}
	}
	return nil
}
