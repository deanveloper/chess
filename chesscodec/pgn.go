package chesscodec

import (
	"fmt"
	"io"
	"strings"

	"github.com/deanveloper/chess"
)

// PGNReader returns a reader for a game that
// reads the data into PGN notation.
func PGNReader(game *chess.Game, tags map[string]string) io.Reader {
	return &pgnReader{
		game: game,
	}
}

func tagsReader(tags map[string]string) io.Reader {
	var builder strings.Builder

	clonedTags := make(map[string]string)
	for k, v := range tags {
		clonedTags[k] = v
	}
	order := []string{"Event", "Site", "Date", "Round", "White", "Black", "Result"}
	for _, key := range order {
		if val, ok := clonedTags[key]; ok {
			builder.WriteString(fullTag(key, val))
			delete(clonedTags, key)
		}
	}
	for key, val := range clonedTags {
		builder.WriteString(fullTag(key, val))
	}

	return strings.NewReader(builder.String())
}

func fullTag(key, value string) string {
	value = strings.NewReplacer("\\", "\\\\", "\"", "\\\"").Replace(value)
	return fmt.Sprintf(`[%s "%s"]\n`, key, value)
}

type pgnReader struct {
	game *chess.Game

	bytesRead int
	err       error
}

func (r *pgnReader) Read(b []byte) (n int, err error) {
	if r.err != nil {
		return 0, r.err
	}

	return 0, nil
}
