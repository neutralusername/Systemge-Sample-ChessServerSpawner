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
	mode960 bool
}

func New(client *Client.Client, args []string) (Application.Application, error) {
	ids := strings.Split(client.GetName(), "-")
	if len(ids) != 2 {
		return nil, Utilities.NewError("Invalid client name", nil)
	}
	app := &App{
		client:  client,
		whiteId: ids[0],
		blackId: ids[1],
		mode960: false,
	}
	if app.mode960 {
		app.board = get960StartingPosition()
	} else {
		app.board = getStandardStartingPosition()
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
