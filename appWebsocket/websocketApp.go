package appWebsocket

import (
	"Systemge/Application"
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"Systemge/WebsocketClient"
	"SystemgeSampleChessServer/topics"
)

type WebsocketApp struct {
	client *Client.Client
}

func New(messageBrokerClient *Client.Client, args []string) (Application.WebsocketApplication, error) {
	return &WebsocketApp{
		client: messageBrokerClient,
	}, nil
}

func (app *WebsocketApp) OnStart() error {
	return nil
}

func (app *WebsocketApp) OnStop() error {
	return nil
}

func (app *WebsocketApp) GetAsyncMessageHandlers() map[string]Application.AsyncMessageHandler {
	return map[string]Application.AsyncMessageHandler{
		topics.PROPAGATE_MOVE: func(message *Message.Message) error {
			println(app.client.GetName() + " received \"" + message.GetPayload() + "\" from: " + message.GetOrigin())
			return nil
		},
	}
}

func (app *WebsocketApp) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{}
}

func (app *WebsocketApp) GetCustomCommandHandlers() map[string]Application.CustomCommandHandler {
	return map[string]Application.CustomCommandHandler{}
}

func (app *WebsocketApp) GetWebsocketMessageHandlers() map[string]Application.WebsocketMessageHandler {
	return map[string]Application.WebsocketMessageHandler{}
}

func (app *WebsocketApp) OnConnectHandler(connection *WebsocketClient.Client) {
	_, err := app.client.SyncMessage(topics.NEW, connection.GetId(), connection.GetId())
	if err != nil {
		panic(Utilities.NewError("Error sending sync message", err))
	}
	err = app.client.AsyncMessage(connection.GetId(), connection.GetId(), "e4e2")
	if err != nil {
		panic(Utilities.NewError("Error sending async message", err))
	}
}

func (app *WebsocketApp) OnDisconnectHandler(connection *WebsocketClient.Client) {
	_, err := app.client.SyncMessage(topics.END, app.client.GetName(), connection.GetId())
	if err != nil {
		panic(err)
	}
}
