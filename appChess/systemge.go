package appChess

import (
	"Systemge/Config"
	"Systemge/Error"
	"Systemge/Helpers"
	"Systemge/Message"
	"Systemge/Node"
	"Systemge/Tcp"
	"strings"
)

func (app *App) GetSystemgeComponentConfig() Config.Systemge {
	return Config.Systemge{
		HandleMessagesSequentially: false,

		BrokerSubscribeDelayMs:    1000,
		TopicResolutionLifetimeMs: 10000,
		SyncResponseTimeoutMs:     10000,
		TcpTimeoutMs:              5000,

		ResolverEndpoint: Tcp.NewEndpoint("127.0.0.1:60000", "example.com", Helpers.GetFileContent("MyCertificate.crt")),
	}
}

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{}
}

func (app *App) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		app.gameId: func(node *Node.Node, message *Message.Message) (string, error) {
			segments := strings.Split(message.GetPayload(), " ")
			if len(segments) != 4 {
				return "", Error.New("Invalid message format", nil)
			}
			chessMove, err := app.handleMoveRequest(message.GetOrigin(), Helpers.StringToInt(segments[0]), Helpers.StringToInt(segments[1]), Helpers.StringToInt(segments[2]), Helpers.StringToInt(segments[3]))
			if err != nil {
				return "", err
			}
			return chessMove.Marshal(), nil
		},
	}
}

func (app *App) handleMoveRequest(playerId string, rowFrom, colFrom, rowTo, colTo int) (*ChessMove, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if app.isWhiteTurn() && playerId != app.whiteId {
		return nil, Error.New("Not your turn", nil)
	}
	if !app.isWhiteTurn() && playerId != app.blackId {
		return nil, Error.New("Not your turn", nil)
	}
	chessMove, err := app.move(rowFrom, colFrom, rowTo, colTo)
	if err != nil {
		return nil, Error.New("Invalid move", err)
	}
	return chessMove, nil
}
