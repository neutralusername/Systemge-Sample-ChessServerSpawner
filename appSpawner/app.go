package appSpawner

import (
	"Systemge/Application"
	"Systemge/Client"
	"sync"
)

type App struct {
	client *Client.Client

	spawnedClients map[string]*Client.Client
	mutex          sync.Mutex
}

func New(client *Client.Client, args []string) (Application.Application, error) {
	app := &App{
		client:         client,
		spawnedClients: make(map[string]*Client.Client),
	}
	return app, nil
}

func (app *App) OnStart() error {
	return nil
}

func (app *App) OnStop() error {
	return nil
}
