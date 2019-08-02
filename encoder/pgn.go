package encoder

import (
	"fmt"
	"io"
	"strings"

	"github.com/deanveloper/chess"
)

// PGNReader returns a reader for a game that
// reads the data into PGN notation.
func PGNReader(tags map[string]string, moves <-chan chess.Move, completion <-chan chess.CompletionState) io.Reader {
	return io.MultiReader(tagsReader(tags), &moveTextReader{moves: moves}, completionStateReader(completion))
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

	builder.WriteByte('\n')
	return strings.NewReader(builder.String())
}

func fullTag(key, value string) string {
	value = strings.NewReplacer("\\", "\\\\", "\"", "\\\"").Replace(value)
	return fmt.Sprintf(`[%s "%s"]`+"\n", key, value)
}

type moveTextReader struct {
	moves <-chan chess.Move

	movesRead int
	curMove   string
	strIndex  int

	err error
}

func (r *moveTextReader) Read(b []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}

	var bytesRead int

	// step 1: finish reading last move
	if r.strIndex < len(r.curMove) {
		copied := copy(b[bytesRead:], r.curMove[r.strIndex:])
		r.strIndex += copied
		bytesRead += copied
	}

	for bytesRead < len(b) {
		var move chess.Move
		move, ok := <-r.moves
		if !ok {
			r.err = io.EOF
			return bytesRead, io.EOF
		}

		alg, err := PGNAlgebraic(move)
		if err != nil {
			r.err = err
			return bytesRead, err
		}
		if r.movesRead%2 == 0 {
			alg = fmt.Sprintf("%d. %s", r.movesRead/2+1, alg)
		}
		if r.movesRead < len(r.moves)-1 {
			alg += " "
		}

		copied := copy(b[bytesRead:], alg[r.strIndex:])
		r.strIndex += copied
		bytesRead += copied

		if r.strIndex == len(alg) {
			r.movesRead++
			r.strIndex = 0
		}
	}

	return bytesRead, nil
}

func completionStateReader(completion <-chan chess.CompletionState) io.Reader {
	var final string

	complete := <-completion
	if complete.Done {
		switch {
		case complete.Draw:
			final = " 1/2-1/2"
		case complete.Winner == chess.White:
			final = " 1-0"
		case complete.Winner == chess.Black:
			final = " 0-1"
		}
	}

	return strings.NewReader(final)
}
