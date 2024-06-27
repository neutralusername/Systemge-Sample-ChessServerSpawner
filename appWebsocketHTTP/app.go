package appWebsocketHTTP

import (
	"Systemge/Node"
	"Systemge/Utilities"
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

func (app *AppWebsocketHTTP) GetApplicationConfig() Node.ApplicationConfig {
	return Node.ApplicationConfig{
		ResolverAddress:            "127.0.0.1:60000",
		ResolverNameIndication:     "127.0.0.1",
		ResolverTLSCert:            Utilities.GetFileContent("MyCertificate.crt"),
		HandleMessagesSequentially: false,
	}
}
