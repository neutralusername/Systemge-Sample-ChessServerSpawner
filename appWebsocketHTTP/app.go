package appWebsocketHTTP

import (
	"Systemge/Config"
	"Systemge/Node"
	"Systemge/Resolution"
	"Systemge/Utilities"
	"sync"
)

type AppWebsocketHTTP struct {
	nodeIds map[string]string
	mutex   sync.Mutex
}

func New() *AppWebsocketHTTP {
	return &AppWebsocketHTTP{
		nodeIds: make(map[string]string),
	}
}

func (app *AppWebsocketHTTP) OnStart(node *Node.Node) error {
	return nil
}

func (app *AppWebsocketHTTP) OnStop(node *Node.Node) error {
	return nil
}

func (app *AppWebsocketHTTP) GetApplicationConfig() Config.Application {
	return Config.Application{
		ResolverResolution:         Resolution.New("resolver", "127.0.0.1:60000", "127.0.0.1", Utilities.GetFileContent("MyCertificate.crt")),
		HandleMessagesSequentially: false,
	}
}
