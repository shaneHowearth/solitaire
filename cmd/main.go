package main

import (
	"log"
	"os"

	"github.com/shanehowearth/solitaire"
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

	instance := solitaire.New()

	if err := instance.Start(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
