package main

import (
	"SystemgeSampleChessServer/appChess"
	"SystemgeSampleChessServer/appWebsocketHTTP"
	"SystemgeSampleChessServer/topics"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Node"
	"github.com/neutralusername/Systemge/Spawner"
	"github.com/neutralusername/Systemge/Tools"
)

const RESOLVER_ADDRESS = "127.0.0.1:60000"
const RESOLVER_NAME_INDICATION = "127.0.0.1"
const RESOLVER_TLS_CERT_PATH = "MyCertificate.crt"
const WEBSOCKET_PORT = ":8443"
const HTTP_PORT = ":8080"

const LOGGER_PATH = "logs.log"

func main() {
	Tools.NewLoggerQueue(LOGGER_PATH, 10000)
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
		NodeSpawnerCounterIntervalMs:   1000,
		NodeHTTPCounterIntervalMs:      1000,
		GoroutineUpdateIntervalMs:      1000,
		AutoStart:                      true,
		AddDashboardToDashboard:        true,
	},
		Node.New(&Config.Node{
			Name:              "nodeResolver",
			RandomizerSeed:    Tools.GetSystemTime(),
			InfoLoggerPath:    LOGGER_PATH,
			WarningLoggerPath: LOGGER_PATH,
			ErrorLoggerPath:   LOGGER_PATH,
		}, Node.NewResolverApplication(&Config.Resolver{
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
			Name:              "nodeBrokerSpawner",
			RandomizerSeed:    Tools.GetSystemTime(),
			InfoLoggerPath:    LOGGER_PATH,
			WarningLoggerPath: LOGGER_PATH,
			ErrorLoggerPath:   LOGGER_PATH,
		}, Node.NewBrokerApplication(&Config.Broker{
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
			SyncTopics:            []string{Spawner.STOP_NODE_SYNC, Spawner.START_NODE_SYNC, Spawner.SPAWN_NODE_SYNC, Spawner.DESPAWN_NODE_SYNC},
			AsyncTopics:           []string{Spawner.STOP_NODE_ASYNC, Spawner.START_NODE_ASYNC, Spawner.SPAWN_NODE_ASYNC, Spawner.DESPAWN_NODE_ASYNC},
			SyncResponseTimeoutMs: 10000,
			TcpTimeoutMs:          5000,
		})),
		Node.New(&Config.Node{
			Name:              "nodeBrokerWebsocketHTTP",
			RandomizerSeed:    Tools.GetSystemTime(),
			InfoLoggerPath:    LOGGER_PATH,
			WarningLoggerPath: LOGGER_PATH,
			ErrorLoggerPath:   LOGGER_PATH,
		}, Node.NewBrokerApplication(&Config.Broker{
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
			Name:              "nodeBrokerChess",
			RandomizerSeed:    Tools.GetSystemTime(),
			InfoLoggerPath:    LOGGER_PATH,
			WarningLoggerPath: LOGGER_PATH,
			ErrorLoggerPath:   LOGGER_PATH,
		}, Node.NewBrokerApplication(&Config.Broker{
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
			Name:              "nodeSpawner",
			RandomizerSeed:    Tools.GetSystemTime(),
			InfoLoggerPath:    LOGGER_PATH,
			WarningLoggerPath: LOGGER_PATH,
			ErrorLoggerPath:   LOGGER_PATH,
		}, Spawner.New(&Config.Spawner{
			InfoLoggerPath:              LOGGER_PATH,
			WarningLoggerPath:           LOGGER_PATH,
			ErrorLoggerPath:             LOGGER_PATH,
			IsSpawnedNodeTopicSync:      true,
			PropagateSpawnedNodeChanges: true,
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
			Name:              "nodeWebsocketHTTP",
			RandomizerSeed:    Tools.GetSystemTime(),
			InfoLoggerPath:    LOGGER_PATH,
			WarningLoggerPath: LOGGER_PATH,
			ErrorLoggerPath:   LOGGER_PATH,
		}, appWebsocketHTTP.New()),
	)).StartBlocking()
}
