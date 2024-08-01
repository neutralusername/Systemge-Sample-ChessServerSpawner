package appWebsocketHTTP

import (
	"sync"
	"sync/atomic"

	"github.com/neutralusername/Systemge/Node"
)

type AppWebsocketHTTP struct {
	gameIds map[string]string // playerId -> gameId
	ports   atomic.Uint32
	mutex   sync.Mutex
}

func New() *AppWebsocketHTTP {
	app := &AppWebsocketHTTP{
		gameIds: make(map[string]string),
	}
	app.ports.Store(60003)
	return app
}

func (app *AppWebsocketHTTP) GetCommandHandlers() map[string]Node.CommandHandler {
	return map[string]Node.CommandHandler{}
}
