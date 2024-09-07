package dto

import "encoding/json"

type Move struct {
	AlgebraicNotation string `json:"algebraicNotation"`
	FromRow           int    `json:"fromRow"`
	FromCol           int    `json:"fromCol"`
	ToRow             int    `json:"toRow"`
	ToCol             int    `json:"toCol"`
	White             bool   `json:"white"`
}

func UnmarshalMove(str string) (*Move, error) {
	move := &Move{}
	err := json.Unmarshal([]byte(str), move)
	return move, err
}
