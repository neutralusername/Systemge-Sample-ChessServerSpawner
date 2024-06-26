package appWebsocketHTTP

import (
	"Systemge/Node"
	"sync"
)

type AppWebsocketHTTP struct {
	nodeIds map[string]string
	mutex   sync.Mutex
}

func New() Node.WebsocketHTTPApplication {
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
