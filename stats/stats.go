package stats

import (
	"bytes"
	"os"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/seq"

	"me/ga-robbie/robbie"
)

// Stats is a struct that tracks the statistics of one Robbie out of every
// generation.
type Stats struct {
	totalMoves  int
	reward      int
	minScore    int
	scores      []float64
	achievement []float64
	cleanliness []float64
	clumsiness  []float64
	scatiness   []float64
	// pickUps       []float64
	// falsePickUps []float64
	// bumps        []float64
	// missedRubbish []float64
	// totalRubbish  []float64
	numGenes map[robbie.Gene][]float64
}

// New returns a new and empty Stats object.
func New(totalMoves int, pickUpReward, pickUpPenalty, bumpPenalty int) *Stats {
	min := pickUpPenalty
	if bumpPenalty > pickUpPenalty {
		min = bumpPenalty
	}
	return &Stats{
		totalMoves:  totalMoves,
		reward:      pickUpReward,
		minScore:    min * totalMoves,
		scores:      make([]float64, 0),
		cleanliness: make([]float64, 0),
		clumsiness:  make([]float64, 0),
		scatiness:   make([]float64, 0),
		achievement: make([]float64, 0),
		// pickUps:       make([]float64, 0),
		// falsePickUps: make([]float64, 0),
		// bumps:        make([]float64, 0),
		// missedRubbish: make([]float64, 0),
		// totalRubbish:  make([]float64, 0),
		numGenes: make(map[robbie.Gene][]float64, 0),
	}
}

// Add adds the given Robbie to the statistics. It is assumed only one Robbie
// from each generation will be added to a single Stats object.
func (s *Stats) Add(r *robbie.Robbie) {
	scoreRange := float64(s.reward*r.TotalRubbish) + float64(s.minScore)
	s.scores = append(s.scores, float64(r.Score))
	s.achievement = append(s.achievement, 100.0*((float64(r.Score)+float64(s.minScore))/scoreRange))
	s.cleanliness = append(s.cleanliness, 100.0*float64(r.PickUps)/float64(r.TotalRubbish))
	s.clumsiness = append(s.clumsiness, 100.0*float64(r.Bumps)/float64(s.totalMoves))
	s.scatiness = append(s.scatiness, 100.0*float64(r.FalsePickUps)/float64(s.totalMoves))
	// s.pickUps = append(s.pickUps, float64(r.PickUps))
	// s.falsePickUps = append(s.falsePickUps, float64(r.FalsePickUps))
	// s.bumps = append(s.bumps, float64(r.Bumps))
	// s.missedRubbish = append(s.missedRubbish, float64(r.MissedRubbish))
	// s.totalRubbish = append(s.totalRubbish, float64(r.TotalRubbish))
	counters := make(map[robbie.Gene]float64)
	for _, g := range r.Genome {
		if _, ok := counters[g]; !ok {
			counters[g] = 1
		} else {
			counters[g]++
		}
	}
	for k, v := range counters {
		s.numGenes[k] = append(s.numGenes[k], v)
	}
}

// GetScoreChart saves the Stats score chart to the filename.
// func (s *Stats) GetScoreChart(filename string) error {
// 	xaxis := seq.Range(0.0, float64(len(s.scores)-1))
// 	score := chart.Chart{
// 		XAxis: chart.XAxis{
// 			Style: chart.StyleShow(),
// 		},
// 		YAxis: chart.YAxis{
// 			Style: chart.StyleShow(),
// 		},
// 		Background: chart.Style{
// 			Padding: chart.Box{
// 				Top:  20,
// 				Left: 120,
// 			},
// 		},
// 		Series: []chart.Series{
// 			chart.ContinuousSeries{
// 				Name:    "score",
// 				XValues: xaxis,
// 				YValues: s.scores,
// 			},
// 			chart.ContinuousSeries{
// 				Name:    "missed-rubbish",
// 				XValues: xaxis,
// 				YValues: s.missedRubbish,
// 			},
// 		},
// 	}
// 	score.Elements = []chart.Renderable{
// 		chart.LegendLeft(&score),
// 	}

// 	return writeChart(filename, score)
// }

// GetGenomeChart saves the Stats genome chart to the filename.
func (s *Stats) GetGenomeChart(filename string) error {
	maxGenerations := len(s.scores)
	xaxis := seq.Range(0.0, float64(maxGenerations-1))
	doNothings := make([]float64, maxGenerations)
	pickUps := make([]float64, maxGenerations)
	moveUps := make([]float64, maxGenerations)
	moveRights := make([]float64, maxGenerations)
	moveDowns := make([]float64, maxGenerations)
	moveLefts := make([]float64, maxGenerations)
	moveRandoms := make([]float64, maxGenerations)
	for i := 0; i < maxGenerations; i++ {
		doNothings[i] = s.numGenes[robbie.DoNothing][i]
		pickUps[i] = s.numGenes[robbie.PickUpRubbish][i] + doNothings[i]
		moveUps[i] = s.numGenes[robbie.MoveUp][i] + pickUps[i]
		moveRights[i] = s.numGenes[robbie.MoveRight][i] + moveUps[i]
		moveDowns[i] = s.numGenes[robbie.MoveDown][i] + moveRights[i]
		moveLefts[i] = s.numGenes[robbie.MoveLeft][i] + moveDowns[i]
		moveRandoms[i] = s.numGenes[robbie.MoveRandom][i] + moveLefts[i]
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
					StrokeColor: chart.ColorBlue,
					FillColor:   chart.ColorBlue,
				},
				XValues: xaxis,
				YValues: moveRandoms,
				Name:    "move-random",
			},
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.ColorYellow,
					FillColor:   chart.ColorYellow,
				},
				XValues: xaxis,
				YValues: moveLefts,
				Name:    "move-left",
			},
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.ColorCyan,
					FillColor:   chart.ColorCyan,
				},
				XValues: xaxis,
				YValues: moveDowns,
				Name:    "move-down",
			},
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.ColorOrange,
					FillColor:   chart.ColorOrange,
				},
				XValues: xaxis,
				YValues: moveRights,
				Name:    "move-right",
			},
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.ColorGreen,
					FillColor:   chart.ColorGreen,
				},
				XValues: xaxis,
				YValues: moveUps,
				Name:    "move-up",
			},
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.ColorRed,
					FillColor:   chart.ColorRed,
				},
				XValues: xaxis,
				YValues: pickUps,
				Name:    "pick-up-rubbish",
			},
			chart.ContinuousSeries{
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.ColorBlack,
					FillColor:   chart.ColorBlack,
				},
				XValues: xaxis,
				YValues: doNothings,
				Name:    "do-nothing",
			},
		},
	}
	genomes.Elements = []chart.Renderable{
		chart.LegendLeft(&genomes),
	}

	return writeChart(filename, genomes)
}

// GetCountersChart saves the Stats couters to the filename.
func (s *Stats) GetCountersChart(filename string) error {
	xaxis := seq.Range(0.0, float64(len(s.scores)-1))
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
				Name:    "cleanliness",
				XValues: xaxis,
				YValues: s.cleanliness,
			},
			chart.ContinuousSeries{
				Name:    "clumsiness",
				XValues: xaxis,
				YValues: s.clumsiness,
			},
			chart.ContinuousSeries{
				Name:    "scatiness",
				XValues: xaxis,
				YValues: s.scatiness,
			},
			chart.ContinuousSeries{
				Name:    "achievement",
				XValues: xaxis,
				YValues: s.achievement,
			},
			// chart.ContinuousSeries{
			// 	Name:    "pick-ups",
			// 	XValues: xaxis,
			// 	YValues: s.pickUps,
			// },
			// chart.ContinuousSeries{
			// 	Name:    "false-pick-ups",
			// 	XValues: xaxis,
			// 	YValues: s.falsePickUps,
			// },
			// chart.ContinuousSeries{
			// 	Name:    "bumps",
			// 	XValues: xaxis,
			// 	YValues: s.bumps,
			// },
			// chart.ContinuousSeries{
			// 	Name:    "missed-rubbish",
			// 	XValues: xaxis,
			// 	YValues: s.missedRubbish,
			// },
			// chart.ContinuousSeries{
			// 	Name:    "total-rubbish",
			// 	XValues: xaxis,
			// 	YValues: s.totalRubbish,
			// },
		},
	}
	counters.Elements = []chart.Renderable{
		chart.LegendLeft(&counters),
	}

	return writeChart(filename, counters)
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
