package chess

import "fmt"

// MoveError represents an error caused by an invalid move.
type MoveError struct {
	Cause   Move
	InCheck bool
}

func (e *MoveError) Error() string {
	if e.InCheck {
		return fmt.Sprintf("cannot move %v to %v: player ends up in check", e.Cause.Moving, e.Cause.To)
	}
	return fmt.Sprintf("cannot move %v to %v: invalid move", e.Cause.Moving, e.Cause.To)
}

// IsInCheckErr returns if the error was caused by the person being in check.
func IsInCheckErr(e error) bool {
	if err, ok := e.(*MoveError); ok {
		return err.InCheck
	}
	return false
}
