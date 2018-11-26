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

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/seq"
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

type statistics struct {
	scores       []float64
	pickUps      []float64
	falsePickUps []float64
	bumps        []float64
	leftovers    []float64
	numGenes     map[gene][]float64
}

func newStatistics() *statistics {
	return &statistics{
		scores:       make([]float64, 0),
		pickUps:      make([]float64, 0),
		falsePickUps: make([]float64, 0),
		bumps:        make([]float64, 0),
		leftovers:    make([]float64, 0),
		numGenes:     make(map[gene][]float64, 0),
	}
}

func (stats *statistics) add(s *strategy) {
	stats.scores = append(stats.scores, float64(s.score))
	stats.pickUps = append(stats.pickUps, float64(s.pickUps))
	stats.falsePickUps = append(stats.falsePickUps, float64(s.falsePickUps))
	stats.bumps = append(stats.bumps, float64(s.bumps))
	stats.leftovers = append(stats.leftovers, float64(s.leftovers))
	counters := make(map[gene]float64)
	for _, g := range s.genome {
		if _, ok := counters[g]; !ok {
			counters[g] = 1
		} else {
			counters[g]++
		}
	}
	for k, v := range counters {
		stats.numGenes[k] = append(stats.numGenes[k], v)
	}
}

func (stats *statistics) getScoreChart() chart.Chart {
	xaxis := seq.Range(0.0, float64(maxGenerations-1))
	score := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 120,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "score",
				XValues: xaxis,
				YValues: stats.scores,
			},
			chart.ContinuousSeries{
				Name:    "missed-rubbish",
				XValues: xaxis,
				YValues: stats.leftovers,
			},
		},
	}
	score.Elements = []chart.Renderable{
		chart.LegendLeft(&score),
	}

	return score
}

func (stats *statistics) getGenomeChart() chart.Chart {
	xaxis := seq.Range(0.0, float64(maxGenerations-1))
	doNothings := make([]float64, maxGenerations)
	pickUps := make([]float64, maxGenerations)
	moveUps := make([]float64, maxGenerations)
	moveRights := make([]float64, maxGenerations)
	moveDowns := make([]float64, maxGenerations)
	moveLefts := make([]float64, maxGenerations)
	moveRandoms := make([]float64, maxGenerations)
	for i := 0; i < maxGenerations; i++ {
		doNothings[i] = stats.numGenes[doNothing][i]
		pickUps[i] = stats.numGenes[pickUpRubbish][i] + doNothings[i]
		moveUps[i] = stats.numGenes[moveUp][i] + pickUps[i]
		moveRights[i] = stats.numGenes[moveRight][i] + moveUps[i]
		moveDowns[i] = stats.numGenes[moveDown][i] + moveRights[i]
		moveLefts[i] = stats.numGenes[moveLeft][i] + moveDowns[i]
		moveRandoms[i] = stats.numGenes[moveRandom][i] + moveLefts[i]
	}
	genomes := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 120,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
				XValues: xaxis,
				YValues: doNothings,
				Name:    "do-nothing",
			},
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					FillColor:   chart.GetDefaultColor(2).WithAlpha(0),
					StrokeColor: chart.GetDefaultColor(2).WithAlpha(0),
				},
				XValues: xaxis,
				YValues: moveUps,
				Name:    "move-up",
			},
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(8).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(8).WithAlpha(64),
				},
				XValues: xaxis,
				YValues: moveRights,
				Name:    "move-right",
			},
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(12).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(12).WithAlpha(64),
				},
				XValues: xaxis,
				YValues: moveDowns,
				Name:    "move-down",
			},
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(16).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(16).WithAlpha(64),
				},
				XValues: xaxis,
				YValues: moveLefts,
				Name:    "move-left",
			},
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(20).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(20).WithAlpha(64),
				},
				XValues: xaxis,
				YValues: moveRandoms,
				Name:    "move-random",
			},
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(24).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(24).WithAlpha(64),
				},
				XValues: xaxis,
				YValues: pickUps,
				Name:    "pick-up-rubbish",
			},
		},
	}
	genomes.Elements = []chart.Renderable{
		chart.LegendLeft(&genomes),
	}

	return genomes
}

func (stats *statistics) getCountersChart() chart.Chart {
	xaxis := seq.Range(0.0, float64(maxGenerations-1))
	counters := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.StyleShow(),
		},
		YAxis: chart.YAxis{
			Style: chart.StyleShow(),
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top:  20,
				Left: 120,
			},
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				Name:    "pick-ups",
				XValues: xaxis,
				YValues: stats.pickUps,
			},
			chart.ContinuousSeries{
				Name:    "false-pick-ups",
				XValues: xaxis,
				YValues: stats.falsePickUps,
			},
			chart.ContinuousSeries{
				Name:    "bumps",
				XValues: xaxis,
				YValues: stats.bumps,
			},
			chart.ContinuousSeries{
				Name:    "missed-rubbish",
				XValues: xaxis,
				YValues: stats.leftovers,
			},
		},
	}
	counters.Elements = []chart.Renderable{
		chart.LegendLeft(&counters),
	}

	return counters
}

func writeChart(name string, c chart.Chart) error {
	buffer := new(bytes.Buffer)
	err := c.Render(chart.PNG, buffer)
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(buffer.Bytes()); err != nil {
		return err
	}
	return nil
}

func evolve() {
	currentGen := createNextGeneration(nil)
	alphas := newStatistics()
	runts := newStatistics()
	final := &strategy{}
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
		alphas.add(alpha)
		runts.add(runt)
		log.Printf("Finished generation %d: alpha %d, runt %d\n", i, alpha.score, runt.score)
		currentGen = createNextGeneration(currentGen)
		final = alpha
	}

	if err := writeChart(fmt.Sprintf("%s-alpha-scores.png", outputFile), alphas.getScoreChart()); err != nil {
		log.Println(err)
	}
	if err := writeChart(fmt.Sprintf("%s-alpha-counters.png", outputFile), alphas.getCountersChart()); err != nil {
		log.Println(err)
	}
	if err := writeChart(fmt.Sprintf("%s-alpha-genome.png", outputFile), alphas.getGenomeChart()); err != nil {
		log.Println(err)
	}
	if err := writeChart(fmt.Sprintf("%s-runt-scores.png", outputFile), runts.getScoreChart()); err != nil {
		log.Println(err)
	}
	if err := writeChart(fmt.Sprintf("%s-runt-counters.png", outputFile), runts.getCountersChart()); err != nil {
		log.Println(err)
	}
	if err := writeChart(fmt.Sprintf("%s-runt-genome.png", outputFile), runts.getGenomeChart()); err != nil {
		log.Println(err)
	}
	log.Println(final)
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
