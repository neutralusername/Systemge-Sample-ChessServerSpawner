package appSpawner

import (
	"Systemge/Error"
	"Systemge/Module"
	"Systemge/Node"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/appChess"
)

func (app *App) EndClient(client *Node.Node, id string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	spawnedClient := app.spawnedClients[id]
	if spawnedClient == nil {
		return Error.New("Node "+id+" does not exist", nil)
	}
	err := spawnedClient.Stop()
	if err != nil {
		return Error.New("Error stopping client "+id, err)
	}
	delete(app.spawnedClients, id)
	err = client.RemoveSyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		client.GetLogger().Log(Error.New("Error removing sync topic \""+id+"\"", err).Error())
	}
	err = client.RemoveResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		client.GetLogger().Log(Error.New("Error unregistering topic \""+id+"\"", err).Error())
	}
	return nil
}

func (app *App) StartClient(client *Node.Node, id string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if _, ok := app.spawnedClients[id]; ok {
		return Error.New("Node "+id+" already exists", nil)
	}
	newClient := Module.NewClient(&Node.Config{
		Name:                   id,
		ResolverAddress:        client.GetResolverAddress(),
		ResolverNameIndication: client.GetResolverNameIndication(),
		ResolverTLSCert:        client.GetResolverTLSCert(),
		LoggerPath:             "error.log",
	}, appChess.New(id), nil, nil)
	err := client.AddSyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		return Error.New("Error adding sync topic \""+id+"\"", err)
	}
	err = client.AddResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), "brokerChess", id)
	if err != nil {
		err = client.RemoveSyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			client.GetLogger().Log(Error.New("Error removing sync topic \""+id+"\"", err).Error())
		}
		return Error.New("Error registering topic", err)
	}
	err = newClient.Start()
	if err != nil {
		err = client.RemoveSyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			client.GetLogger().Log(Error.New("Error removing sync topic \""+id+"\"", err).Error())
		}
		err = client.RemoveResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			client.GetLogger().Log(Error.New("Error unregistering topic \""+id+"\"", err).Error())
		}
		return Error.New("Error starting client", err)
	}
	app.spawnedClients[id] = newClient
	return nil
}
