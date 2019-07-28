package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"strconv"
	"strings"

	a "github.com/logrusorgru/aurora"
	"golang.org/x/xerrors"

	"github.com/deanveloper/chess"
	"github.com/deanveloper/chess/chesscodec"
)

var blackAuto [][]string
var whiteAuto [][]string

func main() {
	game := &chess.Game{}
	game.InitClassic()

	scan := bufio.NewScanner(os.Stdin)

	fmt.Print("> ")
	for scan.Scan() {
		runCmd(game, strings.Fields(strings.TrimSpace(scan.Text())))
		fmt.Print("> ")
	}
}

// TODO modularize this
func runCmd(game *chess.Game, fields []string) bool {

	defer func() {
		if rec := recover(); rec != nil {
			fmt.Printf("panic: %v\n", rec)
			debug.PrintStack()
		}
	}()

	if len(fields) < 1 {
		fields = []string{"help"}
	}

	switch fields[0] {
	case "auto":
		if len(fields) < 2 {
			fmt.Println("command auto:")
			fmt.Println("\truns the command immediately AND at the start of the current player's every turn")
			fmt.Println("\tsyntax: auto [command]")
			fmt.Println("\tex: `auto board` (automatically print board before your turn)")
			fmt.Println("\tex: `auto stockfish move` (automatically have stockfish move on this turn)")
			return false
		}
		var turn = game.Turn()
		var autos *[][]string
		if turn == chess.Black {
			autos = &blackAuto
		} else {
			autos = &whiteAuto
		}

		if runCmd(game, fields[1:]) {
			*autos = append(*autos, fields[1:])
			fmt.Printf("%v will now run before %v plays", strings.Join(fields[1:], " "), turn)
		}

	case "move":
		if len(fields) < 2 {
			fmt.Println("command move:")
			fmt.Println("\tmakes a move using uci notation")
			fmt.Println("\tsyntax: move <from><to>[promotion]")
			fmt.Println("\tex: `move e2e4` or `move a7a8q`")
			return false
		}
		move, err := parseMove(game, fields[1])
		if err != nil {
			fmt.Println("error:", err.Error())
			return false
		}
		err = game.MakeMove(move)
		if err != nil {
			fmt.Println("error:", err)
			return false
		}

		var cmds [][]string
		if game.Turn().Other() == chess.Black {
			cmds = blackAuto
		} else {
			cmds = whiteAuto
		}

		// run auto commands for next player
		for _, cmd := range cmds {
			runCmd(game, cmd)
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
		fmt.Println("`print` deprecated, renamed to `board`")
		fallthrough
	case "board":
		board := game.BoardRankFile()

		var rotated bool
		if game.Turn() == chess.Black {
			board = rotate(board)
			rotated = true
		}

		for rank := 7; rank >= 0; rank-- {
			rankSlice := board[rank]
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
			if rotated {
				fmt.Printf(" %d ", 8-rank)
			} else {
				fmt.Printf(" %d ", rank+1)
			}

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
		if rotated {
			fmt.Println("     h    g    f    e    d    c    b    a  ")
		} else {
			fmt.Println("     a    b    c    d    e    f    g    h  ")
		}
	case "fen":
		all, err := ioutil.ReadAll(chesscodec.FENReader(game))
		if err != nil {
			fmt.Println("error:", err)
			return false
		}
		fmt.Println(string(all))
	case "stockfish":
		difficulty := 20
		if len(fields) >= 3 {
			diff, err := strconv.Atoi(fields[2])
			if err != nil || diff < 1 || diff > 20 {
				fmt.Println("difficulty must be a number between 1 and 20")
				return false
			}
			difficulty = diff
		}

		fen, err := ioutil.ReadAll(chess.FENReader(game))
		if err != nil {
			fmt.Println("error:", err)
			return false
		}

		fmt.Println("running stockfish...")
		sfSuggest, err := runStockfish(string(fen), difficulty)

		if len(fields) >= 2 && fields[1] == "move" {
			fmt.Println("stockfish plays " + sfSuggest)
			runCmd(game, []string{"move", sfSuggest})
		} else {
			fmt.Println("stockfish suggests: " + sfSuggest)
		}

	default:
		fmt.Printf("unknown command: %q\n", fields)
		fmt.Println("available commands:")
		fmt.Println("syntax: auto [command]")
		fmt.Println("\truns the command immediately AND at the start of the current player's every turn")
		fmt.Println("\tex: `auto board` (automatically print board before your turn)")
		fmt.Println("\tex: `auto stockfish move` (automatically have stockfish move on this turn)")
		fmt.Println()
		fmt.Println("move <from><to>[promotion]")
		fmt.Println("\tmakes a move using uci notation")
		fmt.Println("\tex: `move e2e4` or `move a7a8q`")
		fmt.Println()
		fmt.Println("pieces")
		fmt.Println("\tlists remaining pieces in the game")
		fmt.Println()
		fmt.Println("fen")
		fmt.Println("\toutputs the game state in FEN notation")
		fmt.Println()
		fmt.Println("board")
		fmt.Println("\toutputs the game on a human-readable board")
		fmt.Println()
		fmt.Println("stockfish [move [difficulty=20]]")
		fmt.Println("\thas stockfish suggest a move. if `move` is")
		fmt.Println("\tset, stockfish will make the move as well")
		fmt.Println("\tit is the current color's turn.")
		fmt.Println("\tex: `stockfish` (suggest a move)")
		fmt.Println("\tex: `stockfish move 10` (play against stockfish level 10)")
		fmt.Println()
		return false
	}
	return true
}

func rotate(board [8][8]chess.Piece) [8][8]chess.Piece {
	var newBoard [8][8]chess.Piece
	for i, rank := range board {
		for j := range rank {
			newBoard[7-i][7-j] = board[i][j]
		}
	}
	return newBoard
}

func parseMove(g *chess.Game, uci string) (chess.Move, error) {

	if len(uci) < 4 {
		return chess.Move{}, xerrors.New("malformed move")
	}

	var from, to string
	var promotion byte

	if len(uci) >= 4 {
		from = uci[0:2]
		to = uci[2:4]
	}
	if len(uci) == 5 {
		promotion = uci[4]
	}

	fromFile := int(from[0] - 'a')
	toFile := int(to[0] - 'a')

	fromRank := int(from[1] - '1')
	toRank := int(to[1] - '1')

	piece, ok := g.PieceAt(chess.Space{File: fromFile, Rank: fromRank})
	if !ok {
		return chess.Move{}, xerrors.New("no piece at " + from)
	}

	pieces := map[byte]chess.PieceType{
		'r': chess.PieceRook,
		'n': chess.PieceKnight,
		'b': chess.PieceBishop,
		'q': chess.PieceQueen,
	}

	var move chess.Move
	move.Moving = piece
	move.To = chess.Space{File: toFile, Rank: toRank}

	if typ, ok := pieces[promotion]; ok {
		move.Promotion = typ
	} else if promotion != 0 {
		return chess.Move{}, xerrors.New("promotion")
	}

	return move, nil
}
