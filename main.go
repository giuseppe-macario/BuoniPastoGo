package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

// ======================================================================================
// MAIN
// ======================================================================================

func main() {
	// Inizializzazione App Fyne
	a := app.New()
	w := a.NewWindow("Gestione Buoni Pasto - Militari")

	w.Resize(fyne.NewSize(800, 600))

	// Inizializza Controller e costruisci UI
	ctrl := NewAppController(w)
	ctrl.BuildUI()

	w.ShowAndRun()
}
