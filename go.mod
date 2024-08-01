module SystemgeSampleChessServer

go 1.22.3

//replace github.com/neutralusername/Systemge => ../Systemge
require (
	github.com/gorilla/websocket v1.5.3
	github.com/neutralusername/Systemge v0.0.0-20240801073437-6eb2bbc048a5
)

require golang.org/x/oauth2 v0.21.0 // indirect
