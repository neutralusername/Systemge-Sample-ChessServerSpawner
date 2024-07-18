package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Helpers"
	"Systemge/Node"
	"Systemge/Resolver"
	"Systemge/Spawner"
	"Systemge/Tcp"
	"SystemgeSampleChessServer/appChess"
	"SystemgeSampleChessServer/appWebsocketHTTP"
	"SystemgeSampleChessServer/config"
	"SystemgeSampleChessServer/topics"
)

const RESOLVER_ADDRESS = "127.0.0.1:60000"
const RESOLVER_NAME_INDICATION = "127.0.0.1"
const RESOLVER_TLS_CERT_PATH = "MyCertificate.crt"
const WEBSOCKET_PORT = ":8443"
const HTTP_PORT = ":8080"

const ERROR_LOG_FILE_PATH = "error.log"

func main() {
	Node.StartCommandLineInterface(true,
		Node.New(Config.Node{
			Name: "nodeResolver",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, Resolver.New(Config.Resolver{
			Server:       Tcp.NewServer(60000, "MyCertificate.crt", "MyKey.key"),
			ConfigServer: Tcp.NewServer(60001, "MyCertificate.crt", "MyKey.key"),

			TcpTimeoutMs: 5000,
		})),
		Node.New(Config.Node{
			Name: "nodeBrokerSpawner",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, Broker.New(Config.Broker{
			Server:       Tcp.NewServer(60002, "MyCertificate.crt", "MyKey.key"),
			Endpoint:     Tcp.NewEndpoint("127.0.0.1:60002", "example.com", Helpers.GetFileContent("MyCertificate.crt")),
			ConfigServer: Tcp.NewServer(60003, "MyCertificate.crt", "MyKey.key"),

			ResolverConfigEndpoint: Tcp.NewEndpoint("127.0.0.1:60001", "example.com", Helpers.GetFileContent("MyCertificate.crt")),

			SyncTopics:  []string{topics.END_NODE_SYNC, topics.START_NODE_SYNC},
			AsyncTopics: []string{topics.END_NODE_ASYNC, topics.START_NODE_ASYNC},

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(Config.Node{
			Name: "nodeBrokerWebsocketHTTP",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, Broker.New(Config.Broker{
			Server:       Tcp.NewServer(60004, "MyCertificate.crt", "MyKey.key"),
			Endpoint:     Tcp.NewEndpoint("127.0.0.1:60004", "example.com", Helpers.GetFileContent("MyCertificate.crt")),
			ConfigServer: Tcp.NewServer(60005, "MyCertificate.crt", "MyKey.key"),

			ResolverConfigEndpoint: Tcp.NewEndpoint("127.0.0.1:60001", "example.com", Helpers.GetFileContent("MyCertificate.crt")),

			SyncTopics:  []string{topics.PROPAGATE_GAMESTART},
			AsyncTopics: []string{topics.PROPAGATE_GAMEEND},

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(Config.Node{
			Name: "nodeBrokerChess",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, Broker.New(Config.Broker{
			Server:       Tcp.NewServer(60006, "MyCertificate.crt", "MyKey.key"),
			Endpoint:     Tcp.NewEndpoint("127.0.0.1:60006", "example.com", Helpers.GetFileContent("MyCertificate.crt")),
			ConfigServer: Tcp.NewServer(60007, "MyCertificate.crt", "MyKey.key"),

			ResolverConfigEndpoint: Tcp.NewEndpoint("127.0.0.1:60001", "example.com", Helpers.GetFileContent("MyCertificate.crt")),

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(Config.Node{
			Name: "nodeSpawner",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, Spawner.New(Config.Spawner{
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
			IsSpawnedNodeTopicSync: true,
			ResolverEndpoint:       Tcp.NewEndpoint(config.SERVER_IP+":"+Helpers.IntToString(config.RESOLVER_PORT), config.SERVER_NAME_INDICATION, Helpers.GetFileContent(config.CERT_PATH)),
			BrokerConfigEndpoint:   Tcp.NewEndpoint(config.SERVER_IP+":"+Helpers.IntToString(60003), config.SERVER_NAME_INDICATION, Helpers.GetFileContent(config.CERT_PATH)),
		}, Config.Systemge{
			HandleMessagesSequentially: false,

			BrokerSubscribeDelayMs:    1000,
			TopicResolutionLifetimeMs: 10000,
			SyncResponseTimeoutMs:     10000,
			TcpTimeoutMs:              5000,

			ResolverEndpoint: Tcp.NewEndpoint("127.0.0.1:60000", "example.com", Helpers.GetFileContent("MyCertificate.crt")),
		},
			appChess.New)),
		Node.New(Config.Node{
			Name: "nodeWebsocketHTTP",
			Logger: Config.Logger{
				InfoPath:    ERROR_LOG_FILE_PATH,
				DebugPath:   ERROR_LOG_FILE_PATH,
				ErrorPath:   ERROR_LOG_FILE_PATH,
				WarningPath: ERROR_LOG_FILE_PATH,
				QueueBuffer: 10000,
			},
		}, appWebsocketHTTP.New()),
	)
}
