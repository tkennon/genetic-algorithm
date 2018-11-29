package game

import (
	"math/rand"

	"me/ga-robbie/robbie"
)

type cellState int

const (
	empty cellState = iota
	rubbish
	wall
)

// Game is an object that will test how good Robbie is.
type Game interface {
	Play(numTurns int, r *robbie.Robbie)
}

type grid struct {
	// The current contents of the grid.
	cells []cellState
	// The length of one side of the square grid (including both walls).
	size int
	// Robbie's x position in the grid (zero is left-most).
	x int
	// Robbie's y position in the grid (zero is bottom-most).
	y int
	// How much Robbie is rewarded for successfully picking up rubbish.
	pickUpReward int
	// How much Robbie is punished for mistakenly trying to pick up rubbish.
	pickUpPenalty int
	// How much Robbie is punished for bumping into the wall of the grid.
	bumpPenalty int
}

// New returns a new game ready for Robbie to play. The game constists of a
// size*size grid with the cells randomly populated with rubbish. Robbie is
// rewarded with pickUpReward every time it picks up rubbish in a cell,
// punished with pickUpPenalty every time it tries to pick up rubbish in an
// empty cell, and punsihed with bumpPenalty every time it bumps into the wall
// of the grid.
func New(size int, chanceOfRubbish float64, pickUpReward, pickUpPenalty, bumpPenalty int) Game {
	extendedSize := size + 2
	g := &grid{
		cells:         make([]cellState, extendedSize*extendedSize),
		size:          extendedSize,
		x:             rand.Intn(size) + 1,
		y:             rand.Intn(size) + 1,
		pickUpReward:  pickUpReward,
		pickUpPenalty: pickUpPenalty,
		bumpPenalty:   bumpPenalty,
	}
	for i := 0; i < extendedSize; i++ {
		g.setCell(i, 0, wall)
		g.setCell(i, extendedSize-1, wall)
		g.setCell(0, i, wall)
		g.setCell(extendedSize-1, i, wall)
	}
	for x := 1; x <= size; x++ {
		for y := 1; y <= size; y++ {
			if rand.Float64() < chanceOfRubbish {
				g.setCell(x, y, rubbish)
			}
		}
	}

	return g
}

func (g *grid) Play(turns int, r *robbie.Robbie) {
	r.Total(g.countRubbish())
	for i := 0; i < turns; i++ {
		switch r.GetGene(g.getSituation()) {
		case robbie.DoNothing:
		case robbie.MoveUp:
			if !g.moveUp() {
				r.Bump(g.bumpPenalty)
			}
		case robbie.MoveRight:
			if !g.moveRight() {
				r.Bump(g.bumpPenalty)
			}
		case robbie.MoveDown:
			if !g.moveDown() {
				r.Bump(g.bumpPenalty)
			}
		case robbie.MoveLeft:
			if !g.moveLeft() {
				r.Bump(g.bumpPenalty)
			}
		case robbie.MoveRandom:
			if !g.moveRandom() {
				r.Bump(g.bumpPenalty)
			}
		case robbie.PickUpRubbish:
			if g.pickUp() {
				r.PickUp(g.pickUpReward)
			} else {
				r.FalsePickUp(g.pickUpPenalty)
			}
		}
	}
	r.Missed(g.countRubbish())
}

func (g *grid) positionToIndex(x, y int) int {
	return x + y*(g.size)
}

func (g *grid) getCell(x, y int) cellState {
	return g.cells[g.positionToIndex(x, y)]
}

func (g *grid) isCell(x, y int, val cellState) bool {
	return g.cells[g.positionToIndex(x, y)] == val
}

func (g *grid) setCell(x, y int, val cellState) {
	g.cells[g.positionToIndex(x, y)] = val
}

func (g *grid) isCurrentCell(val cellState) bool {
	return g.isCell(g.x, g.y, val)
}

func (g *grid) setCurrentCell(val cellState) {
	g.setCell(g.x, g.y, val)
}

func (g *grid) getSituation() int {
	current := int(g.getCell(g.x, g.y))
	above := int(g.getCell(g.x, g.y+1))
	right := int(g.getCell(g.x+1, g.y))
	below := int(g.getCell(g.x, g.y-1))
	left := int(g.getCell(g.x-1, g.y))
	return current + 3*above + 3*3*right + 3*3*3*below + 3*3*3*3*left
}

func (g *grid) moveUp() bool {
	if g.y == g.size-2 {
		return false
	}
	g.y++
	return true
}

func (g *grid) moveRight() bool {
	if g.x == g.size-2 {
		return false
	}
	g.x++
	return true
}

func (g *grid) moveDown() bool {
	if g.y == 1 {
		return false
	}
	g.y--
	return true
}

func (g *grid) moveLeft() bool {
	if g.x == 1 {
		return false
	}
	g.x--
	return true
}

func (g *grid) moveRandom() bool {
	switch rand.Intn(4) {
	case 0:
		return g.moveUp()
	case 1:
		return g.moveRight()
	case 2:
		return g.moveDown()
	case 3:
		return g.moveLeft()
	default:
		panic("out of range")
	}
}

func (g *grid) pickUp() bool {
	if g.isCurrentCell(rubbish) {
		g.setCurrentCell(empty)
		return true
	}
	return false
}

func (g *grid) countRubbish() int {
	count := 0
	for _, c := range g.cells {
		if c == rubbish {
			count++
		}
	}
	return count
}
