package main

import (
	"bytes"
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/henry11996/fugle-golang/fugle"
	"github.com/wcharczuk/go-chart"
	"github.com/wcharczuk/go-chart/drawing"
)

func newPlot(data fugle.Data) []byte {
	xv := xvalues(data.Chart)
	yv := yvalues(data.Chart, xv)

	priceSeries := chart.TimeSeries{
		Name: "Price",
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(0),
			StrokeWidth: 6,
			Show:        true,
		},
		XValues: xv,
		YValues: yv,
	}

	minSeries := &chart.MinSeries{
		Style: chart.Style{
			Show:            true,
			StrokeColor:     chart.ColorWhite,
			StrokeDashArray: []float64{3.0, 3.0},
			FontColor:       chart.ColorAlternateGreen,
		},
		InnerSeries: priceSeries,
	}
	maxSeries := &chart.MaxSeries{
		Style: chart.Style{
			Show:            true,
			StrokeColor:     chart.ColorWhite,
			StrokeDashArray: []float64{3.0, 3.0},
			FontColor:       chart.ColorRed,
		},
		InnerSeries: priceSeries,
	}

	smaSeries := chart.SMASeries{
		Name:   "MA5",
		Period: 5,
		Style: chart.Style{
			StrokeColor:     drawing.ColorFromHex("ccefff"),
			StrokeWidth:     6,
			StrokeDashArray: []float64{5.0, 5.0},
			Show:            true,
		},
		InnerSeries: priceSeries,
	}

	bbSeries := &chart.BollingerBandsSeries{
		Name: "BBands",
		Style: chart.Style{
			StrokeColor: drawing.ColorFromHex("707070"),
			FillColor:   drawing.ColorFromHex("707070").WithAlpha(64),
			Show:        true,
		},
		InnerSeries: priceSeries,
	}

	base, _ := data.Meta.PriceReference.Float64()
	baseSeries := &BaseSeries{
		Style: chart.Style{
			Show:            true,
			StrokeWidth:     3,
			StrokeColor:     chart.ColorWhite,
			StrokeDashArray: []float64{1.0, 1.0},
		},
		InnerSeries: priceSeries,
		BaseValue:   &base,
	}

	max, _ := data.Meta.PriceHighLimit.Float64()
	min, _ := data.Meta.PriceLowLimit.Float64()
	yticks := priceTicks(base, min, max)
	graph := chart.Chart{
		Title: data.Meta.NameZhTw + "(" + data.Info.SymbolID + ") " + data.Info.Date,
		TitleStyle: chart.Style{
			Show: true,
		},
		Width:        4096,
		Height:       2800,
		DPI:          400,
		Font:         DefaultFont,
		ColorPalette: darkColorPalette,
		XAxis: chart.XAxis{
			TickPosition: chart.TickPositionBetweenTicks,
			Style: chart.Style{
				Show: true,
			},
			ValueFormatter: chart.TimeValueFormatterWithFormat("15:04"),
		},
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{
				Max: max,
				Min: min,
			},
			Ticks: yticks,
			Style: chart.Style{
				Show: true,
			},
		},
		Series: []chart.Series{
			baseSeries,
			maxSeries,
			minSeries,
			bbSeries,
			smaSeries,
			priceSeries,
			chart.LastValueAnnotation(minSeries),
			chart.LastValueAnnotation(maxSeries),
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph, chart.Style{
			FillColor: BlackColor,
			FontColor: chart.ColorWhite,
		}),
	}

	buffer := bytes.NewBuffer([]byte{})
	err := graph.Render(chart.PNG, buffer)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}

func xvalues(data map[time.Time]fugle.Deal) []time.Time {
	times := make([]time.Time, 0)
	for time := range data {
		times = append(times, time)
	}
	sort.Slice(times, func(i, j int) bool {
		return times[i].Before(times[j])
	})
	return times
}

func yvalues(data map[time.Time]fugle.Deal, times []time.Time) []float64 {
	ys := []float64{}
	for _, time := range times {
		y, _ := data[time].Close.Float64()
		ys = append(ys, y)
	}
	return ys
}

func priceTicks(start float64, min float64, max float64) []chart.Tick {
	ticks := make([]chart.Tick, 0)
	ticks = append(ticks, chart.Tick{Value: min, Label: fmt.Sprintf("%f", min)})
	step := (max - min) / 10
	for i := 0; i < 9; i++ {
		v := min + step*float64((i+1))
		if i == 4 {
			ticks = append(ticks, chart.Tick{Value: start, Label: fmt.Sprintf("%f", start)})
		} else {
			ticks = append(ticks, chart.Tick{Value: v, Label: fmt.Sprintf("%f", v)})
		}
	}
	ticks = append(ticks, chart.Tick{Value: max, Label: fmt.Sprintf("%f", max)})
	return ticks
}

var darkColorPalette DarkColorPalette

var BlackColor = drawing.Color{R: 10, G: 10, B: 10, A: 255}

type DarkColorPalette struct{}

func (ap DarkColorPalette) BackgroundColor() drawing.Color {
	return BlackColor
}

func (ap DarkColorPalette) BackgroundStrokeColor() drawing.Color {
	return BlackColor
}

func (ap DarkColorPalette) CanvasColor() drawing.Color {
	return BlackColor
}

func (ap DarkColorPalette) CanvasStrokeColor() drawing.Color {
	return BlackColor
}

func (ap DarkColorPalette) AxisStrokeColor() drawing.Color {
	return chart.ColorWhite
}

func (ap DarkColorPalette) TextColor() drawing.Color {
	return chart.ColorWhite
}

func (ap DarkColorPalette) GetSeriesColor(index int) drawing.Color {
	return chart.GetAlternateColor(index)
}
