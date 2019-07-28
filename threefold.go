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
	entry.canPassant = g.enPassant.Rank != 0

	// check for castle
	if player == Black {
		entry.canKingCastle = g.castles.BlackKing
		entry.canQueenCastle = g.castles.BlackQueen
	} else {
		entry.canKingCastle = g.castles.WhiteKing
		entry.canQueenCastle = g.castles.WhiteQueen
	}
}
