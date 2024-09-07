package appWebsocketHTTP

import (
	"SystemgeSampleChessServer/dto"
	"sync"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Dashboard"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/HTTPServer"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SingleRequestServer"
	"github.com/neutralusername/Systemge/Status"
	"github.com/neutralusername/Systemge/WebsocketServer"
)

type activeGame struct {
	port    uint16
	blackId string
	whiteId string
}

type AppWebsocketHTTP struct {
	status      int
	statusMutex sync.Mutex

	websocketServer *WebsocketServer.WebsocketServer
	httpServer      *HTTPServer.HTTPServer

	websocketIdGames map[string]*activeGame
	mutex            sync.Mutex
}

func New() *AppWebsocketHTTP {
	app := &AppWebsocketHTTP{
		status:           Status.STOPPED,
		websocketIdGames: make(map[string]*activeGame),
	}

	app.websocketServer = WebsocketServer.New("appWebsocketHttp_websocketServer",
		&Config.WebsocketServer{
			ClientWatchdogTimeoutMs: 1000 * 60,
			Pattern:                 "/ws",
			TcpServerConfig: &Config.TcpServer{
				Port: 8443,
			},
		},
		WebsocketServer.MessageHandlers{
			"startGame": func(websocketClient *WebsocketServer.WebsocketClient, message *Message.Message) error {
				whiteId := websocketClient.GetId()
				blackId := message.GetPayload()
				if !app.websocketServer.ClientExists(blackId) { // can theoretically change at any time during this functions execution
					return Error.New("Opponent does not exist", nil)
				}
				if blackId == whiteId {
					return Error.New("You cannot play against yourself", nil)
				}

				app.mutex.Lock()
				if app.websocketIdGames[whiteId] != nil {
					app.mutex.Unlock()
					return Error.New("You are already in a game", nil)
				}
				if app.websocketIdGames[blackId] != nil {
					app.mutex.Unlock()
					return Error.New("Opponent is already in a game", nil)
				}
				activeGame := &activeGame{
					whiteId: whiteId,
					blackId: blackId,
				}
				app.websocketIdGames[whiteId] = activeGame
				app.websocketIdGames[blackId] = activeGame
				app.mutex.Unlock()

				response, err := SingleRequestServer.SyncRequest("appWebsocketHttp",
					&Config.SingleRequestClient{
						TcpConnectionConfig: &Config.TcpSystemgeConnection{},
						TcpClientConfig: &Config.TcpClient{
							Address: "localhost:60001",
							TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
							Domain:  "example.com",
						},
					},
					"spawn", "",
				)
				if err != nil {
					app.mutex.Lock()
					delete(app.websocketIdGames, whiteId)
					delete(app.websocketIdGames, blackId)
					app.mutex.Unlock()
					return Error.New("Error spawning game", err)
				}
				if response.GetTopic() != Message.TOPIC_SUCCESS {
					app.mutex.Lock()
					delete(app.websocketIdGames, whiteId)
					delete(app.websocketIdGames, blackId)
					app.mutex.Unlock()
					return Error.New("Error spawning game", nil)
				}
				activeGame.port = Helpers.StringToUint16(response.GetPayload())

				response, err = SingleRequestServer.SyncRequest("appWebsocketHttp",
					&Config.SingleRequestClient{
						TcpConnectionConfig: &Config.TcpSystemgeConnection{},
						TcpClientConfig: &Config.TcpClient{
							Address: "localhost:" + Helpers.Uint16ToString(activeGame.port),
							TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
							Domain:  "example.com",
						},
					},
					"getBoard", "",
				)
				if err != nil {
					// shouldn't happen in this sample. Should be properly error handled in a real application though
					panic(Error.New("Error getting board", err))
				}
				if response.GetTopic() != Message.TOPIC_SUCCESS {
					// shouldn't happen in this sample. Should be properly error handled in a real application though
					panic(Error.New("Error getting board", nil))
				}

				app.websocketServer.Multicast([]string{blackId, whiteId}, Message.NewAsync("startGame", response.GetPayload()))
				return nil
			},
			"endGame": func(websocketClient *WebsocketServer.WebsocketClient, message *Message.Message) error {
				app.mutex.Lock()
				activeGame := app.websocketIdGames[websocketClient.GetId()]
				if activeGame == nil {
					app.mutex.Unlock()
					return Error.New("You are not in a game", nil)
				}
				delete(app.websocketIdGames, activeGame.whiteId)
				delete(app.websocketIdGames, activeGame.blackId)
				app.mutex.Unlock()
				err := SingleRequestServer.AsyncMessage("appWebsocketHttp",
					&Config.SingleRequestClient{
						TcpConnectionConfig: &Config.TcpSystemgeConnection{
							TcpReceiveTimeoutMs: 1000,
						},
						TcpClientConfig: &Config.TcpClient{
							Address: "localhost:" + Helpers.Uint16ToString(activeGame.port),
							TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
							Domain:  "example.com",
						},
					},
					"end", Helpers.Uint16ToString(activeGame.port),
				)
				if err != nil {
					// shouldn't happen in this sample. Should be properly error handled in a real application though
					panic(Error.New("Error despawning game", err))
				}
				app.websocketServer.Multicast([]string{activeGame.blackId, activeGame.whiteId}, Message.NewAsync("endGame", ""))
				return nil
			},
			"move": func(websocketClient *WebsocketServer.WebsocketClient, message *Message.Message) error {
				app.mutex.Lock()
				activeGame := app.websocketIdGames[websocketClient.GetId()]
				app.mutex.Unlock()
				if activeGame == nil {
					return Error.New("You are not in a game", nil)
				}

				move, err := dto.UnmarshalMove(message.GetPayload())
				if err != nil {
					return Error.New("Error unmarshalling move", err)
				}
				if activeGame.whiteId == websocketClient.GetId() {
					move.White = true
				} else {
					move.White = false
				}
				response, err := SingleRequestServer.SyncRequest("appWebsocketHttp",
					&Config.SingleRequestClient{
						TcpConnectionConfig: &Config.TcpSystemgeConnection{},
						TcpClientConfig: &Config.TcpClient{
							Address: "localhost:" + Helpers.Uint16ToString(activeGame.port),
							TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
							Domain:  "example.com",
						},
					},
					"move", Helpers.JsonMarshal(move),
				)
				if err != nil {
					return err
				}
				if response.GetTopic() != Message.TOPIC_SUCCESS {
					return Error.New("Error making move", nil)
				}
				app.websocketServer.Multicast([]string{activeGame.blackId, activeGame.whiteId}, Message.NewAsync("move", response.GetPayload()))
				return nil
			},
		},
		app.OnConnectHandler, app.OnDisconnectHandler,
	)
	app.httpServer = HTTPServer.New("httpServer",
		&Config.HTTPServer{
			TcpServerConfig: &Config.TcpServer{
				Port: 8080,
			},
		},
		HTTPServer.Handlers{
			"/": HTTPServer.SendDirectory("../frontend"),
		},
	)
	Dashboard.NewClient("appWebsocketHttp_dashboardClient",
		&Config.DashboardClient{
			ConnectionConfig: &Config.TcpSystemgeConnection{},
			ClientConfig: &Config.TcpClient{
				Address: "localhost:60000",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
				Domain:  "example.com",
			},
		},
		app.start, app.stop, nil, app.getStatus,
		nil,
	).Start()
	if err := app.start(); err != nil {
		// shouldn't happen in this sample. Should be properly error handled in a real application though
		panic(Error.New("Failed to start appWebsocketHttp", err))
	}
	return app
}

func (app *AppWebsocketHTTP) getStatus() int {
	return app.status
}

func (app *AppWebsocketHTTP) start() error {
	app.statusMutex.Lock()
	defer app.statusMutex.Unlock()
	if app.status != Status.STOPPED {
		return Error.New("App already started", nil)
	}
	if err := app.websocketServer.Start(); err != nil {
		return Error.New("Failed to start websocketServer", err)
	}
	if err := app.httpServer.Start(); err != nil {
		app.websocketServer.Stop()
		return Error.New("Failed to start httpServer", err)
	}
	app.status = Status.STARTED
	return nil
}

func (app *AppWebsocketHTTP) stop() error {
	app.statusMutex.Lock()
	defer app.statusMutex.Unlock()
	if app.status != Status.STARTED {
		return Error.New("App not started", nil)
	}
	app.httpServer.Stop()
	app.websocketServer.Stop()
	app.status = Status.STOPPED
	return nil
}

func (app *AppWebsocketHTTP) WebsocketPropagate(message *Message.Message) {
	app.websocketServer.Broadcast(message)
}

func (app *AppWebsocketHTTP) OnConnectHandler(websocketClient *WebsocketServer.WebsocketClient) error {
	err := websocketClient.Send(Message.NewAsync("connected", websocketClient.GetId()).Serialize())
	if err != nil {
		return Error.New("Error sending connected message", err)
	}
	return nil
}

func (app *AppWebsocketHTTP) OnDisconnectHandler(websocketClient *WebsocketServer.WebsocketClient) {
	app.mutex.Lock()
	activeGame := app.websocketIdGames[websocketClient.GetId()]
	if activeGame != nil {
		delete(app.websocketIdGames, activeGame.whiteId)
		delete(app.websocketIdGames, activeGame.blackId)
		go func() {
			err := SingleRequestServer.AsyncMessage("appWebsocketHttp",
				&Config.SingleRequestClient{
					TcpConnectionConfig: &Config.TcpSystemgeConnection{},
					TcpClientConfig: &Config.TcpClient{
						Address: "localhost:" + Helpers.Uint16ToString(activeGame.port),
						TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
						Domain:  "example.com",
					},
				},
				"end", Helpers.Uint16ToString(activeGame.port),
			)
			if err != nil {
				// shouldn't happen in this sample. Should be properly error handled in a real application though
				panic(Error.New("Error despawning game", err))
			}
			app.websocketServer.Multicast([]string{activeGame.blackId, activeGame.whiteId}, Message.NewAsync("endGame", ""))
		}()
	}
	app.mutex.Unlock()

}
