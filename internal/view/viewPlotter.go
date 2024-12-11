package view

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"log"
	"math"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func (v *View) openPlot() {
	plotWindow := fyne.CurrentApp().NewWindow(fmt.Sprintf("Plot: %s", getTruncatedLegendLabel(v.displayLabel.Text)))
	v.showPlot(plotWindow)
}

func (v *View) showPlot(mainWindow fyne.Window) {
	if v.displayLabel.Text != "0" {
		v.presenter.SaveHistory()
	}

	xMinEntry, xMaxEntry := widget.NewEntry(), widget.NewEntry()
	xMinEntry.SetText("-10")
	xMaxEntry.SetText("10")
	yMinEntry, yMaxEntry := widget.NewEntry(), widget.NewEntry()
	yMinEntry.SetText("-10")
	yMaxEntry.SetText("10")

	setupEntryValidation(xMinEntry)
	setupEntryValidation(xMaxEntry)
	setupEntryValidation(yMinEntry)
	setupEntryValidation(yMaxEntry)

	plotCanvas := v.createPlotCanvas(xMinEntry, xMaxEntry, yMinEntry, yMaxEntry)

	plotContent := container.New(layout.NewCenterLayout(), plotCanvas)

	buttonBackgroundColor := color.NRGBA{R: 220, G: 185, B: 240, A: 128}

	refreshPlotText := canvas.NewText("Refresh Plot", color.Black)
	refreshPlotText.Alignment = fyne.TextAlignCenter
	refreshPlotText.TextStyle = fyne.TextStyle{Bold: true}

	buttonBackground := canvas.NewRectangle(buttonBackgroundColor)
	buttonBackground.SetMinSize(fyne.NewSize(120, 40))

	refreshPlotButton := widget.NewButton("", func() {
		if !validateNumericEntries(xMinEntry, xMaxEntry, yMinEntry, yMaxEntry) {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Invalid Input",
				Content: "Please enter valid numeric values for all fields.",
			})
			return
		}

		newCanvas := v.createPlotCanvas(xMinEntry, xMaxEntry, yMinEntry, yMaxEntry)
		plotContent.Objects = []fyne.CanvasObject{newCanvas}
		mainWindow.Canvas().Refresh(plotContent)
	})

	buttonWithBackground := container.NewStack(
		refreshPlotButton,
		buttonBackground,
		container.NewCenter(refreshPlotText),
	)

	horizontalLayout := container.NewVBox(
		container.NewHBox(widget.NewLabel("xMin:"), xMinEntry, widget.NewLabel("xMax:"), xMaxEntry),
		container.NewHBox(widget.NewLabel("yMin:"), yMinEntry, widget.NewLabel("yMax:"), yMaxEntry),
		buttonWithBackground,
	)

	verticalLayout := container.NewVBox(plotContent, horizontalLayout)

	mainWindow.SetContent(verticalLayout)
	mainWindow.Resize(fyne.NewSize(550, 550))
	mainWindow.SetFixedSize(true)
	mainWindow.Show()
}

func (v *View) createPlotCanvas(xMinEntry, xMaxEntry, yMinEntry, yMaxEntry *widget.Entry) fyne.CanvasObject {
	p := plot.New()
	p.Title.Text = "Plot"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	minValueX := parseFloatFromEntry(xMinEntry, -10)
	maxValueX := parseFloatFromEntry(xMaxEntry, 10)
	minValueY := parseFloatFromEntry(yMinEntry, -10)
	maxValueY := parseFloatFromEntry(yMaxEntry, 10)

	p.X.Min = math.Max(minValueX, -1000000)
	p.X.Max = math.Min(maxValueX, 1000000)
	p.Y.Min = math.Max(minValueY, -1000000)
	p.Y.Max = math.Min(maxValueY, 1000000)

	p.X.Tick.Marker = createAdaptiveTicks(p.X.Min, p.X.Max)
	p.Y.Tick.Marker = createAdaptiveTicks(p.Y.Min, p.Y.Max)

	p.Add(plotter.NewGrid())

	points := v.generatePlotPoints(minValueX, maxValueX)
	validPoints := filterValidPlotPoints(points)

	err := plotutil.AddLinePoints(p, getTruncatedLegendLabel(v.displayLabel.Text), validPoints)
	if err != nil {
		log.Printf("Failed to plot data: %v", err)
		return widget.NewLabel("Error: Unable to plot data")
	}

	newCanvas, err := plotToCanvas(p)
	if err != nil {
		log.Printf("Failed to render plot: %v", err)
		return widget.NewLabel("Error: Unable to render plot")
	}

	return newCanvas
}

func createAdaptiveTicks(min, max float64) plot.Ticker {
	return plot.TickerFunc(func(min, max float64) []plot.Tick {
		rangeSize := max - min
		var step float64

		if rangeSize <= 10 {
			step = 1
		} else if rangeSize <= 100 {
			step = 5
		} else if rangeSize <= 1000 {
			step = 50
		} else if rangeSize <= 10000 {
			step = 100
		} else if rangeSize <= 100000 {
			step = 500
		} else {
			step = 1000
		}

		ticks := []plot.Tick{}
		for x := math.Ceil(min/step) * step; x <= max; x += step {
			label := fmt.Sprintf("%.0f", x)
			ticks = append(ticks, plot.Tick{Value: x, Label: label})
		}
		return ticks
	})
}

func setupEntryValidation(entry *widget.Entry) {
	entry.OnChanged = func(content string) {
		if _, err := strconv.ParseFloat(content, 64); err != nil && content != "" {
			entry.SetText(content[:len(content)-1])
		}
	}
}

func validateNumericEntries(entries ...*widget.Entry) bool {
	for _, entry := range entries {
		if _, err := strconv.ParseFloat(entry.Text, 64); err != nil {
			return false
		}
	}
	return true
}

func getTruncatedLegendLabel(label string) string {
	if label == "" {
		return "No Formula"
	}
	if len(label) > 50 {
		return label[:47] + "..."
	}
	return label
}

func parseFloatFromEntry(entry *widget.Entry, defaultValue float64) float64 {
	value, err := strconv.ParseFloat(entry.Text, 64)
	if err != nil {
		return defaultValue
	}
	return value
}

func filterValidPlotPoints(points plotter.XYs) plotter.XYs {
	validPoints := make(plotter.XYs, 0, len(points))
	for _, pt := range points {
		if !math.IsNaN(pt.Y) {
			validPoints = append(validPoints, pt)
		}
	}
	return validPoints
}

func (v *View) generatePlotPoints(n float64, m float64) plotter.XYs {
	len := 1000
	pts := make(plotter.XYs, len)
	interval := (m - n) / float64(len)

	currentX := n
	for i := 0; i < len; i++ {
		currentX = math.RoundToEven(currentX*1e6) / 1e6

		pts[i].X = currentX
		s := v.displayLabel.Text

		yStr, err := v.presenter.CalculatePlotResult(&s, fmt.Sprintf("%f", currentX))
		if err != nil {
			pts[i].Y = math.NaN()
		} else {
			y, convErr := strconv.ParseFloat(yStr, 64)
			if convErr != nil {
				pts[i].Y = math.NaN()
			} else {
				pts[i].Y = y
			}
		}
		currentX += interval
	}

	return pts
}

func plotToCanvas(p *plot.Plot) (fyne.CanvasObject, error) {
	img := vgimg.New(vg.Points(500), vg.Points(500))
	dc := draw.New(img)

	p.Draw(dc)

	rawImage := img.Image()

	var buf bytes.Buffer
	if err := png.Encode(&buf, rawImage); err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	fyneImage := canvas.NewImageFromReader(&buf, "plot.png")
	fyneImage.FillMode = canvas.ImageFillOriginal

	return fyneImage, nil
}
