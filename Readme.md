This is a solitaire game that can be played in the terminal.

Currently the game is played using a mouse, there may be plans in the future to
allow the game(s) to be played using the keyboard only.

# Play.
Different variants have different rules of play, however in terms of the UI, the
player clicks on one component (Stock, Waste, Foundation, or Tableau) which will
change colour. The player then clicks on another component (or the same again)
to indicate where the cards will go.

If the move is allowed then the card(s) will move from the first component to
the second.

Movement of card(s) is greedy, that is, if multiple cards can be moved
(generally only from one tableau to another) then the maximum number of cards
that can be moved will be moved.

Happy Playing.

![gameplay](gameplay.gif)

# Build.
Clone the source locally, and run:
```
$ go build cmd/tui/main.go -o irateOils
```

# Installation.
Until builds are created then build the project using the above instructions,
and then move the resulting binary to somewhere on your $PATH.

# Variants.
The [games/Readme](game/README.md) discusses implementing the Variant interface. The Variant interface describes how the game looks and is played. Please feel free to use the Klondike.go example in games/ to create your own
version.
