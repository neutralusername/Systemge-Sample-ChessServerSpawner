package appSpawner

import (
	"Systemge/Error"
	"Systemge/Node"
)

func (app *App) GetCustomCommandHandlers() map[string]Node.CustomCommandHandler {
	return map[string]Node.CustomCommandHandler{
		"spawnedNodes":   app.activeNodes,
		"endSpawnedNode": app.endNode,
	}
}

func (app *App) activeNodes(node *Node.Node, args []string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	for id := range app.spawnedNodes {
		println(id)
	}
	return nil
}

func (app *App) endNode(node *Node.Node, args []string) error {
	if len(args) != 1 {
		return Error.New("No nodeId provided", nil)
	}
	id := args[0]
	err := app.EndNode(node, id)
	if err != nil {
		return Error.New("Error ending node "+id, err)
	}
	return nil
}
