package appWebsocketHTTP

import (
	"Systemge/Node"
	"sync"
)

type AppWebsocketHTTP struct {
	clientGameIds map[string]string
	mutex         sync.Mutex
}

func New() Node.WebsocketHTTPApplication {
	return &AppWebsocketHTTP{
		clientGameIds: make(map[string]string),
	}
}

func (app *AppWebsocketHTTP) OnStart(client *Node.Node) error {
	return nil
}

func (app *AppWebsocketHTTP) OnStop(client *Node.Node) error {
	return nil
}
