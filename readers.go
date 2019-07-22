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

	board := game.Board()
	history := game.History

	var builder strings.Builder

	// 1st field: board state
	for file := 0; file < 8; file++ {
		var emptySpots byte
		for rank := 7; rank >= 0; rank-- {
			p := board[file][rank]
			if p.Type == None {
				emptySpots++
			} else {
				if emptySpots > 0 {
					builder.WriteByte(emptySpots + '0')
					emptySpots = 0
				}
				name := p.Type.ShortName()
				if p.Color == Black {
					name = unicode.ToUpper(name)
				} else {
					name = unicode.ToLower(name)
				}
				builder.WriteByte(byte(name))
			}
		}
		builder.WriteByte('/')
	}

	builder.WriteByte(' ')

	// second field: player to move
	if len(history)%2 == 0 {
		builder.WriteByte('w')
	} else {
		builder.WriteByte('b')
	}

	builder.WriteByte(' ')

	// third field: castling availability
	bKing := game.TypedAlivePieces(Black, King)[0]
	wKing := game.TypedAlivePieces(White, King)[0]
	var bkCastle, bqCastle, wkCastle, wqCastle bool
	for _, space := range bKing.Seeing(game) {
		diff := bKing.Location.File - space.File
		if diff == 2 {
			bqCastle = true
		}
		if diff == -2 {
			bkCastle = true
		}
	}
	for _, space := range wKing.Seeing(game) {
		diff := wKing.Location.File - space.File
		if diff == 2 {
			wqCastle = true
		}
		if diff == -2 {
			wkCastle = true
		}
	}
	if wkCastle {
		builder.WriteByte('K')
	}
	if wqCastle {
		builder.WriteByte('Q')
	}
	if bkCastle {
		builder.WriteByte('k')
	}
	if bqCastle {
		builder.WriteByte('q')
	}

	builder.WriteByte(' ')

	// fourth field: en passant square
	passant := Taken
	if len(history) > 0 {
		move := history[len(history)-1]
		if move.Moving.Type == Pawn {
			diff := move.Moving.Location.Rank - move.To.Rank
			if diff == 2 || diff == -2 {
				passant = Space{
					File: move.To.File,
					Rank: move.To.Rank - (diff / 2),
				}
			}
		}
	}
	if passant == Taken {
		builder.WriteString(passant.String())
	} else {
		builder.WriteByte('-')
	}

	builder.WriteByte(' ')

	// fifth field: halfmove clock
	builder.WriteString(strconv.Itoa(game.halfmoveClock))
	builder.WriteByte(' ')

	// sixth field: fullmove number
	builder.WriteString(strconv.Itoa(len(history) - 1))

	return strings.NewReader(builder.String())
}
