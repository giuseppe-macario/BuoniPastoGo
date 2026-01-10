package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// ======================================================================================
// VIEW
// Costruisce l'interfaccia grafica.
// ======================================================================================

// BuildUI costruisce l'interfaccia grafica
func (c *AppController) BuildUI() {
	// 1. Header
	title := widget.NewLabelWithStyle("Calcolatore Buoni Pasto", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	// 2. Status Label
	c.status = widget.NewLabel("Seleziona un file PDF per iniziare...")
	c.status.Alignment = fyne.TextAlignCenter

	// 3. Table Configuration
	c.table = widget.NewTable(
		// Length: quante righe e colonne
		func() (int, int) {
			return len(c.data), 5 // 5 colonne
		},
		// Create: crea l'elemento grafico della cella
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell Content")
		},
		// Update: aggiorna i dati della cella
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			entry := c.data[id.Row]
			switch id.Col {
			case 0:
				label.SetText(fmt.Sprintf("%s (%s)", entry.Date, entry.DayOfWeek))
			case 1:
				label.SetText(entry.EntryTime)
			case 2:
				label.SetText(entry.ExitTime)
			case 3:
				label.SetText(entry.MealStatus)
				label.TextStyle = fyne.TextStyle{Bold: true}
			case 4:
				label.SetText(entry.Note)
			}
		},
	)

	// Imposta larghezza colonne
	c.table.SetColumnWidth(0, 180)
	c.table.SetColumnWidth(1, 80)
	c.table.SetColumnWidth(2, 80)
	c.table.SetColumnWidth(3, 150)
	c.table.SetColumnWidth(4, 250)

	// 4. Toolbar / Buttons
	btnOpen := widget.NewButton("Apri PDF", c.handleOpenFile)
	btnOpen.Importance = widget.HighImportance

	// Layout
	content := container.NewBorder(
		container.NewVBox(title, btnOpen, c.status, widget.NewSeparator()), // Top
		nil, // Bottom
		nil, // Left
		nil, // Right
		c.table, // Center (si espande)
	)

	c.window.SetContent(content)
}
