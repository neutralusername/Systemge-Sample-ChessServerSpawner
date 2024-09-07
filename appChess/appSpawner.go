package appChess

import (
	"sync"

	"github.com/neutralusername/Systemge/Config"
	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/SingleRequestServer"
	"github.com/neutralusername/Systemge/SystemgeConnection"
)

type AppSpawner struct {
	spawnedApps map[uint16]*AppChess
	mutex       *sync.Mutex

	portPool map[uint16]bool

	singleRequestServer *SingleRequestServer.Server
}

func NewAppSpawner() *AppSpawner {
	app := &AppSpawner{
		spawnedApps: make(map[uint16]*AppChess),
		mutex:       &sync.Mutex{},
		portPool: map[uint16]bool{
			60002: true,
			60003: true,
			60004: true,
			60005: true,
			60006: true,
			60007: true,
			60008: true,
			60009: true,
			60010: true,
			60011: true,
		},
	}
	app.singleRequestServer = SingleRequestServer.NewSingleRequestServer("chessSpawner",
		&Config.SingleRequestServer{
			SystemgeServerConfig: &Config.SystemgeServer{
				ListenerConfig: &Config.TcpSystemgeListener{
					TcpServerConfig: &Config.TcpServer{
						TlsCertPath: "MyCertificate.crt",
						TlsKeyPath:  "MyKey.key",
						Port:        60001,
					},
				},
				ConnectionConfig: &Config.TcpSystemgeConnection{},
			},
			DashboardClientConfig: &Config.DashboardClient{
				ConnectionConfig: &Config.TcpSystemgeConnection{},
				ClientConfig: &Config.TcpClient{
					Address: "localhost:60000",
					TlsCert: Helpers.GetFileContent("MyCertificate.crt"),
					Domain:  "example.com",
				},
			},
		},
		nil, SystemgeConnection.NewConcurrentMessageHandler(
			SystemgeConnection.AsyncMessageHandlers{},
			SystemgeConnection.SyncMessageHandlers{
				"spawn": func(connection SystemgeConnection.SystemgeConnection, message *Message.Message) (string, error) {
					app.mutex.Lock()
					defer app.mutex.Unlock()

					freePort := uint16(0)
					for port, free := range app.portPool {
						if free {
							freePort = port
							app.portPool[port] = false
							break
						}
					}
					if freePort == 0 {
						return "", Error.New("no free ports", nil)
					}

					app.spawnedApps[freePort] = NewAppChess(freePort, func() {
						app.mutex.Lock()
						defer app.mutex.Unlock()

						app.portPool[freePort] = true
						delete(app.spawnedApps, freePort)
					})
					return Helpers.Uint16ToString(freePort), nil
				},
			},
			nil, nil,
		),
	)
	if err := app.singleRequestServer.Start(); err != nil {
		// shouldn't happen in this sample. Should be properly error handled in a real application though
		panic(Error.New("Failed to start singleRequestServer", err))
	}
	return app
}
