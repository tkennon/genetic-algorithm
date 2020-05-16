package stats

import (
	"io"

	"github.com/tkennon/genetic-algorithm/game"
	"github.com/wcharczuk/go-chart"
)

// Stats is a struct that tracks the statistics of one Robbie out of every
// generation.
type Stats struct {
	totalMoves  int
	reward      int
	scores      []float64
	cleanliness []float64
	clumsiness  []float64
	scatiness   []float64
}

// New returns a new and empty Stats object.
func New(totalMoves int, pickUpReward, pickUpPenalty, bumpPenalty int) *Stats {
	return &Stats{
		totalMoves:  totalMoves,
		reward:      pickUpReward,
		scores:      make([]float64, 0),
		cleanliness: make([]float64, 0),
		clumsiness:  make([]float64, 0),
		scatiness:   make([]float64, 0),
	}
}

// Add adds the given Robbie to the statistics. It is assumed only one Robbie
// from each generation will be added to a single Stats object.
func (s *Stats) Add(p *game.Player) {
	s.scores = append(s.scores, float64(p.Score))
	s.cleanliness = append(s.cleanliness, 100.0*float64(p.PickUps)/float64(p.TotalRubbish))
	s.clumsiness = append(s.clumsiness, 100.0*float64(p.Bumps)/float64(s.totalMoves))
	s.scatiness = append(s.scatiness, 100.0*float64(p.FalsePickUps)/float64(s.totalMoves))
}

// GetCountersChart writes the counters chart to the given writer.
func (s *Stats) GetCountersChart(w io.Writer) error {
	xaxis := chart.LinearRange(0.0, float64(len(s.scores)-1))
	counters := chart.Chart{
		XAxis: chart.XAxis{
			Style: chart.Style{},
		},
		YAxis: chart.YAxis{
			Style: chart.Style{},
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
		},
	}
	counters.Elements = []chart.Renderable{
		chart.LegendLeft(&counters),
	}

	return c.Render(chart.PNG, w)
}
