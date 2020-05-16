package game

import (
	"math/rand"
)

type cellState int

const (
	empty cellState = iota
	rubbish
	wall
)

// Game is an object that will test how good a player is.
type Game interface {
	Play(numTurns int, p *Player)
}

type grid struct {
	// The current contents of the grid.
	cells []cellState
	// The length of one side of the square grid (including both walls).
	size int
	// The Player's x position in the grid (zero is left-most).
	x int
	// The Player's y position in the grid (zero is bottom-most).
	y int
	// How much the Player is rewarded for successfully picking up rubbish.
	pickUpReward int
	// How much the Player is punished for mistakenly trying to pick up rubbish.
	pickUpPenalty int
	// How much the Player is punished for bumping into the wall of the grid.
	bumpPenalty int
}

// New returns a new game ready for the Player to play. The game constists of a
// size*size grid with the cells randomly populated with rubbish. The Player is
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

func (g *grid) Play(turns int, p *Player) {
	p.TotalRubbish += g.countRubbish()
	for i := 0; i < turns; i++ {
		switch p.Mover.NextMove(g.GetCurrentPosition()) {
		case DoNothing:
		case MoveUp:
			if !g.moveUp() {
				p.Bump(g.bumpPenalty)
			}
		case MoveRight:
			if !g.moveRight() {
				p.Bump(g.bumpPenalty)
			}
		case MoveDown:
			if !g.moveDown() {
				p.Bump(g.bumpPenalty)
			}
		case MoveLeft:
			if !g.moveLeft() {
				p.Bump(g.bumpPenalty)
			}
		case MoveRandom:
			if !g.moveRandom() {
				p.Bump(g.bumpPenalty)
			}
		case PickUpRubbish:
			if g.pickUp() {
				p.PickUp(g.pickUpReward)
			} else {
				p.FalsePickUp(g.pickUpPenalty)
			}
		}
	}
	p.MissedRubbish += g.countRubbish()
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

// Situation encodes the situation that player is in on a given turn in the game.
// For example: current empty, above wall, left rubbish etc.
type Position struct {
	Current cellState
	Above   cellState
	Right   cellState
	Below   cellState
	Left    cellState
}

func (g *grid) GetCurrentPosition() Position {
	return Position{
		Current: g.getCell(g.x, g.y),
		Above:   g.getCell(g.x, g.y+1),
		Right:   g.getCell(g.x+1, g.y),
		Below:   g.getCell(g.x, g.y-1),
		Left:    g.getCell(g.x-1, g.y),
	}
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
