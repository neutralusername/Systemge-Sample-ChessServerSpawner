package main

import (
	"Systemge/Module"
	"SystemgeSampleChessServer/appSpawner"
	"SystemgeSampleChessServer/appWebsocket"
)

const TOPICRESOLUTIONSERVER_ADDRESS = "127.0.0.1:60000"
const WEBSOCKET_PORT = ":8443"

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	err := Module.NewResolverFromConfig("resolver.systemge", ERROR_LOG_FILE_PATH).Start()
	if err != nil {
		panic(err)
	}
	err = Module.NewBrokerFromConfig("brokerSpawner.systemge", ERROR_LOG_FILE_PATH).Start()
	if err != nil {
		panic(err)
	}
	err = Module.NewBrokerFromConfig("brokerWebsocket.systemge", ERROR_LOG_FILE_PATH).Start()
	if err != nil {
		panic(err)
	}
	err = Module.NewBrokerFromConfig("brokerChess.systemge", ERROR_LOG_FILE_PATH).Start()
	if err != nil {
		panic(err)
	}
	clientSpawner := Module.NewClient("clientSpawner", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, appSpawner.New, nil)
	clientWebsocket := Module.NewWebsocketClient("clientWebsocket", TOPICRESOLUTIONSERVER_ADDRESS, ERROR_LOG_FILE_PATH, "/ws", WEBSOCKET_PORT, "", "", appWebsocket.New, nil)
	Module.StartCommandLineInterface(Module.NewMultiModule(
		clientSpawner,
		clientWebsocket,
		Module.NewHTTPServerFromConfig("httpServe.systemge", ERROR_LOG_FILE_PATH),
	), clientSpawner.GetApplication().GetCustomCommandHandlers(), clientWebsocket.GetApplication().GetCustomCommandHandlers(), clientWebsocket.GetWebsocketServer().GetCustomCommandHandlers())
}
