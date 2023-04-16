package chess

import "fmt"

// Space represents a space on the chess board.
//
// Rank and File are zero-indexed, so a1 is (0,0)
type Space struct {
	File, Rank int
}

// Valid returns if the space is a valid space to be occupying
func (s Space) Valid() bool {
	return s.Rank >= 0 && s.Rank < 8 && s.File >= 0 && s.File < 8
}

// Color returns the color of this space
func (s Space) Color() Color {
	if !s.Valid() {
		panic(fmt.Sprintf("invalid space: %+v", s))
	}

	// (0,0) is black which is represented by `false`
	return (s.Rank+s.File)%2 != 0
}

func (s Space) String() string {
	if !s.Valid() {
		panic(fmt.Sprintf("invalid space {File:%d Rank:%d}", s.File, s.Rank))
	}

	return fmt.Sprintf("%c%d", ('a' + s.File), s.Rank+1)
}
