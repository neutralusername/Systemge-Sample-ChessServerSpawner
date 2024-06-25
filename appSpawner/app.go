package appSpawner

import (
	"Systemge/Client"
	"sync"
)

type App struct {
	spawnedClients map[string]*Client.Client
	mutex          sync.Mutex
}

func New() Client.Application {
	app := &App{
		spawnedClients: make(map[string]*Client.Client),
	}
	return app
}

func (app *App) OnStart(client *Client.Client) error {
	return nil
}

func (app *App) OnStop(client *Client.Client) error {
	return nil
}
