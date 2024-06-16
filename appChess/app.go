package appChess

import (
	"Systemge/Application"
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"SystemgeSampleChessServer/topics"
)

type App struct {
	client *Client.Client

	id string
}

func New(client *Client.Client, args []string) (Application.Application, error) {
	if len(args) != 1 {
		return nil, Utilities.NewError("Expected 1 argument", nil)
	}
	app := &App{
		client: client,

		id: args[0],
	}
	return app, nil
}

func (app *App) OnStart() error {
	return nil
}

func (app *App) OnStop() error {
	return nil
}

func (app *App) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{
		app.id: func(message *Message.Message) error {
			println(app.client.GetName() + " received \"" + message.GetPayload() + "\" from: " + message.GetOrigin())
			err := app.client.AsyncMessage(topics.PROPAGATE_MOVE, app.client.GetName(), message.GetPayload())
			if err != nil {
				panic(err)
			}
			return nil
		},
	}
}

func (app *App) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{}
}

func (app *App) GetCustomCommandHandlers() map[string]Application.CustomCommandHandler {
	return map[string]Application.CustomCommandHandler{}
}
