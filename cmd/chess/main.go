package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/deanveloper/chess"
)

func main() {

	scan := bufio.NewScanner(os.Stdin)

	game := &chess.Game{}
	game.Init()

	fmt.Print("> ")
	for scan.Scan() {
		line := scan.Text()

		fields := strings.Fields(line)
		switch fields[0] {
		case "move":
			move := parseMove(game, fields[1], fields[2])
			err := game.MakeMove(move)
			if err != nil {
				fmt.Println("Error making move:", err)
				fmt.Print("> ")
				continue
			}
		case "pieces":
			fmt.Println("white:")
			for _, piece := range game.AlivePieces(chess.White) {
				fmt.Println("\t", piece)
			}
			fmt.Println("black:")
			for _, piece := range game.AlivePieces(chess.Black) {
				fmt.Println("\t", piece)
			}
		case "print":
			const top = "  ┌───┬───┬───┬───┬───┬───┬───┬───┐"
			const midPieces = " │ X │ X │ X │ X │ X │ X │ X │ X │"
			const mid = "  ├───┼───┼───┼───┼───┼───┼───┼───┤"
			const bottom = "  └───┴───┴───┴───┴───┴───┴───┴───┘"
			const files = "    a   b   c   d   e   f   g   h"

			board := rotate(game.Board())
			fmt.Println(top)
			for i, file := range board {
				if i != 0 {
					fmt.Println(mid)
				}
				fmt.Println(strconv.Itoa(8-i) + replacePieces(midPieces, file))
			}
			fmt.Println(bottom)
			fmt.Println(files)
		case "fen":
			if len(fields) > 1 && fields[1] == "extended" {
				fmt.Println(ioutil.ReadAll(chess.XFENReader(game)))
			} else {
				fmt.Println(ioutil.ReadAll(chess.FENReader(game)))
			}
		default:
			fmt.Println("available commands:")
			fmt.Println("move <from> <to>")
			fmt.Println("print")
		}
		fmt.Print("> ")
	}
}

func rotate(board [8][8]chess.Piece) [8][8]chess.Piece {
	var newBoard [8][8]chess.Piece
	for file, row := range board {
		for rank, piece := range row {
			newBoard[7-rank][file] = piece
		}
	}
	return newBoard
}

func replacePieces(format string, pieces [8]chess.Piece) string {
	final := format
	for _, piece := range pieces {
		final = strings.Replace(final, " X ", fmt.Sprintf(" %c ", piece.Symbol()), 1)
	}
	return final
}

func parseMove(g *chess.Game, from, to string) chess.Move {

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
