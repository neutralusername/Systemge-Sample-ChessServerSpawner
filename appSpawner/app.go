package appSpawner

import (
	"Systemge/Application"
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/topics"
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

func (app *App) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{
		topics.END: app.End,
	}
}

func (app *App) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{
		topics.NEW: app.New,
	}
}

func (app *App) GetCustomCommandHandlers() map[string]Application.CustomCommandHandler {
	return map[string]Application.CustomCommandHandler{
		"activeGames": app.activeGames,
		"endGame":     app.endGame,
	}
}

func (app *App) activeGames(args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for id := range app.spawnedClients {
		println(id)
	}
	return nil
}

func (app *App) endGame(args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if len(args) != 1 {
		return Utilities.NewError("No game id provided", nil)
	}
	id := args[0]
	err := app.EndClient(id)
	if err != nil {
		return Utilities.NewError("Error ending game "+id, err)
	}
	return nil
}

func (app *App) End(message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	id := message.GetPayload()
	err := app.EndClient(id)
	if err != nil {
		return Utilities.NewError("Error ending client "+id, err)
	}
	return nil
}

func (app *App) New(message *Message.Message) (string, error) {
	id := message.GetPayload()
	err := app.StartClient(id)
	if err != nil {
		return "", Utilities.NewError("Error starting client "+id, err)
	}
	return id, nil
}
