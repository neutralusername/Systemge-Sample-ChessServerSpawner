package appSpawner

import (
	"Systemge/Node"
	"sync"
)

type App struct {
	spawnedClients map[string]*Node.Node
	mutex          sync.Mutex
}

func New() Node.Application {
	app := &App{
		spawnedClients: make(map[string]*Node.Node),
	}
	return app
}

func (app *App) OnStart(client *Node.Node) error {
	return nil
}

func (app *App) OnStop(client *Node.Node) error {
	return nil
}
