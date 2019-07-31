package encoder

import (
	"fmt"
	"strings"

	"github.com/deanveloper/chess"
)

type algebraicError struct {
	algebraic string
	reason    string
}

func (a algebraicError) Error() string {
	return fmt.Sprintf("error parsing %q: %s", a.algebraic, a.reason)
}

// FromAlgebraic returns a move from an algebraic string.
func FromAlgebraic(g *chess.Game, algebraic string) (chess.Move, error) {

	// remove checkmate and en passant symbols
	algebraic = strings.TrimSuffix(algebraic, "+")
	algebraic = strings.TrimSuffix(algebraic, "+")
	algebraic = strings.TrimSuffix(algebraic, "#")
	algebraic = strings.TrimSuffix(algebraic, "e.p.")

	if len(algebraic) < 2 {
		return chess.Move{}, algebraicError{algebraic: algebraic, reason: "too short"}
	}

	var promotion chess.PieceType
	switch algebraic[len(algebraic)-1] {
	case 'R':
		promotion = chess.PieceRook
		algebraic = algebraic[:len(algebraic)-1]
	case 'N':
		promotion = chess.PieceKnight
		algebraic = algebraic[:len(algebraic)-1]
	case 'B':
		promotion = chess.PieceBishop
		algebraic = algebraic[:len(algebraic)-1]
	case 'Q':
		promotion = chess.PieceQueen
		algebraic = algebraic[:len(algebraic)-1]
	}

	target := chess.Space{
		File: int(algebraic[len(algebraic)-2] - 'a'),
		Rank: int(algebraic[len(algebraic)-1] - '1'),
	}
	if !target.Valid() {
		return chess.Move{}, algebraicError{
			algebraic: algebraic,
			reason:    "invalid target square " + algebraic[len(algebraic)-2:],
		}
	}

	var pieceType chess.PieceType
	switch algebraic[0] {
	case 'R':
		pieceType = chess.PieceRook
	case 'N':
		pieceType = chess.PieceKnight
	case 'B':
		pieceType = chess.PieceBishop
	case 'Q':
		pieceType = chess.PieceQueen
	case 'K':
		pieceType = chess.PieceKing
	default:
		pieceType = chess.PiecePawn
	}

	file := -1
	rank := -1
	char := algebraic[1]
	if pieceType == chess.PiecePawn {
		char = algebraic[0]
	}
	if char >= 'a' && char <= 'h' {
		file = int(char - 'a')
	}
	if char >= '1' && char <= '8' {
		rank = int(char - '1')
	}

	var pieceFound bool
	var piece chess.Piece
	for _, eachPiece := range g.TypedAlivePieces(g.Turn(), pieceType) {
		for _, space := range eachPiece.Seeing() {
			if space == target {
				if file >= 0 && file != space.File {
					continue
				}
				if rank >= 0 && rank != space.Rank {
					continue
				}
				if pieceFound {
					return chess.Move{}, algebraicError{
						algebraic: algebraic,
						reason:    "move is ambiguous",
					}
				}
				piece = eachPiece
				pieceFound = true
			}
		}
	}

	if !pieceFound {
		return chess.Move{}, algebraicError{
			algebraic: algebraic,
			reason:    "could not find a " + pieceType.String() + " targetting " + target.String(),
		}
	}

	// if a piece is being captured, an x must appear as the 2nd or 3rd character
	algCapturing := algebraic[1] == 'x' || (len(algebraic) > 2 && algebraic[2] == 'x')

	targetPiece, _ := g.PieceAt(target)
	actCapturing := targetPiece.Type != chess.PieceNone && targetPiece.Color != g.Turn()
	if pieceType == chess.PiecePawn && target == g.EnPassant {
		actCapturing = true
	}

	if algCapturing != actCapturing {
		return chess.Move{}, algebraicError{
			algebraic: algebraic,
			reason:    "incorrect capturing information given",
		}
	}

	return chess.Move{
		Snapshot:  *g,
		Moving:    piece,
		To:        target,
		Promotion: promotion,
	}, nil
}

// Algebraic returns the algebraic form for a given move. Does not detect
// if the move puts the other person in check.
func Algebraic(m chess.Move) string {
	game := m.Snapshot
	player := m.Moving.Color
	piece := m.Moving
	from := m.Moving.Location
	to := m.To

	var builder strings.Builder

	// detect castle, uses a `goto` to skip all the standard move notation...
	// please don't get mad at me for using a goto...
	if piece.Type == chess.PieceKing {
		diff := to.File - from.File
		if diff == -2 {
			builder.WriteString("O-O-O")
			return "O-O-O"
		}
		if diff == 2 {
			builder.WriteString("O-O")
			return "O-O-O"
		}
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

	return builder.String()
}

// PGNAlgebraic returns the algebraic notation used for PGN notation. This
// means that checks are included, and a game result is appended if it is the
// final move of the game.
func PGNAlgebraic(m chess.Move) (string, error) {

	player := m.Moving.Color

	alg := Algebraic(m)

	// in check/mate
	nextState := m.Snapshot.Clone()
	nextState.MakeMoveUnconditionally(m)

	if nextState.InCheck(player.Other()) {
		alg += "+"
	}
	if nextState.InCheckmate(player.Other()) {
		if player == chess.White {
			alg += "+ 1-0"
		} else {
			alg += "+ 0-1"
		}
	}
	if nextState.InStalemate(player.Other()) {
		alg += " 1/2-1/2"
	}

	return alg, nil
}
