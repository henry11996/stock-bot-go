package plot

import (
	"fmt"

	"github.com/wcharczuk/go-chart"
)

// BaseSeries draws a horizontal line at the minimum value of the inner series.
type BaseSeries struct {
	Name        string
	Style       chart.Style
	YAxis       chart.YAxisType
	InnerSeries chart.ValuesProvider

	BaseValue *float64
}

// GetName returns the name of the time series.
func (ms BaseSeries) GetName() string {
	return ms.Name
}

// GetStyle returns the line style.
func (ms BaseSeries) GetStyle() chart.Style {
	return ms.Style
}

// GetYAxis returns which YAxis the series draws on.
func (ms BaseSeries) GetYAxis() chart.YAxisType {
	return ms.YAxis
}

// Len returns the number of elements in the series.
func (ms BaseSeries) Len() int {
	return ms.InnerSeries.Len()
}

// GetValues gets a value at a given index.
func (ms *BaseSeries) GetValues(index int) (x, y float64) {
	x, _ = ms.InnerSeries.GetValues(index)
	y = *ms.BaseValue
	return
}

// Render renders the series.
func (ms *BaseSeries) Render(r chart.Renderer, canvasBox chart.Box, xrange, yrange chart.Range, defaults chart.Style) {
	style := ms.Style.InheritFrom(defaults)
	chart.Draw.LineSeries(r, canvasBox, xrange, yrange, style, ms)
}

// Validate validates the series.
func (ms *BaseSeries) Validate() error {
	if ms.InnerSeries == nil {
		return fmt.Errorf("baseline series requires InnerSeries to be set")
	}
	return nil
}
