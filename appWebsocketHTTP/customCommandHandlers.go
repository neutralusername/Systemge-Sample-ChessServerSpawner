package appWebsocketHTTP

import (
	"Systemge/Application"
	"Systemge/Utilities"
)

func (app *AppWebsocketHTTP) GetCustomCommandHandlers() map[string]Application.CustomCommandHandler {
	return map[string]Application.CustomCommandHandler{
		"move": func(args []string) error {
			if len(args) != 6 {
				return Utilities.NewError("Invalid move command", nil)
			}
			gameId := args[0]
			playerId := args[1]
			rowFrom := args[2]
			colFrom := args[3]
			rowTo := args[4]
			colTo := args[5]
			err := app.handleMove(gameId, playerId, rowFrom+" "+colFrom+" "+rowTo+" "+colTo)
			if err != nil {
				return Utilities.NewError("Error handling move", err)
			}
			return nil
		},
	}
}
