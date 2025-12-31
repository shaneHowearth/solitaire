package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/shanehowearth/solitaire"
	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/screen/tui"
	"github.com/shanehowearth/solitaire/state"
)

func main() {
	logFile, err := os.OpenFile("tview_app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, state.DefaultLogPerms)
	if err != nil {
		logFile.Close() // Explicitly close before exiting.
		log.Fatalf("Failed to open log file: %v", err)
	}

	var programLevel = new(slog.LevelVar) // Info by default
	h := slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))
	programLevel.Set(slog.LevelDebug)

	// --- Step 2: Redirect the standard 'log' package's output to this file ---.
	log.SetOutput(logFile)

	// You can also set log flags if you want timestamps, file/line numbers etc.
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	instance := solitaire.New()

	// Available games.
	variants := []game.Variant{}
	variants = append(variants, &game.Klondike{})
	variants = append(variants, &game.KlondikeVegas{})
	variants = append(variants, &game.Accordian{})
	variants = append(variants, &game.AcesAndKings{})
	variants = append(variants, &game.AcesSquare{})
	variants = append(variants, &game.AcesUp{})
	variants = append(variants, &game.Acme{})
	variants = append(variants, &game.Agnes{})
	variants = append(variants, &game.Gaps{})
	variants = append(variants, &game.Russian{})
	variants = append(variants, &game.Yukon{})

	instance.Display = tui.New(variants)

	if err := instance.Start(); err != nil {
		logFile.Close() // Explicitly close before exiting.
		log.Fatalf("Error running application: %v", err)
	}
}
