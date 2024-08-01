package appChess

import (
	"SystemgeSampleChessServer/dto"
	"sync"

	"github.com/neutralusername/Systemge/Node"
)

type App struct {
	whiteId string
	blackId string
	board   [8][8]Piece
	moves   []*dto.Move
	mutex   sync.Mutex
}

func New() Node.Application {
	app := &App{}
	if false {
		app.board = get960StartingPosition()
	} else {
		app.board = getStandardStartingPosition()
	}
	return app
}
