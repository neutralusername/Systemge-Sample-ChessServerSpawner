package dto

import (
	"encoding/json"

	"github.com/neutralusername/Systemge/Config"
)

type GameStart struct {
	Board           string            `json:"board"`
	TcpClientConfig *Config.TcpClient `json:"tcpClientConfig"`
}

func UnmarshalGameStart(str string) *GameStart {
	gameStart := &GameStart{}
	json.Unmarshal([]byte(str), gameStart)
	return gameStart
}
