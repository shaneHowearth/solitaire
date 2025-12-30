package tui

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/shanehowearth/solitaire/state"
)

func (display *Display) ShowHintModal(moves []state.Move) {
	// Create the text view.
	hintList := tview.NewTextView().
		SetDynamicColors(true).
		SetText(display.formatMoveList(moves))
	hintList.SetBorder(true).SetTitle(" Hints ")

	// Wrap it in a flex to centre it.
	modal := tview.NewFlex().
		AddItem(nil, 0, 1, false).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(nil, 0, 1, false).
			AddItem(hintList, 0, 1, true).
			AddItem(nil, 0, 1, false), 0, 4, true).
		AddItem(nil, 0, 1, false)

	// Close logic.
	hintList.SetInputCapture(func(*tcell.EventKey) *tcell.EventKey {
		display.pages.RemovePage("hints")
		// Refocus the board specifically.
		display.App.SetFocus(display.screens[display.Selected.Name()])

		return nil
	})

	// Add the page as a modal (true, true).
	display.pages.AddPage("hints", modal, true, true)
}

func (display *Display) formatMoveList(moves []state.Move) string {
	if len(moves) == 0 {
		return "\n [yellow]No moves found. Try drawing from stock."
	}

	var sBuild strings.Builder
	sBuild.WriteString(" [blue:][b]Available Hints (Any key to close)[-][:]\n\n")

	rows, cols := display.Selected.TableauGridSize()

	for idx, move := range moves {
		card := formatCardForTUI(move.SourceCardTop)
		src := display.prettyLocation(move.Source, rows, cols)
		dst := display.prettyLocation(move.Destination, rows, cols)

		seq := ""
		if move.NumberMoving > 1 {
			seq = fmt.Sprintf(" [gray](+%d)[-]", move.NumberMoving-1)
		}

		sBuild.WriteString(fmt.Sprintf(" [white]%d) [yellow]%-5s%s [white]%-8s [blue]→ [white]%s\n",
			idx+1, card, seq, src, dst))
	}

	return sBuild.String()
}

func (*Display) prettyLocation(stack state.Stack, rows, cols int) string {
	switch stack.Type {
	case state.StackTableau:
		if rows > 1 {
			r := (stack.TableauPosition / cols) + 1
			c := (stack.TableauPosition % cols) + 1

			return fmt.Sprintf("Row %d,Col %d", r, c)
		}

		return fmt.Sprintf("Tableau %d", stack.TableauPosition+1)
	case state.StackFoundation:
		return fmt.Sprintf("Foundation %d", stack.FoundationPosition+1)
	case state.StackWaste:
		return "Waste"
	case state.StackReserve:
		return "Reserve"
	default:
		return "Stock"
	}
}

// formatCardForTUI should handle the red/black colouring.
func formatCardForTUI(c state.SuitedCard) string {
	before := "[-]"
	if c.Suit == state.Hearts || c.Suit == state.Diamonds {
		before = "[red]"
	}

	return before + c.String() + "[-]"
}
