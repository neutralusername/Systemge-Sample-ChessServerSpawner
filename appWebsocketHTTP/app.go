package appWebsocketHTTP

import (
	"Systemge/Application"
	"Systemge/Client"
	"sync"
)

type AppWebsocketHTTP struct {
	client *Client.Client

	clientGameIds map[string]string
	mutex         sync.Mutex
}

func New(messageBrokerClient *Client.Client, args []string) (Application.CompositeApplicationWebsocketHTTP, error) {
	return &AppWebsocketHTTP{
		client: messageBrokerClient,

		clientGameIds: make(map[string]string),
	}, nil
}

func (app *AppWebsocketHTTP) OnStart() error {
	return nil
}

func (app *AppWebsocketHTTP) OnStop() error {
	return nil
}
