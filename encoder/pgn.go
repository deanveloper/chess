package encoder

import (
	"fmt"
	"io"
	"strings"

	"github.com/deanveloper/chess"
)

// PGNReader returns a reader for a game that
// reads the data into PGN notation.
func PGNReader(tags map[string]string, moves []chess.Move) io.Reader {
	return io.MultiReader(tagsReader(tags), &moveTextReader{moves: moves})
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
		} else {
			builder.WriteString(fullTag(key, "omitted"))
		}
	}
	for key, val := range clonedTags {
		builder.WriteString(fullTag(key, val))
	}

	return strings.NewReader(builder.String())
}

type moveTextReader struct {
	moves []chess.Move

	moveIndex int
	strIndex  int

	err error
}

func (r *moveTextReader) Read(b []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}

	var bytesRead int

	for n := len(b); n > 0; {
		if r.moveIndex >= len(r.moves) {
			return bytesRead, io.EOF
		}

		move := r.moves[r.moveIndex]
		alg, err := algShort(move)
		if err != nil {
			r.err = err
			return bytesRead, err
		}
		if r.moveIndex%2 == 0 {
			alg = fmt.Sprintf("%d. %s", r.moveIndex/2+1, alg)
		}
		if r.moveIndex < len(r.moves)-1 {
			alg += " "
		}

		copied := copy(b[bytesRead:], alg[r.strIndex:])
		n -= copied
		r.strIndex += copied
		bytesRead += copied

		if r.strIndex == len(alg) {
			r.moveIndex++
			r.strIndex = 0
		}
	}

	return bytesRead, nil
}

func fullTag(key, value string) string {
	value = strings.NewReplacer("\\", "\\\\", "\"", "\\\"").Replace(value)
	return fmt.Sprintf(`[%s "%s"]`+"\n", key, value)
}
