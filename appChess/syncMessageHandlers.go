package appChess

import (
	"Systemge/Application"
	"Systemge/Message"
	"Systemge/Utilities"
	"strings"
)

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
			if app.isWhiteTurn() && message.GetOrigin() != app.whiteId {
				return "", Utilities.NewError("Not your turn", nil)
			}
			if !app.isWhiteTurn() && message.GetOrigin() != app.blackId {
				return "", Utilities.NewError("Not your turn", nil)
			}
			_, err := app.move(row1, col1, row2, col2)
			if err != nil {
				return "", Utilities.NewError("Invalid move", err)
			}
			return app.marshalBoard(), nil
		},
	}
}
