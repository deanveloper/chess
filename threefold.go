package chess

type threeFoldEntry struct {
	board          [8][8]Piece
	canPassant     bool
	canQueenCastle bool
	canKingCastle  bool
}

type threeFoldTracker map[threeFoldEntry]int

// creates a threeFoldEntry for the current state of g
func (t threeFoldTracker) addToTracker(g *Game) {
	var entry threeFoldEntry
	entry.board = g.board

	if len(g.History) == 0 {
		return
	}

	lastMove := g.History[len(g.History)-1]
	player := lastMove.Moving.Color

	// check for en passant
	entry.canPassant = g.EnPassant.Rank != 0

	// check for castle
	if player == Black {
		entry.canKingCastle = g.Castles.BlackKing
		entry.canQueenCastle = g.Castles.BlackQueen
	} else {
		entry.canKingCastle = g.Castles.WhiteKing
		entry.canQueenCastle = g.Castles.WhiteQueen
	}
}
