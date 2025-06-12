# Games

The definition of a game should be put in this directory.

The current implementation searches this directory for all possible variants of
games.

The game definitions must be an implementation of the Game interface.

In order to define a game the following needs to be done.

## Deck (or pack)

How many decks are required per game. At this point a deck is 52 cards, 4 suits,
13 ranks, but there may be [Piquet](https://en.wikipedia.org/wiki/Piquet_pack)
packs available.

## [Foundations](https://en.wikipedia.org/wiki/Glossary_of_patience_terms#foundations)

How many foundations are required? Note that this will be a multiple of 4, one for each
suit.

What is the [base card](https://en.wikipedia.org/wiki/Glossary_of_patience_terms#foundation_card) to be used. Note that this needs to be expressed as a solitaire.Rank

What is the rule for adding a card to a given foundation. Typically this is
something like "Must be the same suit as the foundation, and must be the next
card UP in rank".

example:
```Go
func (card solitaire.SuitedCard) bool {
    // Handle when the foundation is empty.
    if foundation.Len() == 0 {
        if card.Suit == foundation.BaseCard.Suit && card.Rank == foundation.BaseCard.Rank {
            return true
        }
    }

    // Get the card currently at the top of the foundation.
    topCard, err := foundation.Top()
    if err != nil {
        return false
    }

    // If the card is the same suit, and is one up in rank
    // then it can go onto the foundation.
    if card.Suit == foundation.BaseCard.Suit && (card.Rank - topCard.Rank) == 1 {
        return true
    }

    // All other cases the card should not be added to the foundation.
    return false
}
```

## [Tableau](https://en.wikipedia.org/wiki/Glossary_of_patience_terms#tableau)

How many tableau are required? Note that this can be any number. If omitted then
there will not be any created.

What is the [base card](https://en.wikipedia.org/wiki/Glossary_of_patience_terms#foundation_card) to be used. Note that this needs to be expressed as a solitaire.Rank

What is the rule for adding a card to a given tableau.

example:
```Go
func (card solitaire.SuitedCard) bool {
    // Handle when the tableau is empty.
    if tableau.Len() == 0 {
        if card.Rank == tableau.BaseCard.Rank {
            return true
        }
    }

    // Get the card currently at the top of the tableau.
    topCard, err := tableau.Top()
    if err != nil {
        return false
    }

    // If the card is the same suit, and is one down in rank
    // then it can go onto the tableau.
    if card.Suit == tableau.BaseCard.Suit && (topCard.Rank - card.Rank) == 1 {
        return true
    }

    // All other cases the card should not be added to the tableau.
    return false
}
```
