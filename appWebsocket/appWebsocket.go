package appWebsocket

import (
	"Systemge/Application"
	"Systemge/Client"
	"sync"
)

type AppWebsocket struct {
	client *Client.Client

	clientGameIds map[string]string
	mutex         sync.Mutex
}

func New(messageBrokerClient *Client.Client, args []string) (Application.WebsocketApplication, error) {
	return &AppWebsocket{
		client: messageBrokerClient,

		clientGameIds: make(map[string]string),
	}, nil
}

func (app *AppWebsocket) OnStart() error {
	return nil
}

func (app *AppWebsocket) OnStop() error {
	return nil
}
