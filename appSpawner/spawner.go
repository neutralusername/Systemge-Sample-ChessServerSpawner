package appSpawner

import (
	"Systemge/Module"
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
	err = app.client.RemoveSyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		app.client.GetLogger().Log(Utilities.NewError("Error removing sync topic \""+id+"\"", err).Error())
	}
	err = app.client.RemoveResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		app.client.GetLogger().Log(Utilities.NewError("Error unregistering topic \""+id+"\"", err).Error())
	}
	return nil
}

func (app *App) StartClient(id string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if _, ok := app.spawnedClients[id]; ok {
		return Utilities.NewError("Client "+id+" already exists", nil)
	}
	newClient := Module.NewClient(&Module.ClientConfig{
		Name:                   id,
		ResolverAddress:        app.client.GetResolverResolution().GetAddress(),
		ResolverNameIndication: app.client.GetResolverResolution().GetServerNameIndication(),
		ResolverTLSCertPath:    "MyCertificate.crt",
		LoggerPath:             "error.log",
	}, appChess.New, nil)
	err := app.client.AddSyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		return Utilities.NewError("Error adding sync topic \""+id+"\"", err)
	}
	err = app.client.AddResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), "brokerChess", id)
	if err != nil {
		err = app.client.RemoveSyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			app.client.GetLogger().Log(Utilities.NewError("Error removing sync topic \""+id+"\"", err).Error())
		}
		return Utilities.NewError("Error registering topic", err)
	}
	err = newClient.Start()
	if err != nil {
		err = app.client.RemoveSyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			app.client.GetLogger().Log(Utilities.NewError("Error removing sync topic \""+id+"\"", err).Error())
		}
		err = app.client.RemoveResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			app.client.GetLogger().Log(Utilities.NewError("Error unregistering topic \""+id+"\"", err).Error())
		}
		return Utilities.NewError("Error starting client", err)
	}
	app.spawnedClients[id] = newClient
	return nil
}
