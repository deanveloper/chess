package chess

// Move represents a move in a chess game.
type Move struct {
	Snapshot Game

	Moving Piece
	To     Space

	Promotion PieceType
}
