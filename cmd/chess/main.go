package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/deanveloper/chess"
)

func main() {

	scan := bufio.NewScanner(os.Stdin)

	game := &chess.Game{}
	game.Init()

	for scan.Scan() {
		line := scan.Text()

		fields := strings.Fields(line)
		switch fields[0] {
		case "move":
			move := parseMove(game, line)
			err := game.MakeMove(move)
			if err != nil {
				fmt.Println("Error making move:", err)
				continue
			}
		case "":
		}

		fmt.Println("success")
	}
}

func parseMove(g *chess.Game, line string) chess.Move {
	var from, to string
	_, err := fmt.Sscanf(line, "%s to %s", &from, &to)
	if err != nil {
		panic(err)
	}

	fromFile := int(from[0] - 'a')
	toFile := int(to[0] - 'a')

	fromRank := int(from[1] - '1')
	toRank := int(to[1] - '1')

	piece, ok := g.PieceAt(chess.Space{File: fromFile, Rank: fromRank})
	if !ok {
		panic("no piece at " + from)
	}
	return chess.Move{
		Moving: piece,
		To:     chess.Space{File: toFile, Rank: toRank},
	}
}
