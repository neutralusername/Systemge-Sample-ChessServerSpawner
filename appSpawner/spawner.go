package appSpawner

import (
	"Systemge/Error"
	"Systemge/Module"
	"Systemge/Node"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/appChess"
)

func (app *App) EndNode(node *Node.Node, id string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	spawnedNode := app.spawnedNodes[id]
	if spawnedNode == nil {
		return Error.New("Node "+id+" does not exist", nil)
	}
	err := spawnedNode.Stop()
	if err != nil {
		return Error.New("Error stopping node "+id, err)
	}
	delete(app.spawnedNodes, id)
	err = node.RemoveSyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		node.GetLogger().Log(Error.New("Error removing sync topic \""+id+"\"", err).Error())
	}
	err = node.RemoveResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		node.GetLogger().Log(Error.New("Error unregistering topic \""+id+"\"", err).Error())
	}
	return nil
}

func (app *App) StartNode(node *Node.Node, id string) error {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if _, ok := app.spawnedNodes[id]; ok {
		return Error.New("Node "+id+" already exists", nil)
	}
	newNode := Module.NewNode(&Node.NodeConfig{
		Name:       id,
		LoggerPath: "error.log",
	}, appChess.New(id), nil, nil)
	err := node.AddSyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
	if err != nil {
		return Error.New("Error adding sync topic \""+id+"\"", err)
	}
	err = node.AddResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), "brokerChess", id)
	if err != nil {
		err = node.RemoveSyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			node.GetLogger().Log(Error.New("Error removing sync topic \""+id+"\"", err).Error())
		}
		return Error.New("Error registering topic", err)
	}
	err = newNode.Start()
	if err != nil {
		err = node.RemoveSyncTopicRemotely("127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			node.GetLogger().Log(Error.New("Error removing sync topic \""+id+"\"", err).Error())
		}
		err = node.RemoveResolverTopicsRemotely("127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("./MyCertificate.crt"), id)
		if err != nil {
			node.GetLogger().Log(Error.New("Error unregistering topic \""+id+"\"", err).Error())
		}
		return Error.New("Error starting node", err)
	}
	app.spawnedNodes[id] = newNode
	return nil
}
