package appSpawner

import (
	"Systemge/Client"
	"Systemge/Utilities"
)

func (app *App) GetCustomCommandHandlers() map[string]Client.CustomCommandHandler {
	return map[string]Client.CustomCommandHandler{
		"spawnedClients":   app.activeClients,
		"endSpawnedClient": app.endClient,
	}
}

func (app *App) activeClients(client *Client.Client, args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for id := range app.spawnedClients {
		println(id)
	}
	return nil
}

func (app *App) endClient(client *Client.Client, args []string) error {
	if len(args) != 1 {
		return Utilities.NewError("No client id provided", nil)
	}
	id := args[0]
	err := app.EndClient(client, id)
	if err != nil {
		return Utilities.NewError("Error ending client "+id, err)
	}
	return nil
}
