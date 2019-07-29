package encoder

import (
	"strings"

	"github.com/deanveloper/chess"
)

// AlgebraicShort returns the short algebraic notation for this move.
func AlgebraicShort(m chess.Move) (string, error) {

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
	if piece.Type != chess.PiecePawn {
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
	if game.BoardFileRank()[to.File][to.Rank].Type != chess.PieceNone {
		builder.WriteByte('x')
	} else if piece.Type == chess.PiecePawn && game.EnPassant == to {
		builder.WriteByte('x')
	}

	// target square
	builder.WriteString(to.String())

	// en passant
	if piece.Type == chess.PiecePawn && game.EnPassant == to {
		builder.WriteString("e.p.")
	}

	// pawn promotion
	if m.Promotion != chess.PieceNone {
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
