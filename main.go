package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	chart "github.com/wcharczuk/go-chart"
)

// Parameters of the game.
var (
	pickUpReward     int
	pickUpPenalty    int
	bumpPenalty      int
	maxMoves         int
	chanceOfRubbish  float64
	chanceOfMutation float64
	gridSize         int
	maxGames         int
	maxGenerations   int
	generationSize   int
	numberOfParents  int
)

func init() {
	rand.Seed(time.Now().Unix())
}

type statistics struct {
	s             strategy
	missedPickUps int
}

func evolve() {
	currentGen := createNextGeneration(nil)
	alphas := make([]statistics, 0)
	runts := make([]statistics, 0)
	for i := 0; i < maxGenerations; i++ {
		var wg sync.WaitGroup
		for _, s := range currentGen {
			wg.Add(1)
			go func(s *strategy) {
				defer wg.Done()
				for j := 0; j < maxGames; j++ {
					g := newGrid()
					for k := 0; k < maxMoves; k++ {
						switch s.getMove(g.getSituation()) {
						case doNothing:
						case moveUp:
							g.moveUp(s)
						case moveRight:
							g.moveRight(s)
						case moveDown:
							g.moveDown(s)
						case moveLeft:
							g.moveLeft(s)
						case moveRandom:
							g.moveRandom(s)
						case pickUpRubbish:
							g.pickUp(s)
						}
					}
					s.leftovers += g.getLeftoverRubbish()
				}
			}(s)
		}
		wg.Wait()
		alpha := getAlpha(currentGen)
		runt := getRunt(currentGen)
		alphas = append(alphas, statistics{s: *alpha})
		runts = append(runts, statistics{s: *runt})
		log.Printf("Finished generation %d: alpha %d, runt %d\n", i, alpha.score, runt.score)
		currentGen = createNextGeneration(currentGen)
	}

	alphaScores := make([]float64, len(alphas))
	xaxis := make([]float64, len(alphas))
	for i, a := range alphas {
		xaxis[i] = float64(i)
		alphaScores[i] = float64(a.s.score)
	}
	graph := chart.Chart{
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: xaxis,
				YValues: alphaScores,
			},
			// chart.ContinuousSeries{
			// 	XValues: xaxis,
			// 	yValues:
			// }
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	f, err := os.Create("chart.png")
	if err != nil {
		log.Println(err)
		return
	}
	f.Write(buffer.Bytes())
}

func parseFlags() error {
	flag.IntVar(&pickUpReward, "pick-up-reward", 10, "the reward for picking up rubbish")
	flag.IntVar(&pickUpPenalty, "pick-up-penalty", -5, "the penalty for picking up rubbish where there is none")
	flag.IntVar(&bumpPenalty, "wall-bump=penalty", -1, "the penalty for bumping into a wall")
	flag.IntVar(&maxMoves, "max-moves", 500, "the number of moves in a game")
	flag.Float64Var(&chanceOfRubbish, "chance-of-rubbish", 0.25, "the chance of any given cell being initialised with rubbish")
	flag.Float64Var(&chanceOfMutation, "chance-of-mutation", 0.01, "the chance of genetic mutation occurring for a given gene")
	flag.IntVar(&gridSize, "grid-size", 10, "the size of one side of the square gird (not including walls)")
	flag.IntVar(&generationSize, "generation-size", 200, "the numbr of strategies per generation")
	flag.IntVar(&numberOfParents, "parents", 2, "the number of parents required to create an offspring")
	flag.IntVar(&maxGenerations, "max-generations", 500, "the maximum number of generations to evolve over")
	flag.IntVar(&maxGames, "max-games", 100, "the maximum number of games per strategy")
	flag.Parse()
	if chanceOfRubbish < 0.0 || chanceOfRubbish > 1.0 {
		return fmt.Errorf("chance-of-rubbish not within allowed range [0.0, 1.0]: %f", chanceOfRubbish)
	}
	if gridSize < 1 {
		return fmt.Errorf("grid-size must be greater than zero: %d", gridSize)
	}
	if generationSize < 1 {
		return fmt.Errorf("generation-size must be greater than zero: %d", generationSize)
	}
	if numberOfParents < 1 {
		return fmt.Errorf("number of parents must be greater than zero: %d", numberOfParents)
	}
	return nil
}

func main() {
	log.SetFlags(log.LUTC | log.Ldate | log.Lmicroseconds | log.Lshortfile)
	if err := parseFlags(); err != nil {
		log.Println(err)
		flag.Usage()
		os.Exit(1)
	}
	start := time.Now()
	evolve()
	log.Println("took", time.Since(start))
}
