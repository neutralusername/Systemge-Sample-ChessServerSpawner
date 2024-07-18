package appChess

import "Systemge/Helpers"

func (app *App) generateAlgebraicNotation(fromRow, fromCol, toRow, toCol int) string {
	notation := ""
	piece := app.board[fromRow][fromCol]
	switch piece.(type) {
	case *King:
		if fromCol-toCol == 2 {
			notation = "O-O-O"
		} else if fromCol-toCol == -2 {
			notation = "O-O"
		} else {
			notation = "K"
		}
		return notation
	case *Pawn:
		if fromCol != toCol && app.board[toRow][toCol] == nil {
			notation = app.getColumnLetter(fromCol) + "x" + app.getColumnLetter(toCol) + app.getRowNumber(toRow)
		} else {
			notation = app.getColumnLetter(toCol) + app.getRowNumber(toRow)
		}
		if toRow == 0 || toRow == 7 {
			notation += "=Q"
		}
		return notation
	}
	notation = piece.getLetter()
	if app.board[toRow][toCol] != nil {
		notation += "x"
	}
	notation += app.getColumnLetter(toCol) + app.getRowNumber(toRow)
	return notation
}

func (app *App) getColumnLetter(col int) string {
	switch col {
	case 0:
		return "a"
	case 1:
		return "b"
	case 2:
		return "c"
	case 3:
		return "d"
	case 4:
		return "e"
	case 5:
		return "f"
	case 6:
		return "g"
	case 7:
		return "h"
	}
	return ""
}

func (app *App) getRowNumber(row int) string {
	return Helpers.IntToString(8 - row)
}
