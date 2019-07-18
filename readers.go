package chess

import (
	"io"
)

type fenReader struct {
}

// FENReader returns a reader for a game that reads
// the data in Forsyth-Edwards Notation.
func FENReader(game *Game, info map[string]string) io.Reader {
	return nil
}
