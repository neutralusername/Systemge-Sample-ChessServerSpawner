package appChess

import (
	"SystemgeSampleChessServer/dto"
	"SystemgeSampleChessServer/topics"
	"strings"

	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Node"
)

func (app *App) OnStart(node *Node.Node) error {
	gameId := node.GetName()
	ids := strings.Split(gameId, "-")
	app.whiteId = ids[0]
	app.blackId = ids[1]
	err := node.AsyncMessage(topics.PROPAGATE_GAMESTART, Helpers.JsonMarshal(&dto.GameStart{
		Board:             app.marshalBoard(),
		TcpEndpointConfig: node.GetSystemgeEndpointConfig(),
	}))
	if err != nil {
		panic(Error.New("Error sending async message", err))
	}
	return nil
}

func (app *App) OnStop(node *Node.Node) error {
	err := node.AsyncMessage(topics.PROPAGATE_GAMEEND, Helpers.JsonMarshal(node.GetSystemgeEndpointConfig()))
	if err != nil {
		panic(Error.New("Error sending async message", err))
	}
	return nil
}
