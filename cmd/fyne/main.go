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
	variants = append(variants, &game.Acme{})
	variants = append(variants, &game.AcesAndKings{})
	variants = append(variants, &game.AcesSquare{})
	variants = append(variants, &game.Gaps{})

	// Replace tui.New with your Fyne display
	instance.Display = gui.New(variants)

	if err := instance.Start(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}

	// myApp := app.New()
	// w := myApp.NewWindow("Irate Sol")

	// // Green Rectangle.
	// gameArea := canvas.NewRectangle(color.RGBA{R: 46, G: 125, B: 50, A: 255})

	// // List of all games
	// games := []string{
	// 	"Klondike", "Spider", "FreeCell", "Gaps", "Pyramid", "Golf",
	// 	"Yukon", "Canfield", "Clock", "Accordion", "Baker's Game",
	// 	"Calculation", "Forty Thieves", "Freecell Baker", "Eight Off",
	// 	"Seahaven Towers", "Cruel", "La Belle Lucie", "Fan", "Trefoil",
	// 	// Add more games as needed...
	// }

	// // Main dropdown menu
	// var mainMenu *widget.Select

	// // Function to show all games in a dialog
	// showAllGamesDialog := func() {
	// 	gamesList := widget.NewList(
	// 		func() int { return len(games) },
	// 		func() fyne.CanvasObject { return widget.NewLabel("Game Name") },
	// 		func(id widget.ListItemID, object fyne.CanvasObject) {
	// 			object.(*widget.Label).SetText(games[id])
	// 		},
	// 	)

	// 	gamesList.Resize(fyne.NewSize(300, 250))

	// 	gamesList.OnSelected = func(id widget.ListItemID) {
	// 		selectedGame := games[id]
	// 		fmt.Printf("Starting %s\n", selectedGame)
	// 		mainMenu.SetSelected(selectedGame)
	// 	}

	// 	content := container.NewVBox(
	// 		widget.NewLabel("Choose a game:"),
	// 		gamesList,
	// 	)

	// 	gamesDialog := dialog.NewCustom("All Games", "Close", content, w)
	// 	gamesDialog.Resize(fyne.NewSize(400, 350))
	// 	gamesDialog.Show()
	// }

	// // Function to create games submenu
	// createGamesSubmenu := func() {
	// 	maxGamesToShow := 10
	// 	var menuItems []*fyne.MenuItem

	// 	for i, game := range games {
	// 		if i >= maxGamesToShow {
	// 			break
	// 		}
	// 		gameName := game // Capture for closure
	// 		menuItems = append(menuItems, fyne.NewMenuItem(gameName, func() {
	// 			fmt.Printf("Starting %s\n", gameName)
	// 			mainMenu.SetSelected(gameName)
	// 		}))
	// 	}

	// 	// Add "More games..." option if there are more games
	// 	if len(games) > maxGamesToShow {
	// 		menuItems = append(menuItems, fyne.NewMenuItem("More games...", func() {
	// 			showAllGamesDialog()
	// 		}))
	// 	}

	// 	// Create and show the popup menu
	// 	menu := fyne.NewMenu("", menuItems...)
	// 	popUp := widget.NewPopUpMenu(menu, w.Canvas())
	// 	popUp.ShowAtPosition(fyne.NewPos(10, 60))
	// }

	// // Create the main dropdown with games included as options
	// allOptions := append([]string{"Menu", "New Game", "Games"}, games...)
	// allOptions = append(allOptions, "Statistics", "Options", "Help")

	// mainMenu = widget.NewSelect(allOptions, func(selected string) {
	// 	switch selected {
	// 	case "Menu":
	// 		return // Do nothing for placeholder
	// 	case "Games":
	// 		createGamesSubmenu()
	// 	case "New Game":
	// 		fmt.Println("Starting new game...")
	// 	case "Statistics":
	// 		fmt.Println("Showing statistics...")
	// 	case "Options":
	// 		fmt.Println("Opening options...")
	// 	case "Help":
	// 		fmt.Println("Showing help...")
	// 	default:
	// 		// Check if it's a game name
	// 		for _, game := range games {
	// 			if selected == game {
	// 				fmt.Printf("Starting %s\n", selected)
	// 				return
	// 			}
	// 		}
	// 	}
	// })

	// mainMenu.SetSelected("Menu") // Set initial selection

	// menuContainer := container.NewHBox(
	// 	mainMenu,
	// 	widget.NewButton("Undo", func() {}),
	// 	widget.NewButton("Hint", func() {}),
	// )

	// mainContent := container.NewBorder(
	// 	menuContainer, // Menu at top
	// 	nil,
	// 	nil,
	// 	nil,
	// 	gameArea, // Your card table
	// )

	// w.SetContent(mainContent)
	// width, height := float32(1200), float32(800)
	// w.Resize(fyne.NewSize(width, height))
	// w.ShowAndRun()
}
