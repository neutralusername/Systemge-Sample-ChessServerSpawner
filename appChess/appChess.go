package appChess

import (
	"SystemgeSampleChessServer/dto"
	"sync"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/DashboardClient"
	"github.com/neutralusername/Systemge/DashboardClientCustomService"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SingleRequestServer"
	"github.com/neutralusername/Systemge/SystemgeConnection"
)

type AppChess struct {
	board [8][8]Piece
	moves []*dto.Move
	mutex sync.Mutex

	singleRequestServer *SingleRequestServer.Server
}

func NewAppChess(port uint16, stopFunc func()) *AppChess {
	app := &AppChess{}
	if false {
		app.board = get960StartingPosition()
	} else {
		app.board = getStandardStartingPosition()
	}
	var dashboardClient *DashboardClient.Client
	app.singleRequestServer = SingleRequestServer.NewSingleRequestServer(Helpers.Uint16ToString(port),
		&Config.SingleRequestServer{
			SystemgeServerConfig: &Config.SystemgeServer{
				TcpSystemgeListenerConfig: &Config.TcpSystemgeListener{
					TcpServerConfig: &Config.TcpServer{
						TlsCertPath: "MyCertificate.crt",
						TlsKeyPath:  "MyKey.key",
						Port:        port,
					},
				},
				TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{},
			},
		},
		nil, nil,
		nil,
		SystemgeConnection.NewConcurrentMessageHandler(
			SystemgeConnection.AsyncMessageHandlers{
				"end": func(connection SystemgeConnection.SystemgeConnection, message *Message.Message) {
					if err := app.singleRequestServer.Stop(); err != nil {
						// shouldn't happen in this sample. Should be properly error handled in a real application though
						panic(Error.New("Error stopping singleRequestServer", err))
					}
					if dashboardClient != nil {
						if err := dashboardClient.Stop(); err != nil {
							// shouldn't happen in this sample. Should be properly error handled in a real application though
							panic(Error.New("Error stopping dashboardClient", err))
						}
					}
					stopFunc()
				},
			},
			SystemgeConnection.SyncMessageHandlers{
				"move": func(connection SystemgeConnection.SystemgeConnection, message *Message.Message) (string, error) {
					move, err := dto.UnmarshalMove(message.GetPayload())
					if err != nil {
						return "", Error.New("Error unmarshalling move", err)
					}
					chessMove, err := app.handleMove(move)
					if err != nil {
						return "", err
					}
					return Helpers.JsonMarshal(chessMove), nil
				},
				"getBoard": func(connection SystemgeConnection.SystemgeConnection, message *Message.Message) (string, error) {
					return app.marshalBoard(), nil
				},
			},
			nil, nil,
		),
	)
	if err := app.singleRequestServer.Start(); err != nil {
		// shouldn't happen in this sample. Should be properly error handled in a real application though
		panic(Error.New("Failed to start singleRequestServer", err))
	}

	dashboardClient = DashboardClientCustomService.New_(Helpers.Uint16ToString(port),
		&Config.DashboardClient{
			TcpSystemgeConnectionConfig: &Config.TcpSystemgeConnection{},
			TcpClientConfig: &Config.TcpClient{
				Address: "localhost:60000",
				TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
				Domain:  "example.com",
			},
		},
		app.singleRequestServer.Start,
		func() error {
			if err := app.singleRequestServer.Stop(); err != nil {
				// shouldn't happen in this sample. Should be properly error handled in a real application though
				panic(Error.New("Error stopping singleRequestServer", err))
			}
			return nil
		},
		app.singleRequestServer.GetStatus, app.singleRequestServer.GetMetrics,
		app.singleRequestServer.GetDefaultCommands(),
	)
	if err := dashboardClient.Start(); err != nil {
		// shouldn't happen in this sample. Should be properly error handled in a real application though
		panic(Error.New("Failed to start dashboardClient", err))
	}
	return app
}

func (app *AppChess) handleMove(move *dto.Move) (*dto.Move, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if app.isWhiteTurn() != move.White {
		return nil, Error.New("Not your turn", nil)
	}
	chessMove, err := app.move(move)
	if err != nil {
		return nil, Error.New("Invalid move", err)
	}
	return chessMove, nil
}
