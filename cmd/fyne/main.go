package main

import (
	_ "image/gif"
	"log"
	"os"

	"fyne.io/fyne/v2/app"
	"github.com/shanehowearth/solitaire"
	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/screen/gui" // Updated import
	"github.com/shanehowearth/solitaire/state"
)

func main() {
	// ... Logging Setup (remains identical) ...
	logFile, err := os.OpenFile("gui_app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, state.DefaultLogPerms)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// ... (Rest of logging setup stays the same) ...

	// 1. Initialize the Fyne App with a unique ID for Persistence
	// This ID is the key to your "filing cabinet" for Recently Played games.
	myApp := app.NewWithID("com.solitaire.game.2026")

	// --- Instance Initialization ---
	instance := solitaire.New()

	// Available games
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
		&game.Easthaven{},
		&game.FlowerGarden{},
		&game.Gaps{},
		&game.Russian{},
		&game.SirTommy{},
		&game.WestcliffAmerican{},
		&game.WestcliffClassic{},
		&game.Whitehead{},
		&game.Yukon{},
	}

	// 2. Pass the app instance into your GUI constructor
	// You will need to update your gui.New function signature to accept fyne.App
	instance.Display = gui.New(myApp, variants)

	// Start the game engine
	if err := instance.Start(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
