package encoder

import (
	"io"
	"strconv"
	"strings"

	"github.com/deanveloper/chess"
)

// FENReader returns a reader for a game that reads
// the data in Forsyth-Edwards Notation.
func FENReader(game *chess.Game) io.Reader {

	board := game.BoardRankFile()

	var builder strings.Builder

	// 1st field: board state
	for rank := 7; rank >= 0; rank-- {
		var emptySpots byte
		for file := 0; file < 8; file++ {
			p := board[rank][file]
			if p.Type == chess.PieceNone {
				emptySpots++
			} else {
				if emptySpots > 0 {
					builder.WriteByte(emptySpots + '0')
					emptySpots = 0
				}
				name := p.Type.ShortName()
				if p.Color == chess.Black {
					name = name - 'A' + 'a' // lowercase
				}
				builder.WriteByte(byte(name))
			}
		}
		if emptySpots > 0 {
			builder.WriteByte(emptySpots + '0')
		}
		builder.WriteByte('/')
	}

	builder.WriteByte(' ')

	// second field: player to move
	if game.Turn() == chess.White {
		builder.WriteByte('w')
	} else {
		builder.WriteByte('b')
	}

	builder.WriteByte(' ')

	// third field: castling availability
	var any bool
	if game.Castles.WhiteKing {
		builder.WriteByte('K')
		any = true
	}
	if game.Castles.WhiteQueen {
		builder.WriteByte('Q')
		any = true
	}
	if game.Castles.BlackKing {
		builder.WriteByte('k')
		any = true
	}
	if game.Castles.BlackQueen {
		builder.WriteByte('q')
		any = true
	}
	if !any {
		builder.WriteByte('-')
	}
	builder.WriteByte(' ')

	// fourth field: en passant square
	if game.EnPassant.Rank != 0 {
		builder.WriteString(game.EnPassant.String())
	} else {
		builder.WriteByte('-')
	}

	builder.WriteByte(' ')

	// fifth field: halfmove clock
	builder.WriteString(strconv.Itoa(game.Halfmove))
	builder.WriteByte(' ')

	// sixth field: fullmove number
	fullMove := game.Fullmove
	if fullMove == 0 {
		fullMove = 1
	}
	builder.WriteString(strconv.Itoa(fullMove))

	return strings.NewReader(builder.String())
}
