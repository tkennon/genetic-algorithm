package main

import (
	"flag"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"

	"me/ga-robbie/game"
	"me/ga-robbie/robbie"
	"me/ga-robbie/stats"
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
	outputFile       string
)

func init() {
	rand.Seed(time.Now().Unix())
}

func evolve() {
	gen := robbie.FirstGeneration(generationSize, int(math.Pow(3, 5)))
	alphas := stats.New(maxMoves*maxGames, pickUpReward, pickUpPenalty, bumpPenalty)
	runts := stats.New(maxMoves*maxGames, pickUpReward, pickUpPenalty, bumpPenalty)
	for i := 0; i < maxGenerations; i++ {
		var wg sync.WaitGroup
		for _, r := range gen.Robbies {
			wg.Add(1)
			go func(r *robbie.Robbie) {
				defer wg.Done()
				for j := 0; j < maxGames; j++ {
					g := game.New(gridSize, chanceOfRubbish, pickUpReward, pickUpPenalty, bumpPenalty)
					g.Play(maxMoves, r)
				}
			}(r)
		}
		wg.Wait()
		alpha := gen.GetAlpha()
		runt := gen.GetRunt()
		alphas.Add(alpha)
		runts.Add(runt)
		log.Printf("Finished generation %d: alpha %d, runt %d\n", i, alpha.Score, runt.Score)
		gen = gen.NextGeneration(numberOfParents, chanceOfMutation)
	}

	// if err := alphas.GetScoreChart(fmt.Sprintf("%s-alpha-scores.png", outputFile)); err != nil {
	// 	log.Println(err)
	// }
	if err := alphas.GetCountersChart(fmt.Sprintf("%s-alpha-counters.png", outputFile)); err != nil {
		log.Println(err)
	}
	if err := alphas.GetGenomeChart(fmt.Sprintf("%s-alpha-genome.png", outputFile)); err != nil {
		log.Println(err)
	}
	// if err := runts.GetScoreChart(fmt.Sprintf("%s-runt-scores.png", outputFile)); err != nil {
	// 	log.Println(err)
	// }
	if err := runts.GetCountersChart(fmt.Sprintf("%s-runt-counters.png", outputFile)); err != nil {
		log.Println(err)
	}
	if err := runts.GetGenomeChart(fmt.Sprintf("%s-runt-genome.png", outputFile)); err != nil {
		log.Println(err)
	}
}

func parseFlags() error {
	flag.IntVar(&pickUpReward, "pick-up-reward", 2, "the reward for picking up rubbish")
	flag.IntVar(&pickUpPenalty, "pick-up-penalty", 1, "the penalty for picking up rubbish where there is none")
	flag.IntVar(&bumpPenalty, "wall-bump-penalty", 1, "the penalty for bumping into a wall")
	flag.IntVar(&maxMoves, "max-moves", 500, "the number of moves in a game")
	flag.Float64Var(&chanceOfRubbish, "chance-of-rubbish", 0.25, "the chance of any given cell being initialised with rubbish")
	flag.Float64Var(&chanceOfMutation, "chance-of-mutation", 0.01, "the chance of genetic mutation occurring for a given gene")
	flag.IntVar(&gridSize, "grid-size", 10, "the size of one side of the square gird (not including walls)")
	flag.IntVar(&generationSize, "generation-size", 200, "the number of Robbies per generation")
	flag.IntVar(&numberOfParents, "parents", 2, "the number of parents required to create an offspring")
	flag.IntVar(&maxGenerations, "max-generations", 300, "the maximum number of generations to evolve over")
	flag.IntVar(&maxGames, "max-games", 100, "the maximum number of games each Robbie will play")
	flag.StringVar(&outputFile, "output-file", "chart", "the name of the resulting file (no extension)")
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
