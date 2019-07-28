package chess

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

// PGNReader returns a reader for a game that
// reads the data into PGN notation.
func PGNReader(game *Game, tags map[string]string) io.Reader {
	order := []string{"Event", "Site", "Date", "Round", "White", "Black", "Result"}
	fullTags := make([]string, 0, len(tags))
	for _, key := range order {
		if val, ok := tags["Event"]; ok {
			fullTags = append(fullTags, fullTag(key, val))
		}
	}
	return &pgnReader{
		game: game,
		tags: fullTags,
	}
}

func fullTag(key, value string) string {
	value = strings.NewReplacer("\\", "\\\\", "\"", "\\\"").Replace(value)
	return fmt.Sprintf(`[%s "%s"]\n`, key, value)
}

type pgnReader struct {
	game *Game
	tags []string

	bytesRead int
	err       error
}

func (r *pgnReader) Read(b []byte) (n int, err error) {
	return 0, nil
}

// FENReader returns a reader for a game that reads
// the data in Forsyth-Edwards Notation.
func FENReader(game *Game) io.Reader {

	board := game.BoardRankFile()
	history := game.History

	var builder strings.Builder

	// 1st field: board state
	for rank := 7; rank >= 0; rank-- {
		var emptySpots byte
		for file := 0; file < 8; file++ {
			p := board[rank][file]
			if p.Type == PieceNone {
				emptySpots++
			} else {
				if emptySpots > 0 {
					builder.WriteByte(emptySpots + '0')
					emptySpots = 0
				}
				name := p.Type.ShortName()
				if p.Color == Black {
					name = unicode.ToLower(name)
				} else {
					name = unicode.ToUpper(name)
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
	if game.Turn() == White {
		builder.WriteByte('w')
	} else {
		builder.WriteByte('b')
	}

	builder.WriteByte(' ')

	// third field: castling availability
	var any bool
	if game.castles.WhiteKing {
		builder.WriteByte('K')
		any = true
	}
	if game.castles.WhiteQueen {
		builder.WriteByte('Q')
		any = true
	}
	if game.castles.BlackKing {
		builder.WriteByte('k')
		any = true
	}
	if game.castles.BlackQueen {
		builder.WriteByte('q')
		any = true
	}
	if !any {
		builder.WriteByte('-')
	}
	builder.WriteByte(' ')

	// fourth field: en passant square
	if game.enPassant.Rank != 0 {
		builder.WriteString(game.enPassant.String())
	} else {
		builder.WriteByte('-')
	}

	builder.WriteByte(' ')

	// fifth field: halfmove clock
	builder.WriteString(strconv.Itoa(game.halfmoveClock))
	builder.WriteByte(' ')

	// sixth field: fullmove number
	fullMove := strconv.Itoa(len(history))
	if fullMove == "0" {
		fullMove = "1"
	}
	builder.WriteString(fullMove)

	return strings.NewReader(builder.String())
}
