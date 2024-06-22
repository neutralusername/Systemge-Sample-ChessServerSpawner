package main

import (
	"Systemge/Module"
	"SystemgeSampleChessServer/appSpawner"
	"SystemgeSampleChessServer/appWebsocketHTTP"
)

const TOPICRESOLUTIONSERVER_ADDRESS = "127.0.0.1:60000"
const WEBSOCKET_PORT = ":8443"
const HTTP_PORT = ":8080"

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
	clientSpawner := Module.NewClient(&Module.ClientConfig{
		Name:            "clientSpawner",
		ResolverAddress: TOPICRESOLUTIONSERVER_ADDRESS,
		LoggerPath:      ERROR_LOG_FILE_PATH,
	}, appSpawner.New, nil)
	clientWebsocketHTTP := Module.NewCompositeClientWebsocketHTTP(&Module.ClientConfig{
		Name:             "clientWebsocket",
		ResolverAddress:  TOPICRESOLUTIONSERVER_ADDRESS,
		LoggerPath:       ERROR_LOG_FILE_PATH,
		WebsocketPattern: "/ws",
		WebsocketPort:    WEBSOCKET_PORT,
		WebsocketCert:    "",
		WebsocketKey:     "",
		HTTPPort:         HTTP_PORT,
		HTTPCert:         "",
		HTTPKey:          "",
	}, appWebsocketHTTP.New, nil)
	Module.StartCommandLineInterface(Module.NewMultiModule(
		clientSpawner,
		clientWebsocketHTTP,
	), clientSpawner.GetApplication().GetCustomCommandHandlers(), clientWebsocketHTTP.GetApplication().GetCustomCommandHandlers(), clientWebsocketHTTP.GetWebsocketServer().GetCustomCommandHandlers())
}
