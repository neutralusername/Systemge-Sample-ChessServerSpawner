package appSpawner

import (
	"Systemge/Error"
	"Systemge/Node"
)

func (app *App) GetCustomCommandHandlers() map[string]Node.CustomCommandHandler {
	return map[string]Node.CustomCommandHandler{
		"spawnedClients":   app.activeClients,
		"endSpawnedClient": app.endClient,
	}
}

func (app *App) activeClients(client *Node.Node, args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for id := range app.spawnedClients {
		println(id)
	}
	return nil
}

func (app *App) endClient(client *Node.Node, args []string) error {
	if len(args) != 1 {
		return Error.New("No client id provided", nil)
	}
	id := args[0]
	err := app.EndClient(client, id)
	if err != nil {
		return Error.New("Error ending client "+id, err)
	}
	return nil
}
