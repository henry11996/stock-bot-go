package main

import (
	"bytes"
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
			StrokeWidth: 4,
			Show:        true,
		},
		XValues: xv,
		YValues: yv,
	}

	smaSeries := chart.SMASeries{
		Name:   "MA",
		Period: 5,
		Style: chart.Style{
			StrokeColor:     drawing.ColorRed,
			StrokeWidth:     4,
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

	max, _ := data.Meta.PriceHighLimit.Float64()
	min, _ := data.Meta.PriceLowLimit.Float64()
	graph := chart.Chart{
		Title: data.Meta.NameZhTw + "(" + data.Info.SymbolID + ")",
		TitleStyle: chart.Style{
			Show: true,
		},
		Width:        4096,
		Height:       2800,
		DPI:          400,
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
			Style: chart.Style{
				Show: true,
			},
		},
		Series: []chart.Series{
			bbSeries,
			priceSeries,
			smaSeries,
		},
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
		y, _ := data[time].Open.Float64()
		ys = append(ys, y)
	}
	return ys
}

var darkColorPalette DarkColorPalette

type DarkColorPalette struct{}

func (ap DarkColorPalette) BackgroundColor() drawing.Color {
	return chart.ColorBlack
}

func (ap DarkColorPalette) BackgroundStrokeColor() drawing.Color {
	return chart.ColorBlack
}

func (ap DarkColorPalette) CanvasColor() drawing.Color {
	return chart.ColorBlack
}

func (ap DarkColorPalette) CanvasStrokeColor() drawing.Color {
	return chart.ColorBlack
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
