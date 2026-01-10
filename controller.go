package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

// ======================================================================================
// CONTROLLER
// Gestisce l'UI e orchestra le chiamate al Model/Service.
// ======================================================================================

type AppController struct {
	window  fyne.Window
	service *PDFService
	data    []TimesheetEntry
	table   *widget.Table
	status  *widget.Label
}

func NewAppController(w fyne.Window) *AppController {
	return &AppController{
		window:  w,
		service: &PDFService{},
		data:    []TimesheetEntry{},
	}
}

// handleOpenFile gestisce l'evento click (Command Pattern)
func (c *AppController) handleOpenFile() {
	fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			dialog.ShowError(err, c.window)
			return
		}
		if reader == nil {
			return // Utente ha annullato
		}
		defer reader.Close()

		filePath := reader.URI().Path()
		c.status.SetText("Elaborazione in corso: " + reader.URI().Name())

		// Esegui parsing
		entries, err := c.service.ProcessFile(filePath)
		if err != nil {
			dialog.ShowError(fmt.Errorf("Errore lettura PDF: %v", err), c.window)
			c.status.SetText("Errore durante l'elaborazione.")
			return
		}

		// Aggiorna Dati e UI
		c.data = entries
		c.table.Refresh() // Notifica la tabella che i dati sono cambiati

		if len(entries) == 0 {
			c.status.SetText("Nessun buono pasto trovato nel file.")
		} else {
			c.status.SetText(fmt.Sprintf("Trovati %d giorni con diritto al pasto.", len(entries)))
		}

	}, c.window)

	fd.SetFilter(storage.NewExtensionFileFilter([]string{".pdf"}))
	fd.Show()
}
