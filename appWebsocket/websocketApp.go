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
			app.client.GetWebsocketServer().Groupcast(message.GetOrigin(), message)
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
			app.mutex.Lock()
			defer app.mutex.Unlock()
			if app.clientGameIds[connection.GetId()] != "" {
				return Utilities.NewError("You are already in a game", nil)
			}
			if app.clientGameIds[message.GetPayload()] != "" {
				return Utilities.NewError("Opponent is already in a game", nil)
			}
			gameId := app.client.GetName() + " " + message.GetPayload()
			_, err := app.client.SyncMessage(topics.NEW, app.client.GetName(), gameId)
			if err != nil {
				return Utilities.NewError("Error spawning new game client", err)
			}
			err = app.client.GetWebsocketServer().AddToGroup(gameId, connection.GetId())
			if err != nil {
				_, err := app.client.SyncMessage(topics.END, app.client.GetName(), gameId)
				if err != nil {
					app.client.GetLogger().Log(Utilities.NewError("Error sending end message for game: "+gameId, err).Error())
				}
				return Utilities.NewError("Error adding \""+connection.GetId()+"\" to group \""+gameId+"\"", err)
			}
			err = app.client.GetWebsocketServer().AddToGroup(gameId, message.GetPayload())
			if err != nil {
				err := app.client.GetWebsocketServer().RemoveFromGroup(gameId, connection.GetId())
				if err != nil {
					app.client.GetLogger().Log(Utilities.NewError("Error removing \""+connection.GetId()+"\" from group \""+gameId+"\"", err).Error())
				}
				_, err = app.client.SyncMessage(topics.END, connection.GetId(), gameId)
				if err != nil {
					app.client.GetLogger().Log(Utilities.NewError("Error sending end message for game: "+gameId, err).Error())
				}
				return Utilities.NewError("Error adding \""+message.GetPayload()+"\" to group \""+gameId+"\"", err)
			}
			app.clientGameIds[connection.GetId()] = gameId
			app.clientGameIds[message.GetPayload()] = gameId
			return nil
		},
		"move": func(connection *WebsocketClient.Client, message *Message.Message) error {
			app.mutex.Lock()
			defer app.mutex.Unlock()
			gameId := app.clientGameIds[connection.GetId()]
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
			defer app.mutex.Unlock()
			gameId := app.clientGameIds[connection.GetId()]
			if gameId == "" {
				return Utilities.NewError("You are not in a game", nil)
			}
			_, err := app.client.SyncMessage(topics.END, connection.GetId(), gameId)
			if err != nil {
				return Utilities.NewError("Error sending end message", err)
			}
			ids := strings.Split(gameId, " ")
			delete(app.clientGameIds, ids[0])
			delete(app.clientGameIds, ids[1])
			app.client.GetWebsocketServer().Groupcast(gameId, Message.NewAsync("endGame", app.client.GetName(), ""))
			err = app.client.GetWebsocketServer().RemoveFromGroup(gameId, ids[0])
			if err != nil {
				app.client.GetLogger().Log(Utilities.NewError("Error removing \""+ids[0]+"\" from group \""+gameId+"\"", err).Error())
			}
			err = app.client.GetWebsocketServer().RemoveFromGroup(gameId, ids[1])
			if err != nil {
				app.client.GetLogger().Log(Utilities.NewError("Error removing \""+ids[1]+"\" from group \""+gameId+"\"", err).Error())
			}
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
	defer app.mutex.Unlock()
	gameId := app.clientGameIds[connection.GetId()]
	if gameId == "" {
		return
	}
	_, err := app.client.SyncMessage(topics.END, app.client.GetName(), gameId)
	if err != nil {
		app.client.GetLogger().Log(Utilities.NewError("Error sending end message for game: "+gameId, err).Error())
	}
	ids := strings.Split(gameId, " ")
	delete(app.clientGameIds, ids[0])
	delete(app.clientGameIds, ids[1])
	app.client.GetWebsocketServer().Groupcast(gameId, Message.NewAsync("endGame", app.client.GetName(), ""))
	err = app.client.GetWebsocketServer().RemoveFromGroup(gameId, ids[0])
	if err != nil {
		app.client.GetLogger().Log(Utilities.NewError("Error removing \""+ids[0]+"\" from group \""+gameId+"\"", err).Error())
	}
	err = app.client.GetWebsocketServer().RemoveFromGroup(gameId, ids[1])
	if err != nil {
		app.client.GetLogger().Log(Utilities.NewError("Error removing \""+ids[1]+"\" from group \""+gameId+"\"", err).Error())
	}
}
