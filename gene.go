package main

import (
	"math"
	"math/rand"
	"sort"
)

type gene int

// The basic moves that Robbie can take each turn.
const (
	doNothing gene = iota
	moveUp
	moveRight
	moveDown
	moveLeft
	moveRandom
	pickUpRubbish
)

var genomeSize = int(math.Pow(3, 5))

type strategy struct {
	// Robbie's instruction set for this game.
	genome []gene
	// Cumulative score over all games.
	score int
	// Cumulative successful pickups over all games.
	pickUps int
	// Cumulative failed pickups over all games.
	falsePickUps int
	// Cumlative wall bumps over all games.
	bumps int
	// Cumulative leftover rubbish.
	leftovers int
}

func (s *strategy) getMove(situation int) gene {
	return s.genome[situation]
}

func (s *strategy) pickUp() {
	s.score += pickUpReward
	s.pickUps++
}

func (s *strategy) falsePickUp() {
	s.score += pickUpPenalty
	s.falsePickUps++
}

func (s *strategy) bump() {
	s.score += bumpPenalty
	s.bumps++
}

func getRandomGene() gene {
	genes := []gene{
		doNothing,
		moveUp,
		moveRight,
		moveDown,
		moveLeft,
		moveRandom,
		pickUpRubbish,
	}
	return genes[rand.Intn(len(genes))]
}

// createChild returns a new genome which is the product of it's parents and a
// small amount of random genetic mutation. Note that any number of parents may
// be specified.
func createChild(parents ...*strategy) *strategy {
	genome := make([]gene, genomeSize)
	if len(parents) == 0 {
		// No parents given: create a brand new, totally random genome.
		for i := range genome {
			genome[i] = getRandomGene()
		}
	} else {
		n := len(parents)
		for i := range genome {
			if rand.Float64() < chanceOfMutation {
				// Random genetic mutation.
				genome[i] = getRandomGene()
			} else {
				// Pick a parent and copy the gene at the current location.
				parent := parents[rand.Intn(n)]
				genome[i] = parent.genome[i]
			}
		}
	}
	return &strategy{
		genome:       genome,
		score:        0,
		pickUps:      0,
		falsePickUps: 0,
		bumps:        0,
		leftovers:    0,
	}
}

// createNextgeneration creates the next generation of strategies from the
// current generation. To create a brand new generation, pass in nil.
func createNextGeneration(current []*strategy) []*strategy {
	nextGen := make([]*strategy, 0)
	parents := chooseParents(current)
	for i := 0; i < generationSize; i++ {
		nextGen = append(nextGen, createChild(parents...))
	}
	return nextGen
}

// chooseParents returns a slice of the top numberOfParents in the given
// generation
func chooseParents(generation []*strategy) []*strategy {
	rankGeneration(generation)
	min := numberOfParents
	if len(generation) < numberOfParents {
		min = len(generation)
	}
	return generation[:min]
}

// getAlpha returns the best-performing strategy of the generation.
func getAlpha(generation []*strategy) *strategy {
	rankGeneration(generation)
	return generation[0]
}

// getRunt returns the worst-performing strategy of the generation.
func getRunt(generation []*strategy) *strategy {
	rankGeneration(generation)
	return generation[len(generation)-1]
}

// rankGeneration sorts the strategies into decending order of score.
func rankGeneration(generation []*strategy) {
	lessFn := func(i, j int) bool {
		return generation[i].score > generation[j].score
	}
	sort.SliceStable(generation, lessFn)
}
