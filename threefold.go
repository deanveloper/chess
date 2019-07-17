package chess

type threeFoldEntry struct {
	pieces         [32]Piece
	canPassant     bool
	canQueenCastle bool
	canKingCastle  bool
}

type threeFoldTracker map[threeFoldEntry]int

// creates a threeFoldEntry for the current state of g
func (t threeFoldTracker) addToTracker(g *Game) {
	var entry threeFoldEntry
	copy(entry.pieces[0:16], g.white[:])
	copy(entry.pieces[16:32], g.black[:])

	if len(g.History) == 0 {
		return
	}

	lastMove := g.History[len(g.History)-1]
	player := lastMove.Moving.Color

	// check for en passant
	for _, pawn := range g.TypedAlivePieces(player.Other(), Pawn) {
		for _, space := range pawn.Seeing(g) {
			if space.File-pawn.Location.File == 0 {
				continue
			}
			if _, ok := g.PieceAt(space); !ok {
				entry.canPassant = true
				break
			}
		}
	}

	// check for castle
	king := g.TypedAlivePieces(player.Other(), King)[0]
	for _, space := range king.Seeing(g) {
		fileDiff := space.File - king.Location.File
		if fileDiff == 1 || fileDiff == 0 {
			continue
		}
		if _, ok := g.PieceAt(space); !ok {
			entry.canPassant = true
			break
		}
	}
}
