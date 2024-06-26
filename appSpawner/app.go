package appSpawner

import (
	"Systemge/Node"
	"sync"
)

type App struct {
	spawnedNodes map[string]*Node.Node
	mutex        sync.Mutex
}

func New() Node.Application {
	app := &App{
		spawnedNodes: make(map[string]*Node.Node),
	}
	return app
}

func (app *App) OnStart(node *Node.Node) error {
	return nil
}

func (app *App) OnStop(node *Node.Node) error {
	return nil
}
