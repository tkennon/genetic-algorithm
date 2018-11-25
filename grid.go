package main

import (
	"math/rand"
)

type cellState int

const (
	empty cellState = iota
	rubbish
	wall
)

type grid struct {
	// The current contents of the grid.
	cells []cellState
	// Robbie's x position in the grid (zero if left-most).
	x int
	// Robbie's y position in the grid (zero is bottom-most).
	y int
}

// newGrid returns a new grid object with the outer-most cells in the grid set
// as walls, and the inner cells randomly populated with rubbish. Robbie's
// initial position is also randomly set.
func newGrid() *grid {
	extendedSize := gridSize + 2
	g := &grid{
		cells: make([]cellState, extendedSize*extendedSize),
		x:     rand.Intn(gridSize) + 1,
		y:     rand.Intn(gridSize) + 1,
	}
	for i := 0; i < extendedSize; i++ {
		g.setCell(i, 0, wall)
		g.setCell(i, extendedSize-1, wall)
		g.setCell(0, i, wall)
		g.setCell(extendedSize-1, i, wall)
	}
	for x := 1; x <= gridSize; x++ {
		for y := 1; y <= gridSize; y++ {
			if rand.Float64() < chanceOfRubbish {
				g.setCell(x, y, rubbish)
			}
		}
	}

	return g
}

func positionToIndex(x, y int) int {
	return x + y*(gridSize+2)
}

func (g *grid) getCell(x, y int) cellState {
	return g.cells[positionToIndex(x, y)]
}

// isCell compares the contents of the specified cell to the given value.
func (g *grid) isCell(x, y int, val cellState) bool {
	return g.cells[positionToIndex(x, y)] == val
}

// setCell sets the contents of the specified cell to the given value.
func (g *grid) setCell(x, y int, val cellState) {
	g.cells[positionToIndex(x, y)] = val
}

// isCurrentCell compares the contents of the cell Robbie is currently at to the
// given value.
func (g *grid) isCurrentCell(val cellState) bool {
	return g.isCell(g.x, g.y, val)
}

// setCurrentCell sets the contents of the cell Robbie is currently at the the
// given value.
func (g *grid) setCurrentCell(val cellState) {
	g.setCell(g.x, g.y, val)
}

// getSituation returns Robbie's current situation in the grid as a number in
// the range [0, genomeSize).
func (g *grid) getSituation() int {
	current := int(g.getCell(g.x, g.y))
	above := int(g.getCell(g.x, g.y+1))
	right := int(g.getCell(g.x+1, g.y))
	below := int(g.getCell(g.x, g.y-1))
	left := int(g.getCell(g.x-1, g.y))
	return current + 3*above + 3*3*right + 3*3*3*below + 3*3*3*3*left
}

// moveUp attempts to move Robbie upwards one cell. It returns the change in
// score as a result of this move.
func (g *grid) moveUp() int {
	if g.y == gridSize {
		return bumpPenalty
	}
	g.y++
	return 0
}

// moveRight attempts to move Robbie right one cell. It returns the change in
// score as a result of this move.
func (g *grid) moveRight() int {
	if g.x == gridSize {
		return bumpPenalty
	}
	g.x++
	return 0
}

// moveDown attempts to move Robbie downwards one cell. It returns the change in
// score as a result of this move.
func (g *grid) moveDown() int {
	if g.y == 1 {
		return bumpPenalty
	}
	g.y--
	return 0
}

// moveLeft attempts to move Robbie left one cell. It returns the change in
// score as a result of this move.
func (g *grid) moveLeft() int {
	if g.x == 1 {
		return bumpPenalty
	}
	g.x--
	return 0
}

// moveRandom attempts to move Robbie in a random direction. It returns the
// change in score as a result of this move.
func (g *grid) moveRandom() int {
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

// pickUp attempts to pick up rubbish from Robbie's current position in the
// grid. It returns the change in score as a result of this move.
func (g *grid) pickUp() int {
	if g.isCurrentCell(rubbish) {
		g.setCurrentCell(empty)
		return pickUpReward
	}
	return pickUpPenalty
}
