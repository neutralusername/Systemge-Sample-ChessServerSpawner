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
	loggerQueue := Tools.NewLoggerQueue(LOGGER_PATH, 10000)
	Node.New(&Config.Node{
		Name:           "dashboard",
		RandomizerSeed: Tools.GetSystemTime(),
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
			InfoLogger:     Tools.NewLogger("[Info \"nodeResolver\"]", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeResolver\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeResolver\"] ", loggerQueue),
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
			InfoLogger:     Tools.NewLogger("[Info \"nodeBrokerSpawner\"]", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeBrokerSpawner\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeBrokerSpawner\"] ", loggerQueue),
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
			InfoLogger:     Tools.NewLogger("[Info \"nodeBrokerWebsocketHTTP\"]", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeBrokerWebsocketHTTP\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeBrokerWebsocketHTTP\"] ", loggerQueue),
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
			InfoLogger:     Tools.NewLogger("[Info \"nodeBrokerChess\"]", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeBrokerChess\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeBrokerChess\"] ", loggerQueue),
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
			InfoLogger:     Tools.NewLogger("[Info \"nodeSpawner\"]", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeSpawner\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeSpawner\"] ", loggerQueue),
		}, Spawner.New(&Config.Spawner{
			InfoLogger:             Tools.NewLogger("[Info \"spawnedNode\"]", loggerQueue),
			WarningLogger:          Tools.NewLogger("[Warning \"spawnedNode\"] ", loggerQueue),
			ErrorLogger:            Tools.NewLogger("[Error \"spawnedNode\"] ", loggerQueue),
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
			InfoLogger:     Tools.NewLogger("[Info \"nodeWebsocketHTTP\"]", loggerQueue),
			WarningLogger:  Tools.NewLogger("[Warning \"nodeWebsocketHTTP\"] ", loggerQueue),
			ErrorLogger:    Tools.NewLogger("[Error \"nodeWebsocketHTTP\"] ", loggerQueue),
		}, appWebsocketHTTP.New()),
	)).StartBlocking()
}
