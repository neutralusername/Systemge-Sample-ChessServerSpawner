package appWebsocket

import (
	"Systemge/Application"
	"Systemge/Client"
	"Systemge/Message"
	"Systemge/Utilities"
	"Systemge/WebsocketClient"
	"SystemgeSampleChessServer/topics"
	"strings"
	"sync"
)

type WebsocketApp struct {
	client *Client.Client

	clientGameIds map[string]string
	mutex         sync.Mutex
}

func New(messageBrokerClient *Client.Client, args []string) (Application.WebsocketApplication, error) {
	return &WebsocketApp{
		client: messageBrokerClient,

		clientGameIds: make(map[string]string),
		mutex:         sync.Mutex{},
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
	return map[string]Application.WebsocketMessageHandler{
		"startGame": func(connection *WebsocketClient.Client, message *Message.Message) error {
			gameId, err := func() (string, error) {
				app.mutex.Lock()
				defer app.mutex.Unlock()
				if app.clientGameIds[connection.GetId()] != "" {
					return "", Utilities.NewError("You are already in a game", nil)
				}
				if app.clientGameIds[message.GetPayload()] != "" {
					return "", Utilities.NewError("Opponent is already in a game", nil)
				}
				gameId := app.client.GetName() + " " + message.GetPayload()
				app.clientGameIds[connection.GetId()] = gameId
				app.clientGameIds[message.GetPayload()] = gameId
				return gameId, nil
			}()
			if err != nil || !app.client.GetWebsocketServer().ClientExists(message.GetPayload()) {
				app.mutex.Lock()
				delete(app.clientGameIds, connection.GetId())
				delete(app.clientGameIds, message.GetPayload())
				app.mutex.Unlock()
				return Utilities.NewError("Error starting game", err)
			}
			_, err = app.client.SyncMessage(topics.NEW, connection.GetId(), gameId)
			if err != nil {
				func() {
					app.mutex.Lock()
					defer app.mutex.Unlock()
					delete(app.clientGameIds, connection.GetId())
					delete(app.clientGameIds, message.GetPayload())
				}()
				return Utilities.NewError("Error sending new message", err)
			}
			return nil
		},
		"move": func(connection *WebsocketClient.Client, message *Message.Message) error {
			app.mutex.Lock()
			gameId := app.clientGameIds[connection.GetId()]
			app.mutex.Unlock()
			if gameId == "" {
				return Utilities.NewError("You are not in a game", nil)
			}
			err := app.client.AsyncMessage(gameId, connection.GetId(), message.GetPayload())
			if err != nil {
				return Utilities.NewError("Error sending move message", err)
			}
			return nil
		},
		"endGame": func(connection *WebsocketClient.Client, message *Message.Message) error {
			app.mutex.Lock()
			gameId := app.clientGameIds[connection.GetId()]
			app.mutex.Unlock()
			if gameId == "" {
				return Utilities.NewError("You are not in a game", nil)
			}
			_, err := app.client.SyncMessage(topics.END, connection.GetId(), gameId)
			if err != nil {
				return Utilities.NewError("Error sending end message", err)
			}
			app.mutex.Lock()
			ids := strings.Split(gameId, " ")
			delete(app.clientGameIds, ids[0])
			delete(app.clientGameIds, ids[1])
			app.mutex.Unlock()
			return nil
		},
	}
}

func (app *WebsocketApp) OnConnectHandler(connection *WebsocketClient.Client) {
	err := connection.Send(Message.NewAsync("connected", app.client.GetName(), connection.GetId()).Serialize())
	if err != nil {
		connection.Disconnect()
		app.client.GetLogger().Log(Utilities.NewError("Error sending connected message", err).Error())
	}
}

func (app *WebsocketApp) OnDisconnectHandler(connection *WebsocketClient.Client) {
	app.mutex.Lock()
	gameId := app.clientGameIds[connection.GetId()]
	app.mutex.Unlock()
	if gameId == "" {
		return
	}
	_, err := app.client.SyncMessage(topics.END, app.client.GetName(), gameId)
	if err != nil {
		app.client.GetLogger().Log(Utilities.NewError("Error sending end message for game: "+gameId, err).Error())
	}
	app.mutex.Lock()
	ids := strings.Split(gameId, " ")
	delete(app.clientGameIds, ids[0])
	delete(app.clientGameIds, ids[1])
	app.mutex.Unlock()
}
