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
)

// the difficulty for each stockfish to run at. -1 if it shouldn't run
var blackStockfish, whiteStockfish = -1, -1

func main() {

	game := &chess.Game{}
	game.Init()

	var lines = make(chan string)
	readLines(lines)

	for {
		fmt.Print("> ")
		var line string
		select {
		case l := <-lines:
			line = l
		}
		runCmd(game, strings.Fields(line))
	}
}

// returns the lines as a channel of strings
func readLines(ch chan<- string) {
	go func() {
		scan := bufio.NewScanner(os.Stdin)
		for scan.Scan() {
			ch <- scan.Text()
		}
	}()
}

func runCmd(game *chess.Game, fields []string) {

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
	case "move":
		if len(fields) < 2 {
			fmt.Println("move syntax (uci):")
			fmt.Println("move <from><to>[promotion]")
			fmt.Println("\tex: `move e2e4` or `move a7a8q`")
		}
		move, err := parseMove(game, fields[1])
		if err != nil {
			fmt.Println("error:", err.Error())
		}
		err = game.MakeMove(move)
		if err != nil {
			fmt.Println("Error making move:", err)
			return
		}

		autoStockfish(game)
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
			fmt.Printf(" %d ", rank)

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
		all, err := ioutil.ReadAll(chess.FENReader(game))
		if err != nil {
			fmt.Println("error:", err)
			return
		}
		fmt.Println(string(all))
	case "stockfish":
		difficulty := 20
		var auto bool
		if len(fields) >= 2 {
			diff, err := strconv.Atoi(fields[1])
			if err != nil {
				fmt.Println("first arg to stockfish must be a number")
				return
			}
			difficulty = diff
		}
		if len(fields) >= 3 {
			auto = (fields[2] == "auto")
		}

		fen, err := ioutil.ReadAll(chess.FENReader(game))
		if err != nil {
			fmt.Println("error:", err)
			return
		}

		fmt.Println("running stockfish...")
		sfSuggest, err := runStockfish(string(fen), difficulty)

		if auto {
			if len(game.History)%2 == 0 {
				whiteStockfish = difficulty
			} else {
				blackStockfish = difficulty
			}
			fmt.Println("stockfish plays " + sfSuggest)
			runCmd(game, []string{"move", sfSuggest})
		} else {
			fmt.Println("stockfish suggests: " + sfSuggest)
		}

	default:
		fmt.Println("available commands:")
		fmt.Println("move <from><to>[promotion]")
		fmt.Println("\tmakes a move in the game")
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
		fmt.Println("stockfish [ELO=3000] [auto]")
		fmt.Println("\thas stockfish suggest a move. if `auto` is")
		fmt.Println("\tset, stockfish will automatically run each time")
		fmt.Println("\tit is the current color's turn.")
		fmt.Println("\tex: `stockfish 20` (suggest a move)")
		fmt.Println("\tex: `stockfish 10 auto` (play against stockfish)")
		fmt.Println()
	}
}

func autoStockfish(game *chess.Game) {
	difficulty := -1
	if len(game.History)%2 == 0 {
		difficulty = whiteStockfish
	} else {
		difficulty = blackStockfish
	}

	if difficulty == -1 {
		return
	}

	fmt.Println("running stockfish...")

	fen, err := ioutil.ReadAll(chess.FENReader(game))
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	moveS, err := runStockfish(string(fen), difficulty)
	if err != nil {
		fmt.Println("error running stockfish:", err)
		return
	}

	fmt.Println("stockfish plays " + moveS)

	runCmd(game, []string{"move", moveS})
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
		'r': chess.Rook,
		'n': chess.Knight,
		'b': chess.Bishop,
		'q': chess.Queen,
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
