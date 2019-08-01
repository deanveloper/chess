# chess

`chess` is a pure-Go library that allows the simulation of a game of chess. It also has an `encoder` package which allows someone to retrieve a game as a FEN or PGN, with possibly more coming in the future.

The package's documentation can be found on [godoc](https://godoc.org/github.com/deanveloper/chess).

### CLI

The CLI takes a series of commands from standard input.

| command | syntax | description |
| ------- | ------ | ----------- |
| move | `move <algebraic>` | Moves a piece on the board using algebraic notation |
| board | `board` | Prints the current board |
| pieces | `pieces` | Lists the current pieces on the board |
| stockfish | `stockfish ["move" [difficulty (0-20)]]` | Evaluates the best move with stockfish. If `stockfish move` is run, it will make the move as well |
| auto | `auto [cmd]` | Runs the command at the beginning of the player's turn |
| fen | `fen` | Prints the current FEN |
| pgn | `pgn` | Prints the PGN |
