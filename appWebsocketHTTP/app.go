package appWebsocketHTTP

import (
	"Systemge/Error"
	"Systemge/Node"
	"sync"
)

type AppWebsocketHTTP struct {
	nodeIds map[string]string
	mutex   sync.Mutex
}

func New() *AppWebsocketHTTP {
	return &AppWebsocketHTTP{
		nodeIds: make(map[string]string),
	}
}

func (app *AppWebsocketHTTP) GetCommandHandlers() map[string]Node.CustomCommandHandler {
	return map[string]Node.CustomCommandHandler{
		"move": func(node *Node.Node, args []string) error {
			if len(args) != 6 {
				return Error.New("Invalid move command", nil)
			}
			gameId := args[0]
			playerId := args[1]
			rowFrom := args[2]
			colFrom := args[3]
			rowTo := args[4]
			colTo := args[5]
			err := app.handleMove(node, gameId, playerId, rowFrom+" "+colFrom+" "+rowTo+" "+colTo)
			if err != nil {
				return Error.New("Error handling move", err)
			}
			return nil
		},
	}
}
