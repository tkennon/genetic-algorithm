package game

import (
	"math/rand"
)

type Human struct{}

type HumanBreeder struct{}

func NewHumanBreeder() Breeder {
	return HumanBreeder{}
}

func (HumanBreeder) Name() string {
	return "human"
}

func (HumanBreeder) NewPlayer(float64, ...*Player) *Player {
	return &Player{
		Mover: Human{},
	}
}

func (Human) NextMove(pos Position) Move {
	if pos.Current == rubbish {
		return PickUpRubbish
	}
	var moves []Move

	if pos.Above == rubbish {
		moves = append(moves, MoveUp)
	}
	if pos.Right == rubbish {
		moves = append(moves, MoveRight)
	}
	if pos.Below == rubbish {
		moves = append(moves, MoveDown)
	}
	if pos.Left == rubbish {
		moves = append(moves, MoveLeft)
	}
	if l := len(moves); l > 0 {
		m := moves[rand.Intn(l)]
		return m
	}

	if pos.Above == empty {
		moves = append(moves, MoveUp)
	}
	if pos.Right == empty {
		moves = append(moves, MoveRight)
	}
	if pos.Below == empty {
		moves = append(moves, MoveDown)
	}
	if pos.Left == empty {
		moves = append(moves, MoveLeft)
	}
	if l := len(moves); l > 0 {
		m := moves[rand.Intn(l)]
		return m
	}

	return DoNothing
}
