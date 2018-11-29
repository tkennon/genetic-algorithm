package robbie

import (
	"math/rand"
	"sort"
)

// Gene is an instruction that Robbie will follow.
type Gene int

// The basic moves that Robbie can take each turn.
const (
	DoNothing Gene = iota
	MoveUp
	MoveRight
	MoveDown
	MoveLeft
	MoveRandom
	PickUpRubbish
)

func getRandomGene() Gene {
	genes := []Gene{
		DoNothing,
		MoveUp,
		MoveRight,
		MoveDown,
		MoveLeft,
		MoveRandom,
		PickUpRubbish,
	}
	return genes[rand.Intn(len(genes))]
}

// Robbie represents a simple robot that is trying to clean the floor. Robbie
// follows instructions encoded in its genome. Robbie is rewarded for good
// behaviour (picking up rubbish) and punished for bad behaviour (bumping into
// walls and trying to tidy an already clean cell).
type Robbie struct {
	Genome        []Gene
	Score         int
	PickUps       int
	FalsePickUps  int
	Bumps         int
	MissedRubbish int
	TotalRubbish  int
}

// newRobbie returns a new robbie which is the product of it's parents and a
// small amount of random genetic mutation. Note that any number of parents may
// be specified.
func newRobbie(genomeSize int, chanceOfMutation float64, parents ...*Robbie) *Robbie {
	genome := make([]Gene, genomeSize)
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
				genome[i] = parent.Genome[i]
			}
		}
	}
	return &Robbie{Genome: genome}
}

// GetGene returns a gene at the nth position in the genome.
func (r *Robbie) GetGene(n int) Gene {
	//log.Println(r)
	return r.Genome[n]
}

// PickUp increases Robbie's score by reward.
func (r *Robbie) PickUp(reward int) {
	r.Score += reward
	r.PickUps++
}

// FalsePickUp decreases Robbie's score by penalty.
func (r *Robbie) FalsePickUp(penalty int) {
	r.Score -= penalty
	r.FalsePickUps++
}

// Bump decreases Robbie's score by penalty.
func (r *Robbie) Bump(penalty int) {
	r.Score -= penalty
	r.Bumps++
}

// Missed adds the rubbish to the tally of rubbish Robbie has missed.
func (r *Robbie) Missed(rubbish int) {
	r.MissedRubbish += rubbish
}

// Total adds the rubbish to the tall of total rubbish put down in the game.
func (r *Robbie) Total(rubbish int) {
	r.TotalRubbish += rubbish
}

// Generation represents a group of Robbie's with common parentage.
type Generation struct {
	Robbies []*Robbie
}

func newGeneration(size int) *Generation {
	return &Generation{Robbies: make([]*Robbie, size)}
}

// FirstGeneration returns the first generation of Robbies, from which all other
// generations are bred.
func FirstGeneration(generationSize, genomeSize int) *Generation {
	g := newGeneration(generationSize)
	for i := range g.Robbies {
		g.Robbies[i] = newRobbie(genomeSize, 0.0)
	}
	return g
}

// NextGeneration creates a new generation of Robbies from the current
// generation.
func (g *Generation) NextGeneration(numParents int, chanceOfMutation float64) *Generation {
	next := newGeneration(len(g.Robbies))
	parents := g.chooseParents(numParents)
	for i := range next.Robbies {
		next.Robbies[i] = newRobbie(len(g.Robbies[0].Genome), chanceOfMutation, parents...)
	}
	return next
}

// GetAlpha returns the best-performing Robbie of the generation.
func (g *Generation) GetAlpha() *Robbie {
	g.rank()
	return g.Robbies[0]
}

// GetRunt returns the worst-performing Robbie of the generation.
func (g *Generation) GetRunt() *Robbie {
	g.rank()
	return g.Robbies[len(g.Robbies)-1]
}

// chooseParents returns a slice of the top numParents in the given generation.
func (g *Generation) chooseParents(numParents int) []*Robbie {
	g.rank()
	min := numParents
	if len(g.Robbies) < numParents {
		min = len(g.Robbies)
	}
	return g.Robbies[:min]
}

// ranksorts the robbiesinto decending order of score.
func (g *Generation) rank() {
	lessFn := func(i, j int) bool {
		return g.Robbies[i].Score > g.Robbies[j].Score
	}
	sort.SliceStable(g.Robbies, lessFn)
}
