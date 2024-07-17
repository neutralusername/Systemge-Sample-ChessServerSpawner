package appChess

import (
	"Systemge/Config"
	"Systemge/Error"
	"Systemge/Message"
	"Systemge/Node"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/topics"
	"strings"
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

func (app *App) GetSystemgeConfig() Config.Systemge {
	return Config.Systemge{
		HandleMessagesSequentially: false,
	}
}

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{}
}

func (app *App) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		app.gameId: func(node *Node.Node, message *Message.Message) (string, error) {
			segments := strings.Split(message.GetPayload(), " ")
			if len(segments) != 4 {
				return "", Error.New("Invalid message format", nil)
			}
			chessMove, err := app.handleMoveRequest(message.GetOrigin(), Utilities.StringToInt(segments[0]), Utilities.StringToInt(segments[1]), Utilities.StringToInt(segments[2]), Utilities.StringToInt(segments[3]))
			if err != nil {
				return "", err
			}
			return chessMove.Marshal(), nil
		},
	}
}

func (app *App) handleMoveRequest(playerId string, rowFrom, colFrom, rowTo, colTo int) (*ChessMove, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if app.isWhiteTurn() && playerId != app.whiteId {
		return nil, Error.New("Not your turn", nil)
	}
	if !app.isWhiteTurn() && playerId != app.blackId {
		return nil, Error.New("Not your turn", nil)
	}
	chessMove, err := app.move(rowFrom, colFrom, rowTo, colTo)
	if err != nil {
		return nil, Error.New("Invalid move", err)
	}
	return chessMove, nil
}
