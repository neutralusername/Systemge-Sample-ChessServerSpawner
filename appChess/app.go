package appChess

import (
	"Systemge/Application"
	"Systemge/Client"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/topics"
	"strings"
	"sync"
)

type App struct {
	client *Client.Client

	whiteId string
	blackId string
	board   [8][8]Piece
	moves   []ChessMove
	mutex   sync.Mutex
}

func New(client *Client.Client, args []string) (Application.Application, error) {
	ids := strings.Split(client.GetName(), "-")
	if len(ids) != 2 {
		return nil, Utilities.NewError("Invalid client name", nil)
	}
	app := &App{
		client: client,
		board: [8][8]Piece{ // first index is row, second index is column. white rooks are at 0,0 and 0,7. black rooks are at 7,0 and 7,7
			{&Rook{true, false}, &Knight{true}, &Bishop{true}, &Queen{true}, &King{true, false}, &Bishop{true}, &Knight{true}, &Rook{true, false}},
			{&Pawn{true}, &Pawn{true}, &Pawn{true}, &Pawn{true}, &Pawn{true}, &Pawn{true}, &Pawn{true}, &Pawn{true}},
			{nil, nil, nil, nil, nil, nil, nil, nil},
			{nil, nil, nil, nil, nil, nil, nil, nil},
			{nil, nil, nil, nil, nil, nil, nil, nil},
			{nil, nil, nil, nil, nil, nil, nil, nil},
			{&Pawn{false}, &Pawn{false}, &Pawn{false}, &Pawn{false}, &Pawn{false}, &Pawn{false}, &Pawn{false}, &Pawn{false}},
			{&Rook{false, false}, &Knight{false}, &Bishop{false}, &Queen{false}, &King{false, false}, &Bishop{false}, &Knight{false}, &Rook{false, false}},
		},
		whiteId: ids[0],
		blackId: ids[1],
	}
	return app, nil
}

func (app *App) OnStart() error {
	_, err := app.client.SyncMessage(topics.PROPAGATE_GAMESTART, app.client.GetName(), app.marshalBoard())
	if err != nil {
		app.client.GetLogger().Log(Utilities.NewError("Error sending sync message", err).Error())
		err := app.client.AsyncMessage(topics.END, app.client.GetName(), app.client.GetName())
		if err != nil {
			app.client.GetLogger().Log(Utilities.NewError("Error sending async message", err).Error())
		}
	}
	return nil
}

func (app *App) OnStop() error {
	err := app.client.AsyncMessage(topics.PROPAGATE_GAMEEND, app.client.GetName(), "...gameEndData...")
	if err != nil {
		app.client.GetLogger().Log(Utilities.NewError("Error sending async message", err).Error())
	}
	return nil
}
