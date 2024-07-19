module SystemgeSampleChessServer

go 1.22.3

replace Systemge => ../Systemge

require (
	Systemge v0.0.0-00010101000000-000000000000
	github.com/gorilla/websocket v1.5.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/oauth2 v0.21.0 // indirect
)
