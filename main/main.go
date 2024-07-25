package main

import (
	"Systemge/Broker"
	"Systemge/Config"
	"Systemge/Dashboard"
	"Systemge/Helpers"
	"Systemge/Node"
	"Systemge/Resolver"
	"Systemge/Spawner"
	"Systemge/Tools"
	"SystemgeSampleChessServer/appChess"
	"SystemgeSampleChessServer/appWebsocketHTTP"
	"SystemgeSampleChessServer/topics"
)

const RESOLVER_ADDRESS = "127.0.0.1:60000"
const RESOLVER_NAME_INDICATION = "127.0.0.1"
const RESOLVER_TLS_CERT_PATH = "MyCertificate.crt"
const WEBSOCKET_PORT = ":8443"
const HTTP_PORT = ":8080"

const LOGGER_PATH = "logs.log"

func main() {

	Node.New(&Config.Node{
		Name:           "dashboard",
		RandomizerSeed: Tools.GetSystemTime(),
		/* ErrorLogger: &Config.Logger{
			Path:        LOGGER_PATH,
			QueueBuffer: 10000,
			Prefix:      "[Error \"dashboard\"] ",
		},
		WarningLogger: &Config.Logger{
			Path:        LOGGER_PATH,
			QueueBuffer: 10000,
			Prefix:      "[Warning \"dashboard\"] ",
		},
		InfoLogger: &Config.Logger{
			Path:        LOGGER_PATH,
			QueueBuffer: 10000,
			Prefix:      "[Info \"dashboard\"] ",
		},
		DebugLogger: &Config.Logger{
			Path:        LOGGER_PATH,
			QueueBuffer: 10000,
			Prefix:      "[Debug \"dashboard\"] ",
		}, */
	}, Dashboard.New(&Config.Dashboard{
		Server: &Config.TcpServer{
			Port: 8081,
		},
		NodeStatusIntervalMs:           1000,
		NodeSystemgeCounterIntervalMs:  1000,
		NodeWebsocketCounterIntervalMs: 1000,
		NodeBrokerCounterIntervalMs:    1000,
		NodeResolverCounterIntervalMs:  1000,
		HeapUpdateIntervalMs:           1000,
		AutoStart:                      true,
		AddSpawnedNodesToDashboard:     true,
	},
		Node.New(&Config.Node{
			Name:           "nodeResolver",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Info \"nodeResolver\"] ",
			},
			WarningLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Warning \"nodeResolver\"] ",
			},
			ErrorLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Error \"nodeResolver\"] ",
			},
			DebugLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Debug \"nodeResolver\"] ",
			},
		}, Resolver.New(&Config.Resolver{
			Server: &Config.TcpServer{
				Port:        60000,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			ConfigServer: &Config.TcpServer{
				Port:        60001,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},

			TcpTimeoutMs: 5000,
		})),
		Node.New(&Config.Node{
			Name:           "nodeBrokerSpawner",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Info \"nodeBrokerSpawner\"] ",
			},
			WarningLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Warning \"nodeBrokerSpawner\"] ",
			},
			ErrorLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Error \"nodeBrokerSpawner\"] ",
			},
			DebugLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Debug \"nodeBrokerSpawner\"] ",
			},
		}, Broker.New(&Config.Broker{
			Server: &Config.TcpServer{
				Port:        60002,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			Endpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60002",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			ConfigServer: &Config.TcpServer{
				Port:        60003,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			ResolverConfigEndpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60001",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			SyncTopics:            []string{topics.STOP_NODE_SYNC, topics.START_NODE_SYNC, topics.SPAWN_NODE_SYNC, topics.DESPAWN_NODE_SYNC},
			AsyncTopics:           []string{topics.STOP_NODE_ASYNC, topics.START_NODE_ASYNC, topics.SPAWN_NODE_ASYNC, topics.DESPAWN_NODE_ASYNC},
			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(&Config.Node{
			Name:           "nodeBrokerWebsocketHTTP",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Info \"nodeBrokerWebsocketHTTP\"] ",
			},
			WarningLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Warning \"nodeBrokerWebsocketHTTP\"] ",
			},
			ErrorLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Error \"nodeBrokerWebsocketHTTP\"] ",
			},
			DebugLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Debug \"nodeBrokerWebsocketHTTP\"] ",
			},
		}, Broker.New(&Config.Broker{
			Server: &Config.TcpServer{
				Port:        60004,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			Endpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60004",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			ConfigServer: &Config.TcpServer{
				Port:        60005,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			ResolverConfigEndpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60001",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			SyncTopics:  []string{topics.PROPAGATE_GAMESTART},
			AsyncTopics: []string{topics.PROPAGATE_GAMEEND},

			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(&Config.Node{
			Name:           "nodeBrokerChess",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Info \"nodeBrokerChess\"] ",
			},
			WarningLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Warning \"nodeBrokerChess\"] ",
			},
			ErrorLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Error \"nodeBrokerChess\"] ",
			},
			DebugLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Debug \"nodeBrokerChess\"] ",
			},
		}, Broker.New(&Config.Broker{
			Server: &Config.TcpServer{
				Port:        60006,
				TlsCertPath: "MyCertificate.crt",
				TlsKeyPath:  "MyKey.key",
			},
			Endpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60006",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			ConfigServer: &Config.TcpServer{Port: 60007, TlsCertPath: "MyCertificate.crt", TlsKeyPath: "MyKey.key"},
			ResolverConfigEndpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60001",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(&Config.Node{
			Name:           "nodeSpawner",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Info \"nodeSpawner\"] ",
			},
			WarningLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Warning \"nodeSpawner\"] ",
			},
			ErrorLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Error \"nodeSpawner\"] ",
			},
			DebugLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Debug \"nodeSpawner\"] ",
			},
		}, Spawner.New(&Config.Spawner{
			InfoLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Info \"spawnedNode\"] ",
			},
			WarningLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Warning \"spawnedNode\"] ",
			},
			ErrorLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Error \"spawnedNode\"] ",
			},
			DebugLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Debug \"spawnedNode\"] ",
			},
			IsSpawnedNodeTopicSync: true,
			ResolverEndpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60000",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
			BrokerConfigEndpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60007",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
		}, &Config.Systemge{
			HandleMessagesSequentially: false,

			BrokerSubscribeDelayMs:    1000,
			TopicResolutionLifetimeMs: 10000,
			SyncResponseTimeoutMs:     10000,
			TcpTimeoutMs:              5000,

			ResolverEndpoint: &Config.TcpEndpoint{
				Address: "127.0.0.1:60000",
				Domain:  "example.com",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
			},
		},
			appChess.New)),
		Node.New(&Config.Node{
			Name:           "nodeWebsocketHTTP",
			RandomizerSeed: Tools.GetSystemTime(),
			InfoLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Info \"nodeWebsocketHTTP\"] ",
			},
			WarningLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Warning \"nodeWebsocketHTTP\"] ",
			},
			ErrorLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Error \"nodeWebsocketHTTP\"] ",
			},
			DebugLogger: &Config.Logger{
				Path:        LOGGER_PATH,
				QueueBuffer: 10000,
				Prefix:      "[Debug \"nodeWebsocketHTTP\"] ",
			},
		}, appWebsocketHTTP.New()),
	)).StartBlocking()
}
