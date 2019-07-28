package main

import (
	"bufio"
	"os/exec"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

// takes a FEN and returns the UCI for the best move
func runStockfish(fen string, difficulty int) (string, error) {
	cmd := exec.Command("stockfish")

	in, _ := cmd.StdinPipe()
	out, _ := cmd.StdoutPipe()

	cmd.Start()
	in.Write([]byte("setoption name Skill Level value " + strconv.Itoa(difficulty) + "\n"))
	in.Write([]byte("position fen " + fen + "\n"))
	in.Write([]byte("go movetime 3000\n"))

	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 2 || fields[0] != "bestmove" {
			continue
		}

		return fields[1], nil
	}
	if scanner.Err() != nil {
		return "", xerrors.Errorf("error while running stockfish: %w (stockfish not installed?)", scanner.Err())
	}

	cmd.Wait()
	return "", xerrors.Errorf("could not find output from stockfish")
}
