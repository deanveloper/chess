package chess

import (
	"fmt"
)

// Move represents a move in a chess game.
type Move struct {
	g *Game

	Moving Piece
	To     Space

	Promotion PieceType
}

// AlgebraicString represents the algebraic string for a given move.
// Uses FIDE notation, however it does not deal with "check" or "checkmate" cases.
// Does not validate the move; invalid moves may result in strange strings.
//
// For an implementation which validates moves and handles "check" and "checkmate",
// use `game.Move`.
func (m Move) AlgebraicString() string {

	// promotions
	if m.Promotion != Pawn {
		return fmt.Sprintf("%v%v", m.To, m.Promotion.ShortName())
	}

	// castles
	if m.Moving.Type == King {
		diff := m.Moving.Location.File - m.To.File
		if diff == 2 || diff == -2 {
			if m.To.File == 2 {
				return "O-O-O"
			}
			if m.To.File == 6 {
				return "O-O"
			}
		}
	}

	// captures
	if _, ok := m.g.PieceAt(m.To); ok {
		if m.Moving.Type == Pawn {
			return fmt.Sprintf("%cx%v", 'a'+m.Moving.Location.File, m.To)
		}
		return fmt.Sprintf("%sx%v", m.Moving.Type.ShortName(), m.To)
	}

	// en passant captures
	if m.Moving.Type == Pawn {
		target := Space{Rank: m.To.Rank, File: m.Moving.Location.File}
		if piece, ok := m.g.PieceAt(target); ok && piece.Type == Pawn {
			return fmt.Sprintf("%cx%se.p.", 'a'+m.Moving.Location.File, m.To.String())
		}
	}

	// normal move
	return fmt.Sprintf("%v%v", m.Moving.Type.ShortName(), m.To)
}
