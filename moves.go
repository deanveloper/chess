package chess

// Move represents a move in a chess game.
type Move struct {
	Moving Piece
	To     Space

	Promotion PieceType
}
