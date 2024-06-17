package appChess

import "strings"

func (app *App) marshalBoard() string {
	var builder strings.Builder
	for _, row := range app.board {
		for _, piece := range row {
			if piece == nil {
				builder.WriteString(".")
			} else {
				builder.WriteString(piece.getLetter())
			}
		}
	}
	return builder.String()
}

func (app *App) isWhiteTurn() bool {
	return len(app.moves)%2 == 0
}
