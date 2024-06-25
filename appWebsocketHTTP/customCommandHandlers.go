package appWebsocketHTTP

import (
	"Systemge/Client"
	"Systemge/Error"
)

func (app *AppWebsocketHTTP) GetCustomCommandHandlers() map[string]Client.CustomCommandHandler {
	return map[string]Client.CustomCommandHandler{
		"move": func(client *Client.Client, args []string) error {
			if len(args) != 6 {
				return Error.New("Invalid move command", nil)
			}
			gameId := args[0]
			playerId := args[1]
			rowFrom := args[2]
			colFrom := args[3]
			rowTo := args[4]
			colTo := args[5]
			err := app.handleMove(client, gameId, playerId, rowFrom+" "+colFrom+" "+rowTo+" "+colTo)
			if err != nil {
				return Error.New("Error handling move", err)
			}
			return nil
		},
	}
}
