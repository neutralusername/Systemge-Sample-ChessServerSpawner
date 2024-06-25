package appChess

import (
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"strings"
)

func (app *App) GetSyncMessageHandlers() map[string]Client.SyncMessageHandler {
	return map[string]Client.SyncMessageHandler{
		app.gameId: func(client *Client.Client, message *Message.Message) (string, error) {
			segments := strings.Split(message.GetPayload(), " ")
			if len(segments) != 4 {
				return "", Utilities.NewError("Invalid message format", nil)
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
		return nil, Utilities.NewError("Not your turn", nil)
	}
	if !app.isWhiteTurn() && playerId != app.blackId {
		return nil, Utilities.NewError("Not your turn", nil)
	}
	chessMove, err := app.move(rowFrom, colFrom, rowTo, colTo)
	if err != nil {
		return nil, Utilities.NewError("Invalid move", err)
	}
	return chessMove, nil
}
