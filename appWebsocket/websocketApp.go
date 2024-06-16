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
		topics.PROPAGATE_GAMESTART: func(message *Message.Message) error {
			gameId := message.GetOrigin()
			ids := strings.Split(gameId, "-")
			err := app.client.GetWebsocketServer().AddToGroup(gameId, ids[0])
			if err != nil {
				return Utilities.NewError("Error adding \""+ids[0]+"\" to group \""+gameId+"\"", err)
			}
			err = app.client.GetWebsocketServer().AddToGroup(gameId, ids[1])
			if err != nil {
				err := app.client.GetWebsocketServer().RemoveFromGroup(gameId, ids[0])
				if err != nil {
					app.client.GetLogger().Log(Utilities.NewError("Error removing \""+ids[0]+"\" from group \""+gameId+"\"", err).Error())
				}
				return Utilities.NewError("Error adding \""+ids[1]+"\" to group \""+gameId+"\"", err)
			}
			app.mutex.Lock()
			app.clientGameIds[ids[0]] = gameId
			app.clientGameIds[ids[1]] = gameId
			app.mutex.Unlock()
			app.client.GetWebsocketServer().Groupcast(gameId, message)
			return nil
		},
		topics.PROPAGATE_GAMEEND: func(message *Message.Message) error {
			gameId := message.GetPayload()
			ids := strings.Split(gameId, "-")
			app.mutex.Lock()
			delete(app.clientGameIds, ids[0])
			delete(app.clientGameIds, ids[1])
			app.mutex.Unlock()
			app.client.GetWebsocketServer().Groupcast(message.GetOrigin(), message)
			err := app.client.GetWebsocketServer().RemoveFromGroup(gameId, ids[0])
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

func (app *WebsocketApp) GetSyncMessageHandlers() map[string]Application.SyncMessageHandler {
	return map[string]Application.SyncMessageHandler{}
}

func (app *WebsocketApp) GetCustomCommandHandlers() map[string]Application.CustomCommandHandler {
	return map[string]Application.CustomCommandHandler{}
}

func (app *WebsocketApp) GetWebsocketMessageHandlers() map[string]Application.WebsocketMessageHandler {
	return map[string]Application.WebsocketMessageHandler{
		"startGame": func(client *WebsocketClient.Client, message *Message.Message) error {
			whiteId := client.GetId()
			blackId := message.GetPayload()
			app.mutex.Lock()
			whiteGameId := app.clientGameIds[whiteId]
			blackGameId := app.clientGameIds[blackId]
			app.mutex.Unlock()
			if blackId == whiteId {
				return Utilities.NewError("You cannot play against yourself", nil)
			}
			if whiteGameId != "" {
				return Utilities.NewError("You are already in a game", nil)
			}
			if blackGameId != "" {
				return Utilities.NewError("Opponent is already in a game", nil)
			}
			_, err := app.client.SyncMessage(topics.NEW, app.client.GetName(), whiteId+"-"+blackId)
			if err != nil {
				return Utilities.NewError("Error spawning new game client", err)
			}
			return nil
		},
		"move": func(client *WebsocketClient.Client, message *Message.Message) error {
			app.mutex.Lock()
			defer app.mutex.Unlock()
			gameId := app.clientGameIds[client.GetId()]
			if gameId == "" {
				return Utilities.NewError("You are not in a game", nil)
			}
			err := app.client.AsyncMessage(gameId, client.GetId(), message.GetPayload())
			if err != nil {
				return Utilities.NewError("Error sending move message", err)
			}
			return nil
		},
		"endGame": func(client *WebsocketClient.Client, message *Message.Message) error {
			app.mutex.Lock()
			gameId := app.clientGameIds[client.GetId()]
			app.mutex.Unlock()
			if gameId == "" {
				return Utilities.NewError("You are not in a game", nil)
			}
			_, err := app.client.SyncMessage(topics.END, app.client.GetName(), gameId)
			if err != nil {
				app.client.GetLogger().Log(Utilities.NewError("Error sending end message for game: "+gameId, err).Error())
			}
			return nil
		},
	}
}

func (app *WebsocketApp) OnConnectHandler(client *WebsocketClient.Client) {
	err := client.Send(Message.NewAsync("connected", app.client.GetName(), client.GetId()).Serialize())
	if err != nil {
		client.Disconnect()
		app.client.GetLogger().Log(Utilities.NewError("Error sending connected message", err).Error())
	}
}

func (app *WebsocketApp) OnDisconnectHandler(client *WebsocketClient.Client) {
	app.mutex.Lock()
	gameId := app.clientGameIds[client.GetId()]
	app.mutex.Unlock()
	if gameId == "" {
		return
	}
	_, err := app.client.SyncMessage(topics.END, app.client.GetName(), gameId)
	if err != nil {
		app.client.GetLogger().Log(Utilities.NewError("Error sending end message for game: "+gameId, err).Error())
	}
}
