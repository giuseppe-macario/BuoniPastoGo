package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
)

// ======================================================================================
// BUSINESS LOGIC LAYER (Service)
// Contiene la logica di parsing e le regole di business.
// Implementa il pattern "Adapter" convertendo il PDF grezzo in oggetti Model.
// ======================================================================================

type PDFService struct{}

var (
	dateRe   = regexp.MustCompile(`(\d{2}/\d{2}/\d{4})`)
	orarioRe = regexp.MustCompile(`(\d{1,2}[:.,]\d{2})`)
	days     = []string{"domenica", "lunedì", "martedì", "mercoledì", "giovedì", "venerdì", "sabato"}
)

func (s *PDFService) ProcessFile(path string) ([]TimesheetEntry, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []TimesheetEntry

	// Iteriamo su tutte le pagine
	for pageIndex := 1; pageIndex <= r.NumPage(); pageIndex++ {
		p := r.Page(pageIndex)
		if p.V.IsNull() {
			continue
		}

		rows, _ := p.GetTextByRow()
		for _, row := range rows {
			// Ricostruzione riga
			var lineBuffer bytes.Buffer
			for _, word := range row.Content {
				lineBuffer.WriteString(word.S)
				lineBuffer.WriteString(" ")
			}
			line := lineBuffer.String()

			// Parsing riga
			parsed, err := s.parseLine(line)
			if err != nil {
				continue // Riga non valida o non contiene dati rilevanti
			}

			// Calcolo Pasto
			mealType := s.calculateMeal(parsed.dateStr, parsed.exitTime)
			if mealType == "" {
				continue // Nessun pasto maturato
			}

			// Formattazione
			wdIt := s.getWeekday(parsed.dateStr)
			entryStr := parsed.entryTime.Format("15:04")

			// Logica asterisco per recupero compensativo
			if parsed.causale == "RECUPERO COMPENSATIVO" && parsed.entryTime.Hour() == 7 && parsed.entryTime.Minute() == 30 {
				entryStr = "*07:30"
			}

			entries = append(entries, TimesheetEntry{
				Date:       parsed.dateStr,
				DayOfWeek:  wdIt,
				EntryTime:  entryStr,
				ExitTime:   parsed.exitTime.Format("15:04"),
				MealStatus: mealType,
				Note:       parsed.causale,
			})
		}
	}

	return entries, nil
}

// Struttura interna di supporto per il parsing
type rawParsedData struct {
	dateStr   string
	entryTime time.Time
	exitTime  time.Time
	causale   string
}

func (s *PDFService) parseLine(line string) (*rawParsedData, error) {
	line = strings.TrimLeft(line, "* ")
	line = strings.TrimSpace(line)

	dateMatch := dateRe.FindStringSubmatch(line)
	if dateMatch == nil {
		return nil, fmt.Errorf("no date")
	}

	timeMatches := orarioRe.FindAllString(line, -1)
	timeIndices := orarioRe.FindAllStringIndex(line, -1)
	if len(timeMatches) < 2 {
		return nil, fmt.Errorf("no times")
	}

	oraIng := s.normalizeTime(timeMatches[0])
	oraUsc := s.normalizeTime(timeMatches[1])

	// Ignora 00:00 - 00:00
	if oraIng.IsZero() && oraUsc.IsZero() {
		return nil, fmt.Errorf("zero times")
	}

	// Estrazione causale (euristica posizionale)
	var startIdx, endIdx int
	if len(timeIndices) >= 3 {
		startIdx = timeIndices[2][1]
	} else {
		startIdx = timeIndices[1][1]
	}
	if len(timeIndices) >= 4 {
		endIdx = timeIndices[3][0]
	} else {
		endIdx = len(line)
	}

	if startIdx > len(line) {
		startIdx = len(line)
	}
	if endIdx > len(line) {
		endIdx = len(line)
	}

	causaleRaw := line[startIdx:endIdx]
	causale := strings.Join(strings.Fields(causaleRaw), " ") // Rimuove spazi extra
	causale = strings.Trim(causale, "-.,; ")

	// Pulizia specifica richiesta
	if causale == "COMANDO E LOGISTICA" {
		causale = ""
	}

	return &rawParsedData{
		dateStr:   dateMatch[1],
		entryTime: oraIng,
		exitTime:  oraUsc,
		causale:   causale,
	}, nil
}

func (s *PDFService) normalizeTime(tstr string) time.Time {
	tstr = strings.ReplaceAll(tstr, ",", ":")
	tstr = strings.ReplaceAll(tstr, ".", ":")
	t, _ := time.Parse("15:04", tstr)
	return t
}

func (s *PDFService) getWeekday(dateStr string) string {
	d, _ := time.Parse("02/01/2006", dateStr)
	return days[d.Weekday()]
}

func (s *PDFService) calculateMeal(dateStr string, exitTime time.Time) string {
	d, _ := time.Parse("02/01/2006", dateStr)
	wd := d.Weekday()

	minExit := exitTime.Hour()*60 + exitTime.Minute()

	// Regole:
	// Pranzo: Ven, Sab, Dom se uscita >= 15:30 (930 min)
	isLunchDay := wd == time.Friday || wd == time.Saturday || wd == time.Sunday
	hasLunch := isLunchDay && minExit >= 930

	// Cena: Sempre se uscita >= 20:30 (1230 min)
	hasDinner := minExit >= 1230

	if hasLunch && hasDinner {
		return "PRANZO e CENA"
	}
	if hasDinner {
		return "CENA"
	}
	if hasLunch {
		return "PRANZO"
	}
	return ""
}
