package chess

// Game represents a game of chess
type Game struct {
	white [16]Piece
	black [16]Piece

	fiftyMoveDetector int
	threeFoldDetector [][]Space

	History []Move
}

// Clone returns a new instance of `g` with the same history.
func (g *Game) Clone() *Game {
	var newG = &Game{
		white:             g.white,
		black:             g.black,
		fiftyMoveDetector: g.fiftyMoveDetector,
		threeFoldDetector: make([][]Space, len(g.threeFoldDetector), len(g.threeFoldDetector)+1),
		History:           make([]Move, len(g.History), len(g.History)+1),
	}
	copy(newG.History, g.History)
	copy(newG.threeFoldDetector, g.threeFoldDetector)
	return newG
}

// returns a mutable slice of the pieces for a given color.
// not exported because modifying this slice will modify the game.
func (g *Game) pieces(c Color) []Piece {
	var color []Piece
	if c == Black {
		color = g.black[:]
	} else {
		color = g.white[:]
	}

	return color
}

// AlivePieces returns the pieces that c has on the board. Modifying this slice
// has no impact on the game.
func (g *Game) AlivePieces(c Color) []Piece {

	var color [16]Piece
	if c == Black {
		color = g.black
	} else {
		color = g.white
	}
	pieces := make([]Piece, 0, len(color))

	for _, p := range color {
		if p.Location != Taken {
			pieces = append(pieces, p)
		}
	}

	return pieces
}

// TakenPieces returns the pieces that c no longer has on the board. Modifying this slice
// has no impact on the game.
func (g *Game) TakenPieces(c Color) []Piece {

	var color [16]Piece
	if c == Black {
		color = g.black
	} else {
		color = g.white
	}

	taken := make([]Piece, 0, len(color))

	for _, p := range color {
		if p.Location == Taken {
			taken = append(taken, p)
		}
	}

	return taken
}

// PieceAt returns the piece at a given space, and an `ok`
// boolean on if there was a piece on that space at all.
func (g *Game) PieceAt(s Space) (Piece, bool) {
	for _, p := range g.white {
		if p.Location == s {
			return p, true
		}
	}
	for _, p := range g.black {
		if p.Location == s {
			return p, true
		}
	}
	return Piece{}, false
}

// InCheck returns if `c` is in check.
func (g *Game) InCheck(c Color) bool {
	pieces := g.AlivePieces(c)

	var king Piece
	for _, piece := range pieces {
		if piece.Type == King {
			king = piece
			break
		}
	}

	oPieces := g.AlivePieces(c.Other())
	for _, piece := range oPieces {
		seeing := piece.Seeing(g)
		for _, sees := range seeing {
			if sees == king.Location {
				return true
			}
		}
	}

	return false
}

// InCheckmate returns if `c` is in checkmate.
func (g *Game) InCheckmate(c Color) bool {
	return g.InCheck(c) && !g.canMove(c)
}

// InStalemate returns if `c` is in stalemate.
func (g *Game) InStalemate(c Color) bool {
	return !g.InCheck(c) && !g.canMove(c)
}

// CanDraw returns if the game is now able to be drawn via threefold
// repetition, the 50 move rule, or insufficient material.
// (to draw via threefold repetition, the last position played
// must have been played at least 2 other times).
func (g *Game) CanDraw() bool {
	if len(g.History) < 5 {
		return false
	}

	color := g.History[len(g.History)-1].Moving.Color

	// threefold repetition
	if len(g.threeFoldDetector) >= 3 {
		current := g.threeFoldDetector[len(g.threeFoldDetector)-1]

		var fold int
		for i := len(g.threeFoldDetector) - 1; i >= 0; i -= 2 {
			if slicesEqual(current, g.threeFoldDetector[i]) {
				fold++
				if fold == 3 {
					return true
				}
			}
		}
	}

	// 50 move rule
	if g.fiftyMoveDetector >= 50 && !g.InCheckmate(color.Other()) {
		return true
	}

	// insufficient material... oh no
	whiteAlive := g.AlivePieces(White)
	blackAlive := g.AlivePieces(Black)

	// king vs king
	if len(whiteAlive) == 1 && len(blackAlive) == 1 {
		return true
	}

	// king/(bishop|knight) vs king
	if len(whiteAlive) == 2 && len(blackAlive) == 1 {
		switch {
		case whiteAlive[0].Type == Bishop,
			whiteAlive[0].Type == Knight,
			whiteAlive[1].Type == Bishop,
			whiteAlive[1].Type == Knight:

			return true
		}
	}
	if len(whiteAlive) == 1 && len(blackAlive) == 2 {
		switch {
		case blackAlive[0].Type == Bishop,
			blackAlive[0].Type == Knight,
			blackAlive[1].Type == Bishop,
			blackAlive[1].Type == Knight:

			return true
		}
	}

	// if all that exist are bishops on the same color...
	var bishops int
	var bishopColor Color

	// white bishops
	for _, piece := range whiteAlive {
		if piece.Type == Bishop {
			if bishops == 0 {
				bishopColor = piece.Location.Color()
			} else {
				if bishopColor != piece.Location.Color() {
					break
				}
			}
			bishops++
		}
	}
	if bishops != len(whiteAlive)-1 {
		return false
	}

	bishops = 0
	// black bishops
	for _, piece := range whiteAlive {
		if piece.Type == Bishop {
			if bishops == 0 {
				bishopColor = piece.Location.Color()
			} else {
				if bishopColor != piece.Location.Color() {
					break
				}
			}
			bishops++
		}
	}
	if bishops != len(blackAlive)-1 {
		return false
	}

	return false
}

// InThreefold returns if the game is able to be drawn via
// threefold repitition (if the last position played has appeared 3 times).
func (g *Game) InThreefold() bool {
	if len(g.threeFoldDetector) < 3 {
		return false
	}

	current := g.threeFoldDetector[len(g.threeFoldDetector)-1]

	var fold int
	for i := len(g.threeFoldDetector) - 1; i >= 0; i -= 2 {
		if slicesEqual(current, g.threeFoldDetector[i]) {
			fold++
			if fold == 3 {
				return true
			}
		}
	}
	return false
}

// PeekMove looks at a move and determines if the move is valid or not.
// Returns a new game instance in which the move is made,
// or an error if the move is not valid.
//
// It is safe to set the resulting game pointer to your instance of the game.
func (g *Game) PeekMove(m Move) (*Game, error) {
	newG := g.Clone()

	_, err := newG.MakeMove(m)
	if err != nil {
		return nil, err
	}

	return newG, nil
}

// MakeMove makes a move in the game, or returns an error if the move is not possible.
func (g *Game) MakeMove(m Move) (string, error) {

	player := m.Moving.Color

	// check if the move is valid
	var validMove bool
	seeing := m.Moving.Seeing(g)
	for _, s := range seeing {
		if m.To == s {
			validMove = true
			break
		}
	}
	if !validMove {
		return "", &MoveError{
			Cause:   m,
			InCheck: false,
		}
	}

	// TODO - check to see if the player is in check
	// TODO - implement castling properly

	// variable used later for draw-status-updating
	var pieceTaken bool

	// make the move
	pieces := g.pieces(player)
	otherPieces := g.pieces(player.Other())
	for i, piece := range pieces {
		if piece.Location == m.Moving.Location {
			pieces[i].Location = m.To
		}
	}
	for i, piece := range otherPieces {
		if piece.Location == m.To {
			pieces[i].Location = Taken
			pieceTaken = true
		}
	}
	g.History = append(g.History, m)

	// update draw detectors
	if pieceTaken || m.Moving.Type == Pawn {
		g.fiftyMoveDetector = 0
		g.threeFoldDetector = nil
	} else {
		g.fiftyMoveDetector++
		g.addToThreeFoldDetector()
	}

	algebraicString := m.AlgebraicString()
	switch {
	case g.InCheckmate(player.Other()):
		algebraicString += "#"
	case g.InCheck(player.Other()):
		algebraicString += "+"
	case g.InStalemate(player.Other()):
		break
	}

	return algebraicString, nil
}

func (g *Game) addToThreeFoldDetector() {
	var spaces []Space
	for _, piece := range g.white {
		spaces = append(spaces, piece.Seeing(g)...)
	}
	for _, piece := range g.black {
		spaces = append(spaces, piece.Seeing(g)...)
	}
	g.threeFoldDetector = append(g.threeFoldDetector, spaces)
}

func (g *Game) canMove(c Color) bool {
	newG := g.Clone()

	// func canUncheck(p Piece) bool
	canUncheck := func(p Piece) bool {
		for _, seen := range p.Seeing(newG) {
			_, err := newG.MakeMove(Move{
				g:      newG,
				Moving: p,
				To:     seen,
			})
			if err != nil {
				continue
			}
			if !newG.InCheck(c) {
				return false
			}
		}
		return false
	}
	// end local function

	var king Piece
	for _, piece := range newG.pieces(c) {
		if piece.Type == King {
			king = piece
			break
		}
	}

	// look at the king's moves first
	if canUncheck(king) {
		return true
	}

	// check the rest of the pieces
	for _, piece := range newG.AlivePieces(c) {
		if piece == king {
			continue
		}

		if canUncheck(piece) {
			return true
		}
	}

	return false
}

func slicesEqual(a []Space, b []Space) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// Init initializes a Game object
func (g *Game) Init() {
	*g = Game{
		white: [16]Piece{
			Piece{Type: Rook, Location: Space{0, 0}, Color: White},
			Piece{Type: Knight, Location: Space{0, 1}, Color: White},
			Piece{Type: Bishop, Location: Space{0, 2}, Color: White},
			Piece{Type: Queen, Location: Space{0, 3}, Color: White},
			Piece{Type: King, Location: Space{0, 4}, Color: White},
			Piece{Type: Bishop, Location: Space{0, 5}, Color: White},
			Piece{Type: Knight, Location: Space{0, 6}, Color: White},
			Piece{Type: Rook, Location: Space{0, 7}, Color: White},
			Piece{Type: Pawn, Location: Space{1, 0}, Color: White},
			Piece{Type: Pawn, Location: Space{1, 1}, Color: White},
			Piece{Type: Pawn, Location: Space{1, 2}, Color: White},
			Piece{Type: Pawn, Location: Space{1, 3}, Color: White},
			Piece{Type: Pawn, Location: Space{1, 4}, Color: White},
			Piece{Type: Pawn, Location: Space{1, 5}, Color: White},
			Piece{Type: Pawn, Location: Space{1, 6}, Color: White},
			Piece{Type: Pawn, Location: Space{1, 7}, Color: White},
		},
		black: [16]Piece{
			Piece{Type: Rook, Location: Space{7, 0}, Color: Black},
			Piece{Type: Knight, Location: Space{7, 1}, Color: Black},
			Piece{Type: Bishop, Location: Space{7, 2}, Color: Black},
			Piece{Type: Queen, Location: Space{7, 3}, Color: Black},
			Piece{Type: King, Location: Space{7, 4}, Color: Black},
			Piece{Type: Bishop, Location: Space{7, 5}, Color: Black},
			Piece{Type: Knight, Location: Space{7, 6}, Color: Black},
			Piece{Type: Rook, Location: Space{7, 7}, Color: Black},
			Piece{Type: Pawn, Location: Space{6, 0}, Color: Black},
			Piece{Type: Pawn, Location: Space{6, 1}, Color: Black},
			Piece{Type: Pawn, Location: Space{6, 2}, Color: Black},
			Piece{Type: Pawn, Location: Space{6, 3}, Color: Black},
			Piece{Type: Pawn, Location: Space{6, 4}, Color: Black},
			Piece{Type: Pawn, Location: Space{6, 5}, Color: Black},
			Piece{Type: Pawn, Location: Space{6, 6}, Color: Black},
			Piece{Type: Pawn, Location: Space{6, 7}, Color: Black},
		},
	}
}
