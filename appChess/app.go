package appChess

import (
	"Systemge/Application"
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/topics"
	"strings"
	"sync"
)

type App struct {
	client *Client.Client

	board [8][8]Piece
	moves []ChessMove
	mutex sync.Mutex
}

func New(client *Client.Client, args []string) (Application.Application, error) {
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
	}
	return app, nil
}

func (app *App) marshalBoard() string {
	var builder strings.Builder
	for _, row := range app.board {
		for _, piece := range row {
			if piece == nil {
				builder.WriteString(".")
			} else {
				builder.WriteString(piece.getLetter())
			}
		}
	}
	return builder.String()
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

func (app *App) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{}
}

func (app *App) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{
		app.client.GetName(): func(message *Message.Message) (string, error) {
			segments := strings.Split(message.GetPayload(), " ")
			if len(segments) != 4 {
				return "", Utilities.NewError("Invalid message format", nil)
			}
			row1, col1, row2, col2 := Utilities.StringToInt(segments[0]), Utilities.StringToInt(segments[1]), Utilities.StringToInt(segments[2]), Utilities.StringToInt(segments[3])
			app.mutex.Lock()
			defer app.mutex.Unlock()
			_, err := app.move(row1, col1, row2, col2)
			if err != nil {
				return "", Utilities.NewError("Invalid move", err)
			}
			return app.marshalBoard(), nil
		},
	}
}

func (app *App) GetCustomCommandHandlers() map[string]Application.CustomCommandHandler {
	return map[string]Application.CustomCommandHandler{}
}
