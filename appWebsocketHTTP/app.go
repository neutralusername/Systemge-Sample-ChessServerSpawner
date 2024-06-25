package appWebsocketHTTP

import (
	"Systemge/Client"
	"sync"
)

type AppWebsocketHTTP struct {
	clientGameIds map[string]string
	mutex         sync.Mutex
}

func New() Client.CompositeApplicationWebsocketHTTP {
	return &AppWebsocketHTTP{
		clientGameIds: make(map[string]string),
	}
}

func (app *AppWebsocketHTTP) OnStart(client *Client.Client) error {
	return nil
}

func (app *AppWebsocketHTTP) OnStop(client *Client.Client) error {
	return nil
}
