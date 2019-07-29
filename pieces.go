package chess

import "fmt"

// Color is a type which can be either black or white.
type Color bool

// The enum of colors
const (
	Black Color = false
	White Color = true
)

// Other returns the color that is not c.
func (c Color) Other() Color {
	return !c
}

func (c Color) String() string {
	if c == Black {
		return "Black"
	}
	return "White"
}

// The enum of pieces
const (
	PieceNone PieceType = iota
	PiecePawn
	PieceRook
	PieceKnight
	PieceBishop
	PieceQueen
	PieceKing
)

// PieceType represents a type of piece
type PieceType byte

// Symbol returns a rune representing the piece
func (p PieceType) Symbol() rune {
	return [...]rune{' ', '♟', '♜', '♞', '♝', '♛', '♚'}[p]
}

// ShortName returns the shortname for p used by Forsyth-Edwards Notation.
func (p PieceType) ShortName() byte {
	return [...]byte{'X', 'P', 'R', 'N', 'B', 'Q', 'K'}[p]
}

func (p PieceType) String() string {
	return [...]string{"None", "Pawn", "Rook", "Knight", "Bishop", "Queen", "King"}[p]
}

// Piece represents a chess piece.
type Piece struct {
	Game     *Game
	Type     PieceType
	Location Space
	Color    Color
}

func (p Piece) String() string {
	return fmt.Sprintf("%v %v on %v", p.Color, p.Type, p.Location)
}

// Seeing returns all spaces that this piece can see. Just because
// the piece can see a square does not mean the move is valid; as the
// player may be in check, or moving the piece may put the player in check.
func (p Piece) Seeing() []Space {
	var moveTo []Space

	cur := p.Location

	switch p.Type {
	case PiecePawn:
		// allow moving one up if there is not a piece there
		next := p.Location
		if p.Color == Black {
			next.Rank--
		} else {
			next.Rank++
		}
		if _, ok := p.Game.PieceAt(next); !ok {
			moveTo = append(moveTo, next)
		}

		// if unmoved, allow moving two up
		if (p.Color == White && cur.Rank == 1) || (p.Color == Black && cur.Rank == 6) {
			twoUp := next
			if p.Color == Black {
				twoUp.Rank--
			} else {
				twoUp.Rank++
			}
			if _, ok := p.Game.PieceAt(next); !ok {
				if _, ok2 := p.Game.PieceAt(twoUp); !ok2 {
					moveTo = append(moveTo, twoUp)
				}
			}
		}

		// allow diagonals if it can take
		diagL := Space{Rank: next.Rank, File: next.File - 1}
		diagR := Space{Rank: next.Rank, File: next.File + 1}
		if _, ok := p.Game.PieceAt(diagL); ok {
			moveTo = append(moveTo, diagL)
		}
		if _, ok := p.Game.PieceAt(diagR); ok {
			moveTo = append(moveTo, diagR)
		}

		// include possibility of en passant
		if diagL == p.Game.EnPassant || diagR == p.Game.EnPassant {
			moveTo = append(moveTo, p.Game.EnPassant)
		}

	case PieceRook:
		moveTo = append(moveTo, p.loop(func(s Space) Space {
			return Space{File: s.File + 1, Rank: s.Rank}
		})...)
		moveTo = append(moveTo, p.loop(func(s Space) Space {
			return Space{File: s.File - 1, Rank: s.Rank}
		})...)
		moveTo = append(moveTo, p.loop(func(s Space) Space {
			return Space{File: s.File, Rank: s.Rank - 1}
		})...)
		moveTo = append(moveTo, p.loop(func(s Space) Space {
			return Space{File: s.File, Rank: s.Rank + 1}
		})...)
	case PieceKnight:
		moveTo = []Space{
			Space{File: cur.File + 1, Rank: cur.Rank + 2},
			Space{File: cur.File + 2, Rank: cur.Rank + 1},
			Space{File: cur.File + 1, Rank: cur.Rank - 2},
			Space{File: cur.File + 2, Rank: cur.Rank - 1},
			Space{File: cur.File - 1, Rank: cur.Rank + 2},
			Space{File: cur.File - 2, Rank: cur.Rank + 1},
			Space{File: cur.File - 1, Rank: cur.Rank - 2},
			Space{File: cur.File - 2, Rank: cur.Rank - 1},
		}
	case PieceBishop:
		moveTo = append(moveTo, p.loop(func(s Space) Space {
			return Space{File: s.File + 1, Rank: s.Rank + 1}
		})...)
		moveTo = append(moveTo, p.loop(func(s Space) Space {
			return Space{File: s.File + 1, Rank: s.Rank - 1}
		})...)
		moveTo = append(moveTo, p.loop(func(s Space) Space {
			return Space{File: s.File - 1, Rank: s.Rank + 1}
		})...)
		moveTo = append(moveTo, p.loop(func(s Space) Space {
			return Space{File: s.File - 1, Rank: s.Rank - 1}
		})...)
	case PieceQueen:
		moveTo = append(moveTo, p.loop(func(s Space) Space {
			return Space{File: s.File + 1, Rank: s.Rank}
		})...)
		moveTo = append(moveTo, p.loop(func(s Space) Space {
			return Space{File: s.File - 1, Rank: s.Rank}
		})...)
		moveTo = append(moveTo, p.loop(func(s Space) Space {
			return Space{File: s.File, Rank: s.Rank - 1}
		})...)
		moveTo = append(moveTo, p.loop(func(s Space) Space {
			return Space{File: s.File, Rank: s.Rank + 1}
		})...)
		moveTo = append(moveTo, p.loop(func(s Space) Space {
			return Space{File: s.File + 1, Rank: s.Rank + 1}
		})...)
		moveTo = append(moveTo, p.loop(func(s Space) Space {
			return Space{File: s.File + 1, Rank: s.Rank - 1}
		})...)
		moveTo = append(moveTo, p.loop(func(s Space) Space {
			return Space{File: s.File - 1, Rank: s.Rank + 1}
		})...)
		moveTo = append(moveTo, p.loop(func(s Space) Space {
			return Space{File: s.File - 1, Rank: s.Rank - 1}
		})...)
	case PieceKing:
		// invalid spaces and spaces with own color's pieces
		// are removed later
		moveTo = []Space{
			Space{File: cur.File + 1, Rank: cur.Rank + 1},
			Space{File: cur.File + 1, Rank: cur.Rank - 1},
			Space{File: cur.File - 1, Rank: cur.Rank + 1},
			Space{File: cur.File - 1, Rank: cur.Rank - 1},
			Space{File: cur.File + 1, Rank: cur.Rank},
			Space{File: cur.File - 1, Rank: cur.Rank},
			Space{File: cur.File, Rank: cur.Rank + 1},
			Space{File: cur.File, Rank: cur.Rank - 1},
		}

		// castling:
		if p.Color == White {
			if p.Game.Castles.WhiteQueen {
				moveTo = append(moveTo, Space{File: 2, Rank: cur.Rank})
			}
			if p.Game.Castles.WhiteKing {
				moveTo = append(moveTo, Space{File: 6, Rank: cur.Rank})
			}
		} else {
			if p.Game.Castles.BlackQueen {
				moveTo = append(moveTo, Space{File: 2, Rank: cur.Rank})
			}
			if p.Game.Castles.BlackKing {
				moveTo = append(moveTo, Space{File: 6, Rank: cur.Rank})
			}
		}
	}

	// remove anything where a piece is being taken of the same color,
	// or anything off of the board.
	sees := make([]Space, 0, len(moveTo))
	for i := 0; i < len(moveTo); i++ {
		s := moveTo[i]
		if !s.Valid() {
			continue
		}
		if other, ok := p.Game.PieceAt(s); ok && other.Color == p.Color {
			continue
		}
		sees = append(sees, s)
	}

	return sees
}

// LegalMoves returns all of the legal moves for p.
func (p Piece) LegalMoves() []Space {
	var legal []Space

	for _, space := range p.Seeing() {

		// special case - remove castles in their cases
		if p.Type == PieceKing {
			diff := p.Location.File - space.File

			// remove ability if pieces are between the rook and king
			if diff == -2 {
				if p := p.Game.board[1][p.Location.Rank]; p.Type != PieceNone {
					continue
				}
				if p := p.Game.board[2][p.Location.Rank]; p.Type != PieceNone {
					continue
				}
				if p := p.Game.board[3][p.Location.Rank]; p.Type != PieceNone {
					continue
				}
			}
			if diff == 2 {
				if p := p.Game.board[5][p.Location.Rank]; p.Type != PieceNone {
					continue
				}
				if p := p.Game.board[6][p.Location.Rank]; p.Type != PieceNone {
					continue
				}
			}

			// no castling at all while in check
			if diff == -2 || diff == 2 {
				if p.Game.InCheck(p.Color) {
					continue
				}
			}

			// queen-side castle
			if diff == -2 {
				clone := p.Game.Clone()
				clone.makeMoveUnconditionally(Move{
					Moving: p,
					To:     Space{File: 3, Rank: p.Location.Rank},
				})
				if clone.InCheck(p.Color) {
					continue
				}

				clone = p.Game.Clone()
				clone.makeMoveUnconditionally(Move{
					Moving: p,
					To:     Space{File: 2, Rank: p.Location.Rank},
				})
				if clone.InCheck(p.Color) {
					continue
				}
			}
			// king-side castle
			if diff == 2 {
				clone := p.Game.Clone()
				clone.makeMoveUnconditionally(Move{
					Moving: p,
					To:     Space{File: 5, Rank: p.Location.Rank},
				})
				if clone.InCheck(p.Color) {
					continue
				}

				clone = p.Game.Clone()
				clone.makeMoveUnconditionally(Move{
					Moving: p,
					To:     Space{File: 6, Rank: p.Location.Rank},
				})
				if clone.InCheck(p.Color) {
					continue
				}
			}
		}

		newG := p.Game.Clone()
		newG.makeMoveUnconditionally(Move{
			Snapshot: *newG,
			Moving:   p,
			To:       space,
		})
		if newG.InCheck(p.Color) {
			continue
		}
		legal = append(legal, space)
	}
	return legal
}

func (p Piece) loop(next func(Space) Space) []Space {
	var spaces []Space

	cur := p.Location
	for {
		cur = next(cur)
		if !cur.Valid() {
			break
		}

		spaces = append(spaces, cur)

		if _, ok := p.Game.PieceAt(cur); ok {
			break
		}
	}

	return spaces
}
