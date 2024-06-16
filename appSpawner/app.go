package appSpawner

import (
	"Systemge/Application"
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/appChess"
	"SystemgeSampleChessServer/topics"
	"sync"
)

type App struct {
	client *Client.Client

	spawnedClients map[string]*Client.Client
	mutex          sync.Mutex
}

func New(client *Client.Client, args []string) (Application.Application, error) {
	app := &App{
		client:         client,
		spawnedClients: make(map[string]*Client.Client),
	}
	return app, nil
}

func (app *App) OnStart() error {
	return nil
}

func (app *App) OnStop() error {
	return nil
}

func (app *App) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{
		topics.END: app.End,
	}
}

func (app *App) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{
		topics.NEW: app.New,
	}
}

func (app *App) GetCustomCommandHandlers() map[string]Application.CustomCommandHandler {
	return map[string]Application.CustomCommandHandler{}
}

func (app *App) End(message *Message.Message) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	id := message.GetPayload()
	client := app.spawnedClients[id]
	if client == nil {
		return Utilities.NewError("Client "+id+" does not exist", nil)
	}
	err := client.Stop()
	if err != nil {
		return Utilities.NewError("Error stopping client "+id, err)
	}
	delete(app.spawnedClients, id)
	brokerNetConn, err := Utilities.TlsDial("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"))
	if err != nil {
		return Utilities.NewError("Error dialing broker", err)
	}
	_, err = Utilities.TcpExchange(brokerNetConn, Message.NewAsync("removeAsyncTopic", app.client.GetName(), message.GetPayload()), 5000)
	if err != nil {
		return Utilities.NewError("Error exchanging messages with broker", err)
	}
	resolverNetConn, err := Utilities.TlsDial("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"))
	if err != nil {
		return Utilities.NewError("Error dialing topic resolution server", err)
	}
	_, err = Utilities.TcpExchange(resolverNetConn, Message.NewAsync("unregisterTopics", app.client.GetName(), "brokerChess"+" "+message.GetPayload()), 5000)
	if err != nil {
		return Utilities.NewError("Error exchanging messages with topic resolution server", err)
	}
	return nil
}

func (app *App) New(message *Message.Message) (string, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	id := message.GetPayload()
	if _, ok := app.spawnedClients[id]; ok {
		return "", Utilities.NewError("Client "+id+" already exists", nil)
	}
	newClient := Client.New(id, app.client.GetTopicResolutionServerAddress(), app.client.GetLogger(), nil)
	chessApp, err := appChess.New(newClient, nil)
	if err != nil {
		return "", Utilities.NewError("Error creating app "+id, err)
	}
	newClient.SetApplication(chessApp)
	brokerNetConn, err := Utilities.TlsDial("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"))
	if err != nil {
		return "", Utilities.NewError("Error dialing brokerChess", err)
	}
	_, err = Utilities.TcpExchange(brokerNetConn, Message.NewAsync("addAsyncTopic", app.client.GetName(), id), 5000)
	if err != nil {
		return "", Utilities.NewError("Error exchanging messages with broker", err)
	}
	resolverNetConn, err := Utilities.TlsDial("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"))
	if err != nil {
		_, err := Utilities.TcpExchange(brokerNetConn, Message.NewAsync("removeAsyncTopic", app.client.GetName(), id), 5000)
		if err != nil {
			app.client.GetLogger().Log(Utilities.NewError("Error exchanging messages with broker", err).Error())
		}
		return "", Utilities.NewError("Error dialing topic resolution server", err)
	}
	_, err = Utilities.TcpExchange(resolverNetConn, Message.NewAsync("registerTopics", app.client.GetName(), "brokerChess"+" "+id), 5000)
	if err != nil {
		_, err := Utilities.TcpExchange(brokerNetConn, Message.NewAsync("removeAsyncTopic", app.client.GetName(), id), 5000)
		if err != nil {
			app.client.GetLogger().Log(Utilities.NewError("Error exchanging messages with broker", err).Error())
		}
		return "", Utilities.NewError("Error exchanging messages with topic resolution server", err)
	}
	err = newClient.Start()
	if err != nil {
		_, err := Utilities.TcpExchange(brokerNetConn, Message.NewAsync("removeAsyncTopic", app.client.GetName(), id), 5000)
		if err != nil {
			app.client.GetLogger().Log(Utilities.NewError("Error exchanging messages with broker", err).Error())
		}
		_, err = Utilities.TcpExchange(resolverNetConn, Message.NewAsync("unregisterTopics", app.client.GetName(), "brokerChess"+" "+id), 5000)
		if err != nil {
			app.client.GetLogger().Log(Utilities.NewError("Error exchanging messages with topic resolution server", err).Error())
		}
		return "", Utilities.NewError("Error starting client", err)
	}
	app.spawnedClients[id] = newClient
	return id, nil
}
