package appChess

import (
	"Systemge/Client"
	"Systemge/Error"
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

func New(id string) Client.Application {
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

func (app *App) OnStart(client *Client.Client) error {
	_, err := client.SyncMessage(topics.PROPAGATE_GAMESTART, client.GetName(), app.marshalBoard())
	if err != nil {
		client.GetLogger().Log(Error.New("Error sending sync message", err).Error())
		err := client.AsyncMessage(topics.END, client.GetName(), client.GetName())
		if err != nil {
			client.GetLogger().Log(Error.New("Error sending async message", err).Error())
		}
	}
	return nil
}

func (app *App) OnStop(client *Client.Client) error {
	err := client.AsyncMessage(topics.PROPAGATE_GAMEEND, client.GetName(), "...gameEndData...")
	if err != nil {
		client.GetLogger().Log(Error.New("Error sending async message", err).Error())
	}
	return nil
}
