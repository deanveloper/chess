package chess

type castlingRights struct {
	BlackKing, BlackQueen bool
	WhiteKing, WhiteQueen bool
}

// Game represents a game of chess
type Game struct {
	// stored in [file][rank] form
	board     [8][8]Piece
	Castles   castlingRights
	EnPassant Space
	Halfmove  int
	Fullmove  int
}

// Clone returns a new instance of `g`.
func (g *Game) Clone() *Game {
	var newG = &Game{
		board:     g.board,
		EnPassant: g.EnPassant,
		Castles:   g.Castles,
		Halfmove:  g.Halfmove,
	}
	for i, file := range g.board {
		for j, piece := range file {
			newPiece := piece
			newPiece.Game = newG
			newG.board[i][j] = newPiece
		}
	}
	return newG
}

// BoardFileRank returns the game board in it's current state.
// Access board contents with [file][rank]. Useful for determining
// the position of pieces.
func (g *Game) BoardFileRank() [8][8]Piece {
	return g.board
}

// BoardRankFile returns the game board in it's current state.
// Access board contents with [rank][file]. Useful for printing
// the board rank-by-rank.
func (g *Game) BoardRankFile() [8][8]Piece {
	var board [8][8]Piece

	for i, rank := range g.board {
		for j, piece := range rank {
			board[j][i] = piece
		}
	}

	return board
}

// Turn returns who should move next.
func (g *Game) Turn() Color {
	return g.Fullmove%2 == 0
}

// TypedAlivePieces returns all of c's alive pieces with PieceType t.
func (g *Game) TypedAlivePieces(c Color, t PieceType) []Piece {
	pieces := make([]Piece, 0, 16)

	for _, p := range g.AlivePieces(c) {
		if p.Type == t {
			pieces = append(pieces, p)
		}
	}

	return pieces
}

// AlivePieces returns the pieces that c has on the board.
func (g *Game) AlivePieces(c Color) []Piece {
	pieces := make([]Piece, 0, 16)

	for _, file := range g.board {
		for _, piece := range file {
			if piece.Type != PieceNone && piece.Color == c {
				pieces = append(pieces, piece)
			}
		}
	}

	return pieces
}

// PieceAt returns the piece at a given space, and an `ok`
// boolean on if there was a piece on that space at all.
func (g *Game) PieceAt(s Space) (Piece, bool) {
	if !s.Valid() {
		return Piece{}, false
	}

	piece := g.board[s.File][s.Rank]
	if piece.Type == PieceNone {
		return Piece{}, false
	}

	return piece, true
}

// InCheck returns if `c` is in check.
func (g *Game) InCheck(c Color) bool {

	king := g.TypedAlivePieces(c, PieceKing)[0]

	oPieces := g.AlivePieces(c.Other())
	for _, piece := range oPieces {
		seeing := piece.Seeing()
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
	if g.Fullmove < 5 {
		return false
	}

	// 50 move rule
	if g.Halfmove >= 50 && !g.InCheckmate(g.Turn()) {
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
		case whiteAlive[0].Type == PieceBishop,
			whiteAlive[0].Type == PieceKnight,
			whiteAlive[1].Type == PieceBishop,
			whiteAlive[1].Type == PieceKnight:

			return true
		}
	}
	if len(whiteAlive) == 1 && len(blackAlive) == 2 {
		switch {
		case blackAlive[0].Type == PieceBishop,
			blackAlive[0].Type == PieceKnight,
			blackAlive[1].Type == PieceBishop,
			blackAlive[1].Type == PieceKnight:

			return true
		}
	}

	// if all that exist are bishops on the same color...
	var bishops int
	var bishopColor Color

	// white bishops
	for _, piece := range whiteAlive {
		if piece.Type == PieceBishop {
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
	for _, piece := range blackAlive {
		if piece.Type == PieceBishop {
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

	return true
}

func (g *Game) makeMoveUnconditionally(m Move) {
	var pieceTaken bool

	var target *Piece
	target = &g.board[m.To.File][m.To.Rank]

	// update piece
	*target = m.Moving
	target.Location = m.To

	// update piece type for promotions
	if m.Promotion != PieceNone {
		target.Type = m.Promotion
	}

	// handle en passant
	if m.Moving.Type == PiecePawn && m.To == g.EnPassant {
		deadSpace := Space{File: m.To.File, Rank: m.Moving.Location.Rank}
		g.board[deadSpace.File][deadSpace.Rank] = Piece{}
	}

	// update where piece came from
	from := m.Moving.Location
	g.board[from.File][from.Rank] = Piece{}

	// update castling rights
	switch m.Moving.Location {
	case Space{File: 0, Rank: 0}:
		g.Castles.WhiteQueen = false
	case Space{File: 7, Rank: 0}:
		g.Castles.WhiteKing = false
	case Space{File: 0, Rank: 7}:
		g.Castles.BlackQueen = false
	case Space{File: 7, Rank: 7}:
		g.Castles.BlackKing = false
	case Space{File: 4, Rank: 0}:
		g.Castles.WhiteKing = false
		g.Castles.WhiteQueen = false
	case Space{File: 4, Rank: 7}:
		g.Castles.BlackKing = false
		g.Castles.BlackQueen = false
	}

	// update halfmove clock
	if pieceTaken || m.Moving.Type == PiecePawn {
		g.Halfmove = 0
	} else {
		g.Halfmove++
	}
}

// MakeMove makes a move in the game, or returns an error if the move is not possible.
func (g *Game) MakeMove(m Move) error {

	// check to make sure correct color is moving
	var colorToMove Color
	if g.Turn() == White {
		colorToMove = White
	} else {
		colorToMove = Black
	}
	if m.Moving.Color != colorToMove {
		return &MoveError{
			Cause:  m,
			Reason: "it is " + colorToMove.String() + "'s turn",
		}
	}
	// check if the move is valid
	var validMove bool
	seeing := m.Moving.Seeing()
	for _, s := range seeing {
		if m.To == s {
			validMove = true
			break
		}
	}
	if !validMove {
		return &MoveError{
			Cause:  m,
			Reason: "piece cannot see space",
		}
	}

	if m.Moving.Type == PiecePawn && (m.To.Rank == 0 || m.To.Rank == 7) {
		switch m.Promotion {
		case PieceRook, PieceKnight, PieceBishop, PieceQueen:
			break
		default:
			return &MoveError{
				Cause:  m,
				Reason: "cannot promote to " + m.Promotion.String(),
			}
		}
	} else if m.Promotion != PieceNone {
		return &MoveError{
			Cause:  m,
			Reason: "piece cannot promote",
		}
	}

	var inCheck bool
	legalMoves := m.Moving.LegalMoves()
	for _, s := range legalMoves {
		if m.To == s {
			inCheck = true
			break
		}
	}
	if inCheck {
		return &MoveError{
			Cause:   m,
			Reason:  "player is in check",
			InCheck: true,
		}
	}

	// make the move
	g.makeMoveUnconditionally(m)

	// move rook in castles
	if m.Moving.Type == PieceKing {
		diff := m.Moving.Location.File - m.To.File

		if diff == -2 {
			rook, _ := g.PieceAt(Space{File: 0, Rank: m.To.Rank})
			g.makeMoveUnconditionally(Move{
				Moving: rook,
				To:     Space{File: 3, Rank: m.To.Rank},
			})
		}
		if diff == 2 {
			rook, _ := g.PieceAt(Space{File: 7, Rank: m.To.Rank})
			g.makeMoveUnconditionally(Move{
				Moving: rook,
				To:     Space{File: 5, Rank: m.To.Rank},
			})
		}
	}

	return nil
}

func (g *Game) canMove(c Color) bool {
	for _, piece := range g.AlivePieces(c) {
		if len(piece.LegalMoves()) > 0 {
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

// InitClassic initializes a Game object to a classic chess layout,
func (g *Game) InitClassic() {
	*g = Game{}
	rank := [8]PieceType{
		PieceRook,
		PieceKnight,
		PieceBishop,
		PieceQueen,
		PieceKing,
		PieceBishop,
		PieceKnight,
		PieceRook,
	}
	for i, pieceType := range rank {
		putPiece(g, pieceType, White, Space{File: i, Rank: 0})
		putPiece(g, PiecePawn, White, Space{File: i, Rank: 1})

		putPiece(g, pieceType, Black, Space{File: i, Rank: 7})
		putPiece(g, PiecePawn, Black, Space{File: i, Rank: 6})
	}
}

func putPiece(g *Game, p PieceType, c Color, s Space) {
	g.board[s.File][s.Rank] = Piece{
		Game:     g,
		Type:     p,
		Location: s,
		Color:    c,
	}
}
