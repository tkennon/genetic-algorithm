package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/tkennon/genetic-algorithm/game"
	"github.com/tkennon/genetic-algorithm/stats"
)

// TODO(tk):
// - fix problem of slow writes at high generation count.

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

func evolve(breeder game.Breeder) error {
	gen := game.FirstGeneration(generationSize, breeder)
	alphas := stats.New(maxMoves*maxGames, pickUpReward, pickUpPenalty, bumpPenalty)
	for i := 0; i < maxGenerations; i++ {
		var wg sync.WaitGroup
		for _, p := range gen.Players {
			wg.Add(1)
			go func(p *game.Player) {
				defer wg.Done()
				for j := 0; j < maxGames; j++ {
					g := game.New(gridSize, chanceOfRubbish, pickUpReward, pickUpPenalty, bumpPenalty)
					g.Play(maxMoves, p)
				}
			}(p)
		}
		wg.Wait()
		alpha := gen.GetAlpha()
		alphas.Add(alpha)
		log.Printf("Finished generation %d: alpha %s\n", i, alpha)
		gen = gen.NextGeneration(numberOfParents, chanceOfMutation)
	}

	f, err := os.Create(fmt.Sprintf("%s-%s-alpha-counters.png", breeder.Name(), outputFile))
	if err != nil {
		return err
	}
	err = alphas.GetCountersChart(f)
	cerr := f.Close()
	if err != nil {
		return err
	}
	return cerr
}

func parseFlags() error {
	flag.IntVar(&pickUpReward, "pick-up-reward", 2, "the reward for picking up rubbish")
	flag.IntVar(&pickUpPenalty, "pick-up-penalty", 1, "the penalty for picking up rubbish where there is none")
	flag.IntVar(&bumpPenalty, "wall-bump-penalty", 1, "the penalty for bumping into a wall")
	flag.IntVar(&maxMoves, "max-moves", 200, "the number of moves in a game")
	flag.Float64Var(&chanceOfRubbish, "chance-of-rubbish", 0.25, "the chance of any given cell being initialised with rubbish")
	flag.Float64Var(&chanceOfMutation, "chance-of-mutation", 0.01, "the chance of genetic mutation occurring for a given gene")
	flag.IntVar(&gridSize, "grid-size", 10, "the size of one side of the square grid (not including walls)")
	flag.IntVar(&generationSize, "generation-size", 200, "the number of Robbies per generation")
	flag.IntVar(&numberOfParents, "parents", 2, "the number of parents required to create an offspring")
	flag.IntVar(&maxGenerations, "max-generations", 300, "the maximum number of generations to evolve over")
	flag.IntVar(&maxGames, "max-games", 50, "the maximum number of games each Robbie will play")
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
	if err := evolve(game.NewRobbieBreeder()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	gen := game.FirstGeneration(generationSize, game.NewHumanBreeder())
	var wg sync.WaitGroup
	for _, p := range gen.Players {
		wg.Add(1)
		go func(p *game.Player) {
			defer wg.Done()
			for j := 0; j < maxGames; j++ {
				g := game.New(gridSize, chanceOfRubbish, pickUpReward, pickUpPenalty, bumpPenalty)
				g.Play(maxMoves, p)
			}
		}(p)
	}
	wg.Wait()
	fmt.Printf("For comparison, a human did: %s\n", gen.GetAlpha())
	log.Println("took", time.Since(start))
}
