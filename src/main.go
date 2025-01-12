package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"net/http"
	"slices"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

var (
	baseHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "rt_base",
		Buckets: []float64{
			.001,
			.005,
			.015,
			.050,
		},
	}, []string{"type"})

	preciseHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "rt_precise",
		Buckets: []float64{
			.001,
			.005,
			.010,
			.020,
			.030,
			.040,
			.050,
		},
	}, []string{"type"})

	summary = promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: "rt_summary",
		Objectives: map[float64]float64{
			0.5:  0.0001,
			0.99: 0.0001,
			1:    0.0001,
		},
	},
		[]string{"type"},
	)

	nativeHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:                        "rt_native_hist",
		Buckets:                     nil,
		NativeHistogramBucketFactor: 1.1,
	}, []string{"type"})
)

const nSamples = 1000

func main() {
	const fastMean = 13 // 20
	const slowMean = 15 // 40
	fast, slow := generateDistribution(fastMean), generateDistribution(slowMean)

	visualise(fast, slow)

	quantiles := []int{50, 99, 100}

	printStatistics(quantiles, "fast", fast)
	printStatistics(quantiles, "slow", slow)

	go func() {
		for {
			for i := 0; i < nSamples; i++ {
				baseHistogram.WithLabelValues("slow").Observe(float64(slow[i]) / 1000.0)
				baseHistogram.WithLabelValues("fast").Observe(float64(fast[i]) / 1000.0)
			}
			for i := 0; i < nSamples; i++ {
				preciseHistogram.WithLabelValues("slow").Observe(float64(slow[i]) / 1000.0)
				preciseHistogram.WithLabelValues("fast").Observe(float64(fast[i]) / 1000.0)
			}
			for i := 0; i < nSamples; i++ {
				summary.WithLabelValues("slow").Observe(float64(slow[i]) / 1000.0)
				summary.WithLabelValues("fast").Observe(float64(fast[i]) / 1000.0)
			}
			for i := 0; i < nSamples; i++ {
				nativeHistogram.WithLabelValues("slow").Observe(float64(slow[i]) / 1000.0)
				nativeHistogram.WithLabelValues("fast").Observe(float64(fast[i]) / 1000.0)
			}

			time.Sleep(1 * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe("localhost:8006", nil))
}

func generateDistribution(mean float64) []int {
	r := rand.New(rand.NewSource(42))

	vals := make([]int, 0, nSamples)
	for range nSamples {
		vals = append(vals, int(r.NormFloat64()+mean))

	}
	return vals
}

func calcQuantile(values []int, q int) int {
	slices.Sort(values)
	pos := int(math.Ceil(float64(len(values)) * float64(q) / 100))
	return values[pos-1]
}

func visualise(fast []int, slow []int) {
	p := plot.New()
	p.Title.Text = "fast/slow"

	const nBuckets = 10
	v := make(plotter.Values, len(fast))
	for i := 0; i < len(fast); i++ {
		v[i] = float64(fast[i])
	}

	h, err := plotter.NewHist(v, nBuckets)
	if err != nil {
		panic(err)
	}
	h.Color = color.RGBA{R: 255, A: 255}
	h.FillColor = nil
	p.Add(h)

	v = make(plotter.Values, len(slow))
	for i := 0; i < len(slow); i++ {
		v[i] = float64(slow[i])
	}

	h, err = plotter.NewHist(v, nBuckets)
	if err != nil {
		panic(err)
	}
	h.Color = color.RGBA{B: 255, A: 255}
	h.FillColor = nil
	p.Add(h)

	p.X.Scale = plot.LinearScale{}
	p.X.Min = -5
	p.X.Max = 50

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "hist.png"); err != nil {
		panic(err)
	}
}

func printStatistics(qs []int, title string, vals []int) {
	fmt.Print(title + ": ")
	for _, quantile := range qs {
		fmt.Printf("p%d=%d ", quantile, calcQuantile(vals, quantile))
	}
	fmt.Println()
}
