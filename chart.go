package main

import (
	"bytes"
	"os"

	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/seq"
)

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
