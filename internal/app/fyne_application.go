package app

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// LiftMetricsApp encapsulates the main application structure
type LiftMetricsApp struct {
	app     fyne.App
	window  fyne.Window
	db      *sql.DB
	dataDir string
}

// New creates and returns a new instance of LiftMetricsApp
func New(database *sql.DB, dataDir string) *LiftMetricsApp {
	fyneApp := app.New()
	liftApp := &LiftMetricsApp{
		app:     fyneApp,
		window:  fyneApp.NewWindow("LiftMetrics"),
		db:      database,
		dataDir: dataDir,
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
	lifterNames, err := a.loadLifterNames()
	if err != nil {
		dialog.ShowError(err, a.window)
		return
	}

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder("Search lifters...")

	popup := widget.NewPopUp(nil, a.window.Canvas())
	var suggestionList *widget.List

	updateSuggestions := func(searchTerm string) {
		filtered := a.filterLifters(lifterNames, searchTerm)
		log.Printf("Filtered lifters for '%s': %d results", searchTerm, len(filtered))

		if len(filtered) > 0 {
			suggestionList = widget.NewList(
				func() int { return len(filtered) },
				func() fyne.CanvasObject {
					return widget.NewLabel("placeholder")
				},
				func(id widget.ListItemID, item fyne.CanvasObject) {
					item.(*widget.Label).SetText(filtered[id])
				},
			)
			suggestionList.OnSelected = func(id widget.ListItemID) {
				searchEntry.SetText(filtered[id])
				popup.Hide()
				// TODO: Update the UI with selected lifter's data
				log.Printf("Selected lifter: %s", filtered[id])
			}
			popup.Content = container.NewVBox(suggestionList)
			popup.Resize(fyne.NewSize(searchEntry.Size().Width, 200))
			popup.Move(fyne.NewPos(searchEntry.Position().X, searchEntry.Position().Y+searchEntry.Size().Height))
			popup.Show()
		} else {
			popup.Hide()
		}
	}

	searchEntry.OnChanged = func(s string) {
		if s == "" {
			popup.Hide()
		} else {
			updateSuggestions(s)
		}
	}

	lifterInfo := container.NewVBox(
		widget.NewLabel("Name: "),
		widget.NewLabel("Age: "),
		widget.NewLabel("Weight: "),
	)

	metricsInfo := container.NewVBox(
		widget.NewLabel("Squat: "),
		widget.NewLabel("Bench: "),
		widget.NewLabel("Deadlift: "),
		widget.NewLabel("Total: "),
	)

	chartPlaceholder := widget.NewLabel("Chart will be displayed here")

	content := container.NewBorder(
		searchEntry,
		nil,
		nil,
		nil,
		container.NewGridWithColumns(2,
			lifterInfo,
			metricsInfo,
			chartPlaceholder,
			widget.NewLabel("Additional info or visualization can go here"),
		),
	)

	a.window.SetContent(content)
}

func (a *LiftMetricsApp) loadLifterNames() ([]string, error) {
	jsonFilePath := filepath.Join(a.dataDir, "lifters.json")
	jsonData, err := os.ReadFile(jsonFilePath)
	if err != nil {
		log.Printf("Error reading lifters.json: %v", err)
		return nil, err
	}

	var lifterNames []string
	err = json.Unmarshal(jsonData, &lifterNames)
	if err != nil {
		log.Printf("Error unmarshaling lifter names: %v", err)
		return nil, err
	}

	log.Printf("Loaded %d lifter names", len(lifterNames))
	if len(lifterNames) > 0 {
		log.Printf("First lifter name: %s", lifterNames[0])
	}

	return lifterNames, nil
}

// Add this new method to filter lifters based on input
func (a *LiftMetricsApp) filterLifters(allLifters []string, input string) []string {
	var filtered []string
	input = strings.ToLower(input)
	for _, name := range allLifters {
		if strings.Contains(strings.ToLower(name), input) {
			filtered = append(filtered, name)
			if len(filtered) >= 5 { // Limit to 5 results for better performance
				break
			}
		}
	}
	return filtered
}
