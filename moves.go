package chess

// Move represents a move in a chess game.
type Move struct {
	Snapshot Game

	Moving Piece
	To     Space

	Promotion PieceType
}

func (m Move) String() string {
	return "Move{" + m.Moving.String() + " to " + m.To.String()
}
