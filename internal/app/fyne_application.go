package app

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

// LiftMetricsApp encapsulates the main application structure
type LiftMetricsApp struct {
	app    fyne.App
	window fyne.Window
}

// New creates and returns a new instance of LiftMetricsApp
func New() *LiftMetricsApp {
	fyneApp := app.New()
	liftApp := &LiftMetricsApp{
		app:    fyneApp,
		window: fyneApp.NewWindow("LiftMetrics"),
	}
	liftApp.window.Resize(fyne.NewSize(800, 600))
	return liftApp
}

// Run sets up the UI and starts the application
func (a *LiftMetricsApp) Run() {
	a.setupUI()
	a.window.ShowAndRun()
}

// setupUI constructs the user interface of the application
func (a *LiftMetricsApp) setupUI() {
	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search lifters...")
	searchButton := widget.NewButton("Search", func() {
		// TODO: Implement search functionality
	})
	searchContainer := container.NewBorder(nil, nil, nil, searchButton, searchEntry)

	lifterInfo := container.NewVBox(
		widget.NewLabel("Name: John Doe"),
		widget.NewLabel("Age: 30"),
		widget.NewLabel("Weight: 80 kg"),
	)

	metricsInfo := container.NewVBox(
		widget.NewLabel("Squat: 150 kg"),
		widget.NewLabel("Bench: 100 kg"),
		widget.NewLabel("Deadlift: 180 kg"),
		widget.NewLabel("Total: 430 kg"),
	)

	chartPlaceholder := widget.NewLabel("Chart will be displayed here")

	grid := container.New(layout.NewGridLayout(2),
		lifterInfo,
		metricsInfo,
		chartPlaceholder,
		widget.NewLabel("Additional info or visualization can go here"),
	)

	content := container.NewBorder(searchContainer, nil, nil, nil, grid)
	a.window.SetContent(content)
}
