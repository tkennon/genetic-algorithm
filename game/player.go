package game

import (
	"fmt"
	"sort"
)

// Player is a player of the game. It decides where to move next and keeps track
// of the ongoing score. Players are rewarded for good behaviour (picking up
// rubbish) and punished for bad behaviour (bumping into walls and trying to
// tidy an already clean cell).
type Player struct {
	Mover         Mover
	Score         int
	PickUps       int
	FalsePickUps  int
	Bumps         int
	MissedRubbish int
	TotalRubbish  int
}

func (p *Player) String() string {
	return fmt.Sprintf("&Player{Score: %d, PickUps: %d, FalsePickups: %d, Bumps: %d, MissedRubbish: %d, TotalRubbish: %d}", p.Score, p.PickUps, p.FalsePickUps, p.Bumps, p.MissedRubbish, p.TotalRubbish)
}

// PickUp increases the Player's score by reward.
func (p *Player) PickUp(reward int) {
	p.Score += reward
	p.PickUps++
}

// FalsePickUp decreases the Player's score by penalty.
func (p *Player) FalsePickUp(penalty int) {
	p.Score -= penalty
	p.FalsePickUps++
}

// Bump decreases the Player's score by penalty.
func (p *Player) Bump(penalty int) {
	p.Score -= penalty
	p.Bumps++
}

type Breeder interface {
	Name() string
	NewPlayer(chanceOfMutation float64, parents ...*Player) *Player
}

// Generation represents a group of Player's with common parentage.
type Generation struct {
	Players []*Player
	Breeder Breeder
}

func newGeneration(size int, breeder Breeder) *Generation {
	return &Generation{
		Players: make([]*Player, size),
		Breeder: breeder,
	}
}

// FirstGeneration returns the first generation of Players, from which all other
// generations are bred.
func FirstGeneration(generationSize int, breeder Breeder) *Generation {
	g := newGeneration(generationSize, breeder)
	for i := range g.Players {
		g.Players[i] = breeder.NewPlayer(0.0)
	}
	return g
}

// NextGeneration creates a new generation of Robbies from the current
// generation.
func (g *Generation) NextGeneration(numParents int, chanceOfMutation float64) *Generation {
	next := newGeneration(len(g.Players), g.Breeder)
	parents := g.chooseParents(numParents)
	for i := range next.Players {
		next.Players[i] = g.Breeder.NewPlayer(chanceOfMutation, parents...)
	}
	return next
}

// GetAlpha returns the best-performing Player of the generation.
func (g *Generation) GetAlpha() *Player {
	g.rank()
	return g.Players[0]
}

// GetRunt returns the worst-performing Player of the generation.
func (g *Generation) GetRunt() *Player {
	g.rank()
	return g.Players[len(g.Players)-1]
}

// chooseParents returns a slice of the top numParents in the given generation.
func (g *Generation) chooseParents(numParents int) []*Player {
	g.rank()
	min := numParents
	if len(g.Players) < numParents {
		min = len(g.Players)
	}
	return g.Players[:min]
}

// rank sorts the Players into decending order of score.
func (g *Generation) rank() {
	lessFn := func(i, j int) bool {
		return g.Players[i].Score > g.Players[j].Score
	}
	sort.SliceStable(g.Players, lessFn)
}
