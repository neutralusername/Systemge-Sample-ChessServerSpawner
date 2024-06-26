package main

import (
	"Systemge/Module"
	"Systemge/Node"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/appSpawner"
	"SystemgeSampleChessServer/appWebsocketHTTP"
)

const RESOLVER_ADDRESS = "127.0.0.1:60000"
const RESOLVER_NAME_INDICATION = "127.0.0.1"
const RESOLVER_TLS_CERT_PATH = "MyCertificate.crt"
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
	nodeSpawner := Module.NewNode(&Node.Config{
		Name:                   "nodeSpawner",
		ResolverAddress:        RESOLVER_ADDRESS,
		ResolverNameIndication: RESOLVER_NAME_INDICATION,
		ResolverTLSCert:        Utilities.GetFileContent(RESOLVER_TLS_CERT_PATH),
		LoggerPath:             ERROR_LOG_FILE_PATH,
	}, appSpawner.New(), nil, nil)
	applicationWebsocketHTTP := appWebsocketHTTP.New()
	nodeWebsocketHTTP := Module.NewNode(&Node.Config{
		Name:                   "nodeWebsocketHTTP",
		ResolverAddress:        RESOLVER_ADDRESS,
		ResolverNameIndication: RESOLVER_NAME_INDICATION,
		ResolverTLSCert:        Utilities.GetFileContent(RESOLVER_TLS_CERT_PATH),
		LoggerPath:             ERROR_LOG_FILE_PATH,
		WebsocketPattern:       "/ws",
		WebsocketPort:          WEBSOCKET_PORT,
		HTTPPort:               HTTP_PORT,
	}, applicationWebsocketHTTP, applicationWebsocketHTTP, applicationWebsocketHTTP)
	Module.StartCommandLineInterface(Module.NewMultiModule(
		nodeSpawner,
		nodeWebsocketHTTP,
	))
}
