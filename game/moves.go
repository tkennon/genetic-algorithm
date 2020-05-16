package game

import "math/rand"

// Move is a move a player can make.
type Move int

// The moves that a player can take each turn.
const (
	DoNothing Move = iota
	MoveUp
	MoveRight
	MoveDown
	MoveLeft
	MoveRandom
	PickUpRubbish
)

// GetRandomMove returns a random move from the set of possible moves.
func GetRandomMove() Move {
	moves := []Move{
		DoNothing,
		MoveUp,
		MoveRight,
		MoveDown,
		MoveLeft,
		MoveRandom,
		PickUpRubbish,
	}
	return moves[rand.Intn(len(moves))]
}

// Mover is an object that decides the next move to take given the current
// position in the game.
type Mover interface {
	NextMove(Position) Move
}
