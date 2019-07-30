package chess

import "fmt"

// Move represents a move in a chess game.
type Move struct {
	Snapshot Game

	Moving Piece
	To     Space

	Promotion PieceType
}

func (m Move) String() string {
	return fmt.Sprintf("Move{%s to %s}", m.Moving, m.To)
}
