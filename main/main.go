package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Module"
	"Systemge/Node"
	"Systemge/Resolver"
	"Systemge/Spawner"
	"Systemge/TcpEndpoint"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/appChess"
	"SystemgeSampleChessServer/appWebsocketHTTP"
	"SystemgeSampleChessServer/config"
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
			Name:                      "nodeWebsocketHTTP",
			Logger:                    Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
			ResolverEndpoint:          TcpEndpoint.New(config.SERVER_IP+":"+Utilities.IntToString(config.RESOLVER_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
			SyncResponseTimeoutMs:     1000,
			TopicResolutionLifetimeMs: 10000,
			BrokerSubscribeDelayMs:    1000,
		}, appWebsocketHTTP.New()),
		Node.New(Config.Node{
			Name:                      "nodePingSpawner",
			Logger:                    Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
			ResolverEndpoint:          TcpEndpoint.New(config.SERVER_IP+":"+Utilities.IntToString(config.RESOLVER_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
			SyncResponseTimeoutMs:     1000,
			TopicResolutionLifetimeMs: 10000,
			BrokerSubscribeDelayMs:    1000,
		}, Spawner.New(Config.Application{
			HandleMessagesSequentially: false,
		}, Config.Spawner{
			SpawnedNodeLogger:      Utilities.NewLogger(ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH, ERROR_LOG_FILE_PATH),
			IsSpawnedNodeTopicSync: true,
			ResolverEndpoint:       TcpEndpoint.New(config.SERVER_IP+":"+Utilities.IntToString(config.RESOLVER_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
			BrokerConfigEndpoint:   TcpEndpoint.New(config.SERVER_IP+":"+Utilities.IntToString(config.BROKER_CONFIG_PORT), config.SERVER_NAME_INDICATION, Utilities.GetFileContent(config.CERT_PATH)),
		}, appChess.New)),
	))
}
