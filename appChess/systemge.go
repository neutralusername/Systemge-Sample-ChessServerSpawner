package appChess

import (
	"SystemgeSampleChessServer/dto"

	"github.com/neutralusername/Systemge/Error"
	"github.com/neutralusername/Systemge/Helpers"
	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{}
}

func (app *App) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{}
}

func (app *App) moveMessageHandler(node *Node.Node, message *Message.Message) (string, error) {
	move, err := dto.UnmarshalMove(message.GetPayload())
	if err != nil {
		return "", Error.New("Error unmarshalling move", err)
	}
	chessMove, err := app.handleMove(move)
	if err != nil {
		return "", err
	}
	return Helpers.JsonMarshal(chessMove), nil

}

func (app *App) handleMove(move *dto.Move) (*dto.Move, error) {
	app.mutex.Lock()
	defer app.mutex.Unlock()
	if app.isWhiteTurn() && move.PlayerId != app.whiteId {
		return nil, Error.New("Not your turn", nil)
	}
	if !app.isWhiteTurn() && move.PlayerId != app.blackId {
		return nil, Error.New("Not your turn", nil)
	}
	chessMove, err := app.move(move)
	if err != nil {
		return nil, Error.New("Invalid move", err)
	}
	return chessMove, nil
}
