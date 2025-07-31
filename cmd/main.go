package main

import (
	"log"
	"os"

	"github.com/shanehowearth/solitaire"
	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/screen/tui"
)

func main() {

	logFile, err := os.OpenFile("tview_app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close() // Ensure the log file is closed when main() exits

	// --- Step 2: Redirect the standard 'log' package's output to this file ---
	log.SetOutput(logFile)

	// You can also set log flags if you want timestamps, file/line numbers etc.
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	// Which variants of state are available to play.
	variants := []game.Variant{}
	variants = append(variants, &game.Klondike{})
	variants = append(variants, &game.Klondike2{})

	instance := solitaire.New(tui.New(variants))
	concreteDisplay, ok := instance.Display.(*tui.Display)
	if !ok {
		// This should theoretically not happen if solitaire.New always returns *tui.Display
		// but it's good practice for type assertions.
		log.Fatalf("Error: The Display implementation is not of type *tui.Display")
	}

	// Now you can access 'App' on the concrete type
	if err := concreteDisplay.App.Run(); err != nil {
		log.Fatalf("Error running tview application: %v", err)
	}
}
