package appSpawner

import (
	"Systemge/Application"
	"Systemge/Utilities"
)

func (app *App) GetCustomCommandHandlers() map[string]Application.CustomCommandHandler {
	return map[string]Application.CustomCommandHandler{
		"activeClients": app.activeClients,
		"endClient":     app.endClient,
	}
}

func (app *App) activeClients(args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for id := range app.spawnedClients {
		println(id)
	}
	return nil
}

func (app *App) endClient(args []string) error {
	if len(args) != 1 {
		return Utilities.NewError("No client id provided", nil)
	}
	id := args[0]
	err := app.EndClient(id)
	if err != nil {
		return Utilities.NewError("Error ending client "+id, err)
	}
	return nil
}
