package appWebsocketHTTP

import (
	"SystemgeSampleChessServer/dto"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
	"github.com/neutralusername/Systemge/Spawner"
	"github.com/neutralusername/Systemge/Tools"
)

func (app *AppWebsocketHTTP) GetWebsocketMessageHandlers() map[string]Node.WebsocketMessageHandler {
	return map[string]Node.WebsocketMessageHandler{
		"startGame": func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			app.mutex.Lock()
			whiteId := websocketClient.GetId()
			blackId := message.GetPayload()
			if !node.WebsocketClientExists(blackId) {
				app.mutex.Unlock()
				return Error.New("Opponent does not exist", nil)
			}
			if blackId == whiteId {
				app.mutex.Unlock()
				return Error.New("You cannot play against yourself", nil)
			}
			if app.gameIds[whiteId] != "" {
				app.mutex.Unlock()
				return Error.New("You are already in a game", nil)
			}
			if app.gameIds[blackId] != "" {
				app.mutex.Unlock()
				return Error.New("Opponent is already in a game", nil)
			}
			app.mutex.Unlock()
			gameId := whiteId + "-" + blackId
			port := app.ports.Add(1)
			responseChannel, err := node.SyncMessage(Spawner.SPAWN_AND_START_NODE_SYNC, Helpers.JsonMarshal(&Config.NewNode{
				NodeConfig: &Config.Node{
					Name:              gameId,
					RandomizerSeed:    Tools.GetSystemTime(),
					InfoLoggerPath:    "logs.log",
					WarningLoggerPath: "logs.log",
					ErrorLoggerPath:   "logs.log",
				},
				SystemgeConfig: &Config.Systemge{
					ProcessMessagesOfEachConnectionSequentially: true,
					ProcessAllMessagesSequentially:              false,

					SyncRequestTimeoutMs:            10000,
					TcpTimeoutMs:                    5000,
					MaxConnectionAttempts:           0,
					ConnectionAttemptDelayMs:        1000,
					StopAfterOutgoingConnectionLoss: true,
					ServerConfig: &Config.TcpServer{
						Port:        uint16(port),
						TlsCertPath: "MyCertificate.crt",
						TlsKeyPath:  "MyKey.key",
					},
					Endpoint: &Config.TcpEndpoint{
						Address: "127.0.0.1:" + Helpers.IntToString(int(port)),
						TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
						Domain:  "example.com",
					},
					EndpointConfigs: []*Config.TcpEndpoint{
						{
							Address: "localhost:60001",
							TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
							Domain:  "example.com",
						},
						{
							Address: "localhost:60002",
							TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
							Domain:  "example.com",
						},
					},
					IncomingMessageByteLimit: 0,
					MaxPayloadSize:           0,
					MaxTopicSize:             0,
					MaxSyncTokenSize:         0,
					MaxNodeNameSize:          0,
				},
			}))
			if err != nil {
				return Error.New("Error spawning game", err)
			}
			response, err := responseChannel.ReceiveResponse()
			if err != nil {
				return Error.New("Error receiving game response", err)
			}
			if response.GetTopic() == Message.TOPIC_FAILURE {
				return Error.New(response.GetPayload(), nil)
			}
			return nil
		},
		"endGame": func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			app.mutex.Lock()
			gameId := app.gameIds[websocketClient.GetId()]
			app.mutex.Unlock()
			if gameId == "" {
				return Error.New("You are not in a game", nil)
			}
			responseChannel, err := node.SyncMessage(Spawner.DESPAWN_NODE_SYNC, gameId)
			if err != nil {
				return Error.New("Error sending end message", err)
			}
			response, err := responseChannel.ReceiveResponse()
			if err != nil {
				return Error.New("Error receiving end response", err)
			}
			if response.GetTopic() == Message.TOPIC_FAILURE {
				return Error.New(response.GetPayload(), nil)
			}
			return nil
		},
		"move": func(node *Node.Node, websocketClient *Node.WebsocketClient, message *Message.Message) error {
			app.mutex.Lock()
			gameId := app.gameIds[websocketClient.GetId()]
			app.mutex.Unlock()
			if gameId == "" {
				return Error.New("You are not in a game", nil)
			}
			move, err := dto.UnmarshalMove(message.GetPayload())
			if err != nil {
				return Error.New("Error unmarshalling move", err)
			}
			move.PlayerId = websocketClient.GetId()
			responseChannel, err := node.SyncMessage(gameId, Helpers.JsonMarshal(move))
			if err != nil {
				return Error.New("Error sending move message", err)
			}
			response, err := responseChannel.ReceiveResponse()
			if err != nil {
				return Error.New("Error receiving move response", err)
			}
			if response.GetTopic() == Message.TOPIC_FAILURE {
				return Error.New(response.GetPayload(), nil)
			}
			node.WebsocketGroupcast(gameId, Message.NewAsync("propagate_move", response.GetPayload()))
			return nil
		},
	}
}

func (app *AppWebsocketHTTP) OnConnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	err := websocketClient.Send(Message.NewAsync("connected", websocketClient.GetId()).Serialize())
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log(Error.New("Error sending connected message", err).Error())
		}
		websocketClient.Disconnect()
	}
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(node *Node.Node, websocketClient *Node.WebsocketClient) {
	app.mutex.Lock()
	gameId := app.gameIds[websocketClient.GetId()]
	app.mutex.Unlock()
	if gameId == "" {
		return
	}
	err := node.AsyncMessage(Spawner.DESPAWN_NODE_ASYNC, gameId)
	if err != nil {
		if errorLogger := node.GetErrorLogger(); errorLogger != nil {
			errorLogger.Log(Error.New("Error sending end message for game: "+gameId, err).Error())
		}
	}
}
