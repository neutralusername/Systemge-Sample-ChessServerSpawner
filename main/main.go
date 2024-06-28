package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Module"
	"Systemge/Node"
	"Systemge/Resolution"
	"Systemge/Resolver"
	"Systemge/Spawner"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/appChess"
	"SystemgeSampleChessServer/appWebsocketHTTP"
)

const RESOLVER_ADDRESS = "127.0.0.1:60000"
const RESOLVER_NAME_INDICATION = "127.0.0.1"
const RESOLVER_TLS_CERT_PATH = "MyCertificate.crt"
const WEBSOCKET_PORT = ":8443"
const HTTP_PORT = ":8080"

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	err := Resolver.New(Config.ParseResolverConfigFromFile("resolver.systemge")).Start()
	if err != nil {
		panic(err)
	}
	err = Broker.New(Config.ParseBrokerConfigFromFile("brokerSpawner.systemge")).Start()
	if err != nil {
		panic(err)
	}
	err = Broker.New(Config.ParseBrokerConfigFromFile("brokerWebsocketHTTP.systemge")).Start()
	if err != nil {
		panic(err)
	}
	err = Broker.New(Config.ParseBrokerConfigFromFile("brokerChess.systemge")).Start()
	if err != nil {
		panic(err)
	}
	Module.StartCommandLineInterface(Module.NewMultiModule(
		Node.New(Config.Node{
			Name:       "nodePingSpawner",
			LoggerPath: ERROR_LOG_FILE_PATH,
		}, Spawner.New(Config.Application{
			ResolverResolution:         Resolution.New("resolver", "127.0.0.1:60000", "127.0.0.1", Utilities.GetFileContent("MyCertificate.crt")),
			HandleMessagesSequentially: false,
		}, Config.Spawner{
			IsSpawnedNodeTopicSync:       true,
			SpawnedNodeLoggerPath:        ERROR_LOG_FILE_PATH,
			ResolverConfigResolution:     Resolution.New("resolverConfig", "127.0.0.1:60001", "127.0.0.1", Utilities.GetFileContent("MyCertificate.crt")),
			BrokerConfigResolution:       Resolution.New("pingBrokerConfig", "127.0.0.1:60008", "127.0.0.1", Utilities.GetFileContent("MyCertificate.crt")),
			BrokerSubscriptionResolution: Resolution.New("pingBroker", "127.0.0.1:60007", "127.0.0.1", Utilities.GetFileContent("MyCertificate.crt")),
		}, appChess.New)),
		Node.New(Config.Node{
			Name:       "nodeWebsocketHTTP",
			LoggerPath: ERROR_LOG_FILE_PATH,
		}, appWebsocketHTTP.New()),
	))
}
