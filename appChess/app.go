package appChess

import (
	"Systemge/Node"
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
