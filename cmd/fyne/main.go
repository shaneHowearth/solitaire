package main

import (
	_ "image/gif"
	"log"
	"log/slog"
	"os"

	"github.com/shanehowearth/solitaire"
	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/screen/gui" // Updated import
	"github.com/shanehowearth/solitaire/state"
)

func main() {
	// --- Logging Setup (remains identical) ---
	logFile, err := os.OpenFile("gui_app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, state.DefaultLogPerms)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	var programLevel = new(slog.LevelVar)
	h := slog.NewJSONHandler(logFile, &slog.HandlerOptions{Level: programLevel})
	slog.SetDefault(slog.New(h))
	programLevel.Set(slog.LevelDebug)

	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// --- Instance Initialization ---
	instance := solitaire.New()

	// Available games (remains identical)
	variants := []game.Variant{
		&game.Klondike{},
		&game.KlondikeVegas{},
		&game.Accordion{},
		&game.Addiction{},
		&game.AcesAndKings{},
		&game.AcesSquare{},
		&game.AcesUp{},
		&game.Acme{},
		&game.Agnes{},
		&game.Algerian{},
		&game.AmericanToad{},
		&game.Gaps{},
		&game.Russian{},
		&game.Yukon{},
	}

	// Swap tui.New for gui.New
	// This satisfies the Display interface requirement of your solitaire.Instance
	instance.Display = gui.New(variants)

	// Start the game engine
	if err := instance.Start(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
