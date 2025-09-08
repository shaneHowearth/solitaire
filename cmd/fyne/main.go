package main

import (
	"log"
	"os"

	"github.com/shanehowearth/solitaire"
	"github.com/shanehowearth/solitaire/game"
	"github.com/shanehowearth/solitaire/screen/gui"
)

func main() {
	logFile, err := os.OpenFile("fyne_app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Redirect the standard 'log' package's output to this file
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	instance := solitaire.New()

	// Available games
	variants := []game.Variant{}
	variants = append(variants, &game.Klondike{})
	variants = append(variants, &game.KlondikeVegas{})
	variants = append(variants, &game.Accordian{})
	variants = append(variants, &game.AcesAndKings{})
	variants = append(variants, &game.AcesSquare{})
	variants = append(variants, &game.Acme{})
	variants = append(variants, &game.Agnes{})
	variants = append(variants, &game.Gaps{})
	variants = append(variants, &game.Russian{})
	variants = append(variants, &game.Yukon{})

	instance.Display = gui.New(variants)

	if err := instance.Start(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
