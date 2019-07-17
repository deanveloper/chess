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
	None PieceType = iota
	Pawn
	Rook
	Knight
	Bishop
	Queen
	King
)

// PieceType represents a type of piece
type PieceType byte

// ShortName returns the short name of this piece used in alegbraic notation.
func (p PieceType) ShortName() string {
	return [...]string{"none", "", "R", "N", "B", "Q", "K"}[p]
}

func (p PieceType) String() string {
	return [...]string{"None", "Pawn", "Rook", "Knight", "Bishop", "Queen", "King"}[p]
}

// Piece represents a chess piece.
type Piece struct {
	Type     PieceType
	Location Space
	Color    Color
}

func (p Piece) String() string {
	return fmt.Sprintf("%v %v on %v", p.Color, p.Type, p.Location)
}

// Symbol returns a rune representing the piece
func (p Piece) Symbol() rune {
	if p.Color == White {
		return [...]rune{' ', '♙', '♖', '♘', '♗', '♕', '♔'}[p.Type]
	}

	return [...]rune{' ', '♟', '♜', '♞', '♝', '♛', '♚'}[p.Type]
}

// History returns all of the movement history that this piece has made.
func (p Piece) History(g *Game) []Move {
	var trackedSpace = p.Location
	var history []Move

	// traverse backward through history
	for i := len(g.History) - 1; i >= 0; i-- {
		move := g.History[i]
		if trackedSpace == move.To {
			history = append([]Move{move}, history...)
		}
	}

	return history
}

// Seeing returns all spaces that this piece can see. Just because
// the piece can see a square does not mean the move is valid; as the
// player may be in check, or moving the piece may put the player in check.
func (p Piece) Seeing(g *Game) []Space {
	var moveTo []Space

	cur := p.Location

	switch p.Type {
	case Pawn:
		// allow moving one up if there is not a piece there
		next := p.Location
		if p.Color == Black {
			next.Rank--
		} else {
			next.Rank++
		}
		if _, ok := g.PieceAt(next); !ok {
			moveTo = append(moveTo, next)
		}

		// if unmoved, allow moving two up
		if len(p.History(g)) == 0 {
			twoUp := next
			if p.Color == Black {
				twoUp.Rank--
			} else {
				twoUp.Rank++
			}
			if _, ok := g.PieceAt(next); !ok {
				if _, ok2 := g.PieceAt(twoUp); !ok2 {
					moveTo = append(moveTo, twoUp)
				}
			}
		}

		// allow diagonals if it can take
		diagL := Space{Rank: next.Rank, File: next.File - 1}
		diagR := Space{Rank: next.Rank, File: next.File + 1}
		if _, ok := g.PieceAt(diagL); ok {
			moveTo = append(moveTo, diagL)
		}
		if _, ok := g.PieceAt(diagR); ok {
			moveTo = append(moveTo, diagR)
		}

		// include possibility of en passant
		if len(g.History) > 0 {
			var validRank int
			if p.Color == Black {
				validRank = 3
			} else {
				validRank = 4
			}
			lastMove := g.History[len(g.History)-1]
			if lastMove.To.Rank == validRank && p.Location.Rank == validRank {
				lastOpposingLoc := lastMove.Moving.Location
				switch {
				case p.Color == Black && lastOpposingLoc.Rank == 1:
					moveTo = append(moveTo, Space{Rank: 2, File: lastOpposingLoc.File})
				case p.Color == White && lastOpposingLoc.Rank == 6:
					moveTo = append(moveTo, Space{Rank: 5, File: lastOpposingLoc.File})
				}
			}
		}

	case Rook:
		moveTo = append(moveTo, p.loop(g, func(s Space) Space {
			return Space{File: s.File + 1, Rank: s.Rank}
		})...)
		moveTo = append(moveTo, p.loop(g, func(s Space) Space {
			return Space{File: s.File - 1, Rank: s.Rank}
		})...)
		moveTo = append(moveTo, p.loop(g, func(s Space) Space {
			return Space{File: s.File, Rank: s.Rank - 1}
		})...)
		moveTo = append(moveTo, p.loop(g, func(s Space) Space {
			return Space{File: s.File, Rank: s.Rank + 1}
		})...)
	case Knight:
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
	case Bishop:
		moveTo = append(moveTo, p.loop(g, func(s Space) Space {
			return Space{File: s.File + 1, Rank: s.Rank + 1}
		})...)
		moveTo = append(moveTo, p.loop(g, func(s Space) Space {
			return Space{File: s.File + 1, Rank: s.Rank - 1}
		})...)
		moveTo = append(moveTo, p.loop(g, func(s Space) Space {
			return Space{File: s.File - 1, Rank: s.Rank + 1}
		})...)
		moveTo = append(moveTo, p.loop(g, func(s Space) Space {
			return Space{File: s.File - 1, Rank: s.Rank - 1}
		})...)
	case Queen:
		moveTo = append(moveTo, p.loop(g, func(s Space) Space {
			return Space{File: s.File + 1, Rank: s.Rank}
		})...)
		moveTo = append(moveTo, p.loop(g, func(s Space) Space {
			return Space{File: s.File - 1, Rank: s.Rank}
		})...)
		moveTo = append(moveTo, p.loop(g, func(s Space) Space {
			return Space{File: s.File, Rank: s.Rank - 1}
		})...)
		moveTo = append(moveTo, p.loop(g, func(s Space) Space {
			return Space{File: s.File, Rank: s.Rank + 1}
		})...)
		moveTo = append(moveTo, p.loop(g, func(s Space) Space {
			return Space{File: s.File + 1, Rank: s.Rank + 1}
		})...)
		moveTo = append(moveTo, p.loop(g, func(s Space) Space {
			return Space{File: s.File + 1, Rank: s.Rank - 1}
		})...)
		moveTo = append(moveTo, p.loop(g, func(s Space) Space {
			return Space{File: s.File - 1, Rank: s.Rank + 1}
		})...)
		moveTo = append(moveTo, p.loop(g, func(s Space) Space {
			return Space{File: s.File - 1, Rank: s.Rank - 1}
		})...)
	case King:
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
		if len(p.History(g)) > 0 {
			break
		}
		if qRook, ok := g.PieceAt(Space{File: 0, Rank: cur.Rank}); ok {
			if qRook.Type == Rook && len(qRook.History(g)) == 0 {
				moveTo = append(moveTo, Space{File: 2, Rank: cur.Rank})
			}
		}
		if kRook, ok := g.PieceAt(Space{File: 7, Rank: cur.Rank}); ok {
			if kRook.Type == Rook && len(kRook.History(g)) == 0 {
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
		if other, ok := g.PieceAt(s); ok && other.Color == p.Color {
			continue
		}
		sees = append(sees, s)
	}

	return sees
}

// LegalMoves returns all of the legal moves for p.
func (p Piece) LegalMoves(g *Game) []Space {
	var legal []Space

	for _, space := range p.Seeing(g) {

		// special case - remove castles if in check or the
		// middle square of the castle puts you in check
		if p.Type == King {
			diff := p.Location.File - space.File

			// no castling at all while in check
			if diff == -2 || diff == 2 {
				if g.InCheck(p.Color) {
					continue
				}
			}

			// queen-side castle
			if diff == -2 {
				clone := g.Clone(false)
				clone.makeMoveUnconditionally(Move{
					Moving: p,
					To:     Space{File: 3, Rank: p.Location.Rank},
				})
				if clone.InCheck(p.Color) {
					continue
				}

				clone = g.Clone(false)
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
				clone := g.Clone(false)
				clone.makeMoveUnconditionally(Move{
					Moving: p,
					To:     Space{File: 5, Rank: p.Location.Rank},
				})
				if clone.InCheck(p.Color) {
					continue
				}

				clone = g.Clone(false)
				clone.makeMoveUnconditionally(Move{
					Moving: p,
					To:     Space{File: 6, Rank: p.Location.Rank},
				})
				if clone.InCheck(p.Color) {
					continue
				}
			}
		}

		newG := g.Clone(false)
		newG.makeMoveUnconditionally(Move{
			Moving: p,
			To:     space,
		})
		if newG.InCheck(p.Color) {
			continue
		}
	}
	return legal
}

func (p Piece) loop(g *Game, next func(Space) Space) []Space {
	var spaces []Space

	cur := p.Location
	for {
		cur = next(cur)
		if !cur.Valid() {
			break
		}

		spaces = append(spaces, cur)

		if _, ok := g.PieceAt(cur); ok {
			break
		}
	}

	return spaces
}
