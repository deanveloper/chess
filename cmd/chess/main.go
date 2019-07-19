package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	a "github.com/logrusorgru/aurora"

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
			board := rotate(game.Board())
			for rank, rankSlice := range board {

				const black, white = 5, 15

				fmt.Print("   ")

				for file := 0; file < 8; file++ {
					space := chess.Space{Rank: rank, File: file}
					if space.Color() == chess.White {
						fmt.Print(a.BgGray(white, "     "))
					} else {
						fmt.Print(a.BgGray(black, "     "))
					}
				}

				fmt.Println()
				fmt.Printf(" %d ", 8-rank)

				for file, piece := range rankSlice {
					space := chess.Space{Rank: rank, File: file}

					var symbol a.Value
					if piece.Color == chess.White {
						symbol = a.White(string(piece.Type.Symbol()))
					} else {
						symbol = a.Black(string(piece.Type.Symbol()))
					}

					if space.Color() == chess.White {
						fmt.Print(a.Sprintf(a.BgGray(white, "  %s  "), a.BgGray(white, symbol)))
					} else {
						fmt.Print(a.Sprintf(a.BgGray(black, "  %s  "), a.BgGray(black, symbol)))
					}
				}

				fmt.Println()
				fmt.Print("   ")

				for file := 0; file < 8; file++ {
					space := chess.Space{Rank: rank, File: file}
					if space.Color() == chess.White {
						fmt.Print(a.BgGray(white, "     "))
					} else {
						fmt.Print(a.BgGray(black, "     "))
					}
				}

				fmt.Println()
			}
			fmt.Println("     a    b    c    d    e    f    g    h  ")
		case "fen":
			if len(fields) > 1 && fields[1] == "extended" {
				fmt.Println(ioutil.ReadAll(chess.XFENReader(game)))
			} else {
				fmt.Println(ioutil.ReadAll(chess.FENReader(game)))
			}
		default:
			fmt.Println("available commands:")
			fmt.Println("move <from> <to>")
			fmt.Println("fen [extended]")
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
