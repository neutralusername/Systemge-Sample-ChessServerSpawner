package appChess

import (
	"Systemge/Error"
	"Systemge/Node"
	"SystemgeSampleChessServer/topics"
	"strings"
	"sync"
)

type App struct {
	gameId  string
	whiteId string
	blackId string
	board   [8][8]Piece
	moves   []ChessMove
	mutex   sync.Mutex
	mode960 bool
}

func New(id string) Node.Application {
	ids := strings.Split(id, "-")
	app := &App{
		gameId:  id,
		whiteId: ids[0],
		blackId: ids[1],
		mode960: false,
	}
	if app.mode960 {
		app.board = get960StartingPosition()
	} else {
		app.board = getStandardStartingPosition()
	}
	return app
}

func (app *App) OnStart(node *Node.Node) error {
	_, err := node.SyncMessage(topics.PROPAGATE_GAMESTART, node.GetName(), app.marshalBoard())
	if err != nil {
		node.GetLogger().Log(Error.New("Error sending sync message", err).Error())
		err := node.AsyncMessage(topics.END, node.GetName(), node.GetName())
		if err != nil {
			node.GetLogger().Log(Error.New("Error sending async message", err).Error())
		}
	}
	return nil
}

func (app *App) OnStop(node *Node.Node) error {
	err := node.AsyncMessage(topics.PROPAGATE_GAMEEND, node.GetName(), "...gameEndData...")
	if err != nil {
		node.GetLogger().Log(Error.New("Error sending async message", err).Error())
	}
	return nil
}
