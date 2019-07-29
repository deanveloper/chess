package chess

import (
	"strings"
)

// Move represents a move in a chess game.
type Move struct {
	Snapshot Game

	Moving Piece
	To     Space

	Promotion PieceType
}

// AlgebraicShort returns the short algebraic notation for this move.
func (m Move) AlgebraicShort() (string, error) {

	game := m.Snapshot
	player := m.Moving.Color
	piece := m.Moving
	from := m.Moving.Location
	to := m.To

	var builder strings.Builder

	// detect castle, uses a `goto` to skip all the standard move notation...
	// please don't get mad at me for using a goto...
	diff := to.File - from.File
	if diff == -2 {
		builder.WriteString("O-O-O")
		goto checkDetect
	}
	if diff == 2 {
		builder.WriteString("O-O")
		goto checkDetect
	}

	// piece to move
	if piece.Type != PiecePawn {
		builder.WriteByte(byte(piece.Type.ShortName()))
	}

	// disambiguate the piece if needed
	for _, each := range game.TypedAlivePieces(player, piece.Type) {
		if each != piece {
			var seesTarget bool
			for _, space := range each.Seeing() {
				if space == to {
					seesTarget = true
					break
				}
			}
			if seesTarget {
				if from.File != each.Location.File {
					builder.WriteByte(byte(from.File + 'a'))
				} else {
					builder.WriteByte(byte(from.Rank + '1'))
				}
				break
			}
		}
	}

	// if it is a capture
	if game.board[to.File][to.Rank].Type != PieceNone {
		builder.WriteByte('x')
	} else if piece.Type == PiecePawn && game.EnPassant == to {
		builder.WriteByte('x')
	}

	// target square
	builder.WriteString(to.String())

	// en passant
	if piece.Type == PiecePawn && game.EnPassant == to {
		builder.WriteString("e.p.")
	}

	// pawn promotion
	if m.Promotion != PieceNone {
		builder.WriteByte(m.Promotion.ShortName())
	}

checkDetect:
	// in check/mate
	nextState := m.Snapshot.Clone()
	err := nextState.MakeMove(m)
	if err != nil {
		return "", err
	}

	if nextState.InCheck(player.Other()) {
		builder.WriteByte('+')
	}
	if nextState.InCheckmate(player.Other()) {
		builder.WriteByte('+')
	}

	return builder.String(), nil
}
