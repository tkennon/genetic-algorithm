package game

import (
	"math/rand"
)

const robbieGenomeSize = 243 // 3^5.

// Robbie represents a simple robot that is trying to clean the floor. Robbie
// follows instructions encoded in its genome.
type Robbie struct {
	Genome []Move
}

type RobbieBreeder struct{}

func NewRobbieBreeder() Breeder {
	return RobbieBreeder{}
}

func (RobbieBreeder) Name() string {
	return "robbie"
}

// NewRobbie returns a new robbie which is the product of it's parents and a
// small amount of random genetic mutation. Note that any number of parents may
// be specified.
func (RobbieBreeder) NewPlayer(chanceOfMutation float64, parents ...*Player) *Player {
	genome := make([]Move, robbieGenomeSize)
	if len(parents) == 0 {
		// No parents given: create a brand new, totally random genome.
		for i := range genome {
			genome[i] = GetRandomMove()
		}
	} else {
		n := len(parents)
		for i := range genome {
			if rand.Float64() < chanceOfMutation {
				// Random genetic mutation.
				genome[i] = GetRandomMove()
			} else {
				// Pick a parent and copy the gene at the current location.
				parent := parents[rand.Intn(n)]
				r := parent.Mover.(*Robbie)
				genome[i] = r.Genome[i]
			}
		}
	}
	return &Player{
		Mover: &Robbie{Genome: genome},
	}
}

// NextMove returns the next move that Robbie will take given the current
// position in the game.
func (r *Robbie) NextMove(pos Position) Move {
	idx := pos.Current + 3*pos.Above + 3*3*pos.Right + 3*3*3*pos.Below + 3*3*3*3*pos.Left
	return r.Genome[idx]
}
