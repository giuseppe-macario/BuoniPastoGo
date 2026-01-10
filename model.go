package main

// ======================================================================================
// MODEL (Domain Layer)
// Rappresenta i dati puri, senza conoscenza della GUI.
// ======================================================================================

type TimesheetEntry struct {
	Date       string
	DayOfWeek  string
	EntryTime  string
	ExitTime   string
	MealStatus string // "PRANZO", "CENA", "PRANZO e CENA"
	Note       string
}
