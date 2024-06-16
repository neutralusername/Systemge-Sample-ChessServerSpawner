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
}

func New(client *Client.Client, args []string) (Application.Application, error) {
	app := &App{
		client: client,
	}
	return app, nil
}

func (app *App) OnStart() error {
	_, err := app.client.SyncMessage(topics.PROPAGATE_GAMESTART, app.client.GetName(), "...gameStartData...")
	if err != nil {
		err := app.client.AsyncMessage(topics.END, app.client.GetName(), app.client.GetName())
		if err != nil {
			app.client.GetLogger().Log(Utilities.NewError("Error sending async message", err).Error())
		}
	}
	return nil
}

func (app *App) OnStop() error {
	err := app.client.AsyncMessage(topics.PROPAGATE_GAMEEND, app.client.GetName(), "...gameEndData...")
	if err != nil {
		app.client.GetLogger().Log(Utilities.NewError("Error sending async message", err).Error())
	}
	return nil
}

func (app *App) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{
		app.client.GetName(): func(message *Message.Message) error {
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
