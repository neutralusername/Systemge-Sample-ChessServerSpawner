package appSpawner

import (
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/appChess"
)

func (app *App) EndClient(id string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
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
	_, err = Utilities.TcpExchange(brokerNetConn, Message.NewAsync("removeAsyncTopic", app.client.GetName(), id), 5000)
	if err != nil {
		return Utilities.NewError("Error exchanging messages with broker", err)
	}
	resolverNetConn, err := Utilities.TlsDial("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"))
	if err != nil {
		return Utilities.NewError("Error dialing topic resolution server", err)
	}
	_, err = Utilities.TcpExchange(resolverNetConn, Message.NewAsync("unregisterTopics", app.client.GetName(), "brokerChess"+" "+id), 5000)
	if err != nil {
		return Utilities.NewError("Error exchanging messages with topic resolution server", err)
	}
	return nil
}

func (app *App) StartClient(id string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if _, ok := app.spawnedClients[id]; ok {
		return Utilities.NewError("Client "+id+" already exists", nil)
	}
	newClient := Client.New(id, app.client.GetTopicResolutionServerAddress(), app.client.GetLogger(), nil)
	chessApp, err := appChess.New(newClient, nil)
	if err != nil {
		return Utilities.NewError("Error creating app "+id, err)
	}
	newClient.SetApplication(chessApp)
	brokerNetConn, err := Utilities.TlsDial("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"))
	if err != nil {
		return Utilities.NewError("Error dialing brokerChess", err)
	}
	_, err = Utilities.TcpExchange(brokerNetConn, Message.NewAsync("addSyncTopic", app.client.GetName(), id), 5000)
	if err != nil {
		return Utilities.NewError("Error exchanging messages with broker", err)
	}
	resolverNetConn, err := Utilities.TlsDial("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"))
	if err != nil {
		_, err := Utilities.TcpExchange(brokerNetConn, Message.NewAsync("removeSyncTopic", app.client.GetName(), id), 5000)
		if err != nil {
			app.client.GetLogger().Log(Utilities.NewError("Error exchanging messages with broker", err).Error())
		}
		return Utilities.NewError("Error dialing topic resolution server", err)
	}
	_, err = Utilities.TcpExchange(resolverNetConn, Message.NewAsync("registerTopics", app.client.GetName(), "brokerChess"+" "+id), 5000)
	if err != nil {
		_, err := Utilities.TcpExchange(brokerNetConn, Message.NewAsync("removeAsyncTopic", app.client.GetName(), id), 5000)
		if err != nil {
			app.client.GetLogger().Log(Utilities.NewError("Error exchanging messages with broker", err).Error())
		}
		return Utilities.NewError("Error exchanging messages with topic resolution server", err)
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
		return Utilities.NewError("Error starting client", err)
	}
	app.spawnedClients[id] = newClient
	return nil
}
