# Port Go di [BuoniPastoCS](https://github.com/giuseppe-macario/BuoniPastoCS)

Rispetto alla versione C#.NET, ha il vantaggio di produrre un solo file, e non un insieme di file (eseguibile e librerie). Inoltre, l'eseguibile è circa dieci volte più piccolo di quello prodotto da C#.NET, quindi molto più leggero.

Lo svantaggio è una GUI (in Fyne) più spartana, non nativa e per il momento senza drag and drop.

## Su macOS

Per compilare: `go build`, oppure per eseguire senza produrre un eseguibile: `go run .` (attenzione al punto finale).

Per creare un `.app` per macOS:
```
go install fyne.io/tools/cmd/fyne@latest
fyne package -os darwin -icon icon.png -release
```
