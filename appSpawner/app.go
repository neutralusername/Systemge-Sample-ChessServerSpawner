package appSpawner

import (
	"Systemge/Node"
	"Systemge/Utilities"
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

func (app *App) GetApplicationConfig() Node.ApplicationConfig {
	return Node.ApplicationConfig{
		ResolverAddress:            "127.0.0.1:60000",
		ResolverNameIndication:     "127.0.0.1",
		ResolverTLSCert:            Utilities.GetFileContent("MyCertificate.crt"),
		HandleMessagesSequentially: false,
	}
}
