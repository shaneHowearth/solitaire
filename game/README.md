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
func (foundation solitaire.Foundation, card solitaire.SuitedCard) bool {
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
func (tableau solitaire.Tableau, card solitaire.SuitedCard) bool {
    // Handle when the tableau is empty.
    if tableau.Len() == 0 {
        if card.Rank == tableau.Card.Rank {
            return true
        }
    }

    // Get the card currently at the top of the tableau.
    topCard, err := tableau.Top()
    if err != nil {
        return false
    }

	// If the card is the opposite colour, and is one down in rank
	// then it can go onto the tableau.
	if ((card.Suit+topCard.Suit)%2 == 1) && (topCard.Rank-card.Rank) == 1 {
		return true
	}

    // All other cases the card should not be added to the tableau.
    return false
}
```

## Possible Variants [Source](https://en.wikipedia.org/wiki/List_of_patience_games)
### A
- [x] [Accordion](https://en.wikipedia.org/wiki/Accordion_(solitaire))
- [x] [Aces and Kings](https://en.wikipedia.org/wiki/Aces_and_Kings)
- [x] [Aces Square](https://en.wikipedia.org/wiki/Aces_Square_(solitaire))
- [x] [Aces Up](https://en.wikipedia.org/wiki/Aces_Up)
- [x] [Acme](https://en.wikipedia.org/wiki/Acme_(solitaire))
- [x] [Addiction](https://en.wikipedia.org/wiki/Addiction_(solitaire))
- [x] [Agnes](https://en.wikipedia.org/wiki/Agnes_(card_game))
- [ ] [Alaska](https://en.wikipedia.org/wiki/Alaska_(solitaire))
- [x] [Algerian](https://en.wikipedia.org/wiki/Algerian_(solitaire))
- [ ] [Alhambra](https://en.wikipedia.org/wiki/Alhambra_(solitaire))
- [ ] [Amazons](https://en.wikipedia.org/wiki/Amazons_(solitaire))
- [x] [American Toad](https://en.wikipedia.org/wiki/American_Toad_(solitaire))
- [ ] [Apophis](https://en.wikipedia.org/wiki/Apophis_(solitaire))
- [x] [Appreciate](https://en.wikipedia.org/wiki/Appreciate_(solitaire))
- [ ] [Acquaintance](https://en.wikipedia.org/wiki/Acquaintance_(solitaire))
- [ ] [Archway](https://en.wikipedia.org/wiki/Archway_(solitaire))
- [ ] [Auld Lang Syne](https://en.wikipedia.org/wiki/Auld_Lang_Syne_(solitaire))
- [x] [Australian Patience (Canberra)](https://en.wikipedia.org/wiki/Australian_Patience)
- [x] [Australian Patience (Tasmanian)](https://en.wikipedia.org/wiki/Australian_Patience)
### B
- [ ] [Babette](https://en.wikipedia.org/wiki/Babette_(card_game))
- [ ] [Backbone](https://en.wikipedia.org/wiki/Backbone_(solitaire))
- [ ] [Baker's Dozen](https://en.wikipedia.org/wiki/Baker%27s_Dozen_(solitaire))
- [ ] [Baker's Game](https://en.wikipedia.org/wiki/Baker%27s_Game)
- [ ] [Baroness](https://en.wikipedia.org/wiki/Baroness_(solitaire))
- [ ] [Batsford](https://en.wikipedia.org/wiki/Batsford_(solitaire))
- [ ] [Beetle](https://en.wikipedia.org/wiki/Beetle_(solitaire))
- [ ] [Beleaguered Castle](https://en.wikipedia.org/wiki/Beleaguered_Castle)
- [ ] [Belvedere](https://en.wikipedia.org/wiki/Belvedere_(solitaire))
- [ ] [Betsy Ross](https://en.wikipedia.org/wiki/Betsy_Ross_(solitaire))
- [ ] [Big Ben](https://en.wikipedia.org/wiki/Big_Ben_(solitaire))
- [ ] [Big Forty](https://en.wikipedia.org/wiki/Big_Forty_(solitaire))
- [ ] [Big Harp](https://en.wikipedia.org/wiki/Big_Harp_(solitaire))
- [ ] [Birthday](https://en.wikipedia.org/wiki/Birthday_(patience))
- [ ] [Bisley](https://en.wikipedia.org/wiki/Bisley_(solitaire))
- [ ] [Black Hole](https://en.wikipedia.org/wiki/Black_Hole_(solitaire))
- [ ] [Block 10](https://en.wikipedia.org/wiki/Block_10_(solitaire))
- [ ] [Blockade](https://en.wikipedia.org/wiki/Blockade_(solitaire))
- [ ] [Bowling Solitaire](https://en.wikipedia.org/wiki/Bowling_(solitaire))
- [ ] [Box Kite](https://en.wikipedia.org/wiki/Box_Kite_(solitaire))
- [ ] [Braid](https://en.wikipedia.org/wiki/Braid_(solitaire))
- [ ] [Brigade](https://en.wikipedia.org/wiki/Brigade_(solitaire))
- [ ] [Bristol](https://en.wikipedia.org/wiki/Bristol_(solitaire))
- [ ] [British Constitution](https://en.wikipedia.org/wiki/British_Constitution_(solitaire))
- [ ] [British Square](https://en.wikipedia.org/wiki/British_Square_(card_game))
- [ ] [Broken Intervals](https://en.wikipedia.org/wiki/Broken_Intervals)
- [ ] [Busy Aces](https://en.wikipedia.org/wiki/Busy_Aces_(solitaire))
### C
- [x] [Calculation](https://en.wikipedia.org/wiki/Calculation_(card_game))
- [ ] [Canfield](https://en.wikipedia.org/wiki/Canfield_(solitaire))
- [ ] [Capricieuse](https://en.wikipedia.org/wiki/Capricieuse)
- [ ] [Carpet](https://en.wikipedia.org/wiki/Carpet_(solitaire))
- [ ] [Carthage](https://en.wikipedia.org/wiki/Carthage_(solitaire))
- [ ] [Casket](https://en.wikipedia.org/wiki/Casket_(solitaire))
- [ ] [Castles in Spain](https://en.wikipedia.org/wiki/Castles_in_Spain_(solitaire))
- [ ] [Chameleon](https://en.wikipedia.org/wiki/Canfield_(solitaire)#Variations)
- [ ] [Chessboard](https://en.wikipedia.org/wiki/Chessboard_(solitaire))
- [ ] [Cicely](https://en.wikipedia.org/wiki/Cicely_(solitaire))
- [ ] [Citadel](https://en.wikipedia.org/wiki/Citadel_(solitaire))
- [ ] [Clock Patience](https://en.wikipedia.org/wiki/Clock_Patience)
- [ ] [Colorado](https://en.wikipedia.org/wiki/Colorado_(game))
- [ ] [Colours](https://en.wikipedia.org/wiki/Colours_(solitaire))
- [ ] [Concentration](https://en.wikipedia.org/wiki/Concentration_(card_game))
- [ ] [Congress](https://en.wikipedia.org/wiki/Congress_(solitaire))
- [ ] [Contradance](https://en.wikipedia.org/wiki/Contradance_(solitaire))
- [ ] [Corner Card](https://en.wikipedia.org/wiki/Corner_Card_(solitaire))
- [ ] [Corner Patience](https://en.wikipedia.org/wiki/Corner_Patience_(solitaire))
- [ ] [Corners](https://en.wikipedia.org/wiki/Corners_(patience))
- [ ] [Corona](https://en.wikipedia.org/wiki/Corona_(solitaire))
- [ ] [Constitution](https://en.wikipedia.org/wiki/British_Constitution_(solitaire))
- [ ] [Cotillion](https://en.wikipedia.org/wiki/Contradance_(solitaire))
- [ ] [Crapette+](https://en.wikipedia.org/wiki/Crapette)
- [ ] [Courtyard](https://en.wikipedia.org/wiki/Deuces_(solitaire))
- [ ] [Crazy Quilt](https://en.wikipedia.org/wiki/Crazy_Quilt_(solitaire))
- [ ] [Crescent](https://en.wikipedia.org/wiki/Crescent_(solitaire))
- [ ] [Cribbage Solitaire](https://en.wikipedia.org/wiki/Cribbage_solitaire)
- [ ] [Cribbage Squares](https://en.wikipedia.org/wiki/Cribbage_Squares)
- [ ] [Cruel](https://en.wikipedia.org/wiki/Cruel_(solitaire))
- [ ] [Curds and Whey](https://en.wikipedia.org/wiki/Curds_and_Whey_(solitaire))
- [ ] [Czarina](https://en.wikipedia.org/wiki/Czarina_(solitaire))
### D
- [ ] [Decade](https://en.wikipedia.org/wiki/Decade_(solitaire))
- [ ] [Deuces](https://en.wikipedia.org/wiki/Deuces_(solitaire))
- [ ] [Devil's Grip](https://en.wikipedia.org/wiki/Devil%27s_Grip)
- [ ] [Diplomat](https://en.wikipedia.org/wiki/Diplomat_(solitaire))
- [ ] [Double Canfield](https://en.wikipedia.org/wiki/Double_Canfield_(solitaire))
- [ ] [Double Klondike+](https://en.wikipedia.org/wiki/Double_Klondike)
- [ ] [Double Solitaire+](https://en.wikipedia.org/wiki/Double_Klondike)
- [ ] [Doublets](https://en.wikipedia.org/wiki/Doublets_(solitaire))
- [ ] [Downing Street](https://en.wikipedia.org/wiki/Downing_Street_(solitaire))
- [ ] [Dress Parade](https://en.wikipedia.org/wiki/Dress_Parade_(solitaire))
- [ ] [Duchess](https://en.wikipedia.org/wiki/Duchess_(solitaire))
### E
- [ ] [Eagle Wing](https://en.wikipedia.org/wiki/Eagle_Wing)
- [x] [Easthaven](https://en.wikipedia.org/wiki/Easthaven_(solitaire))
- [ ] [Eight Cards](https://en.wikipedia.org/wiki/Eight_Cards)
- [ ] [Eight Off](https://en.wikipedia.org/wiki/Eight_Off)
- [ ] [Eighteens](https://en.wikipedia.org/wiki/Eighteens_(solitaire))
- [ ] [Elevens](https://en.wikipedia.org/wiki/Elevens)
- [ ] [Emperor](https://en.wikipedia.org/wiki/Emperor_(solitaire))
- [ ] [Emperor of Germany](https://en.wikipedia.org/wiki/Emperor_of_Germany_(solitaire))
- [ ] [Escalator](https://en.wikipedia.org/wiki/Escalator_(solitaire))
- [ ] [Exit](https://en.wikipedia.org/wiki/Exit_(solitaire))
### F
- [ ] [Faerie Queen](https://en.wikipedia.org/wiki/Faerie_Queen_(solitaire))
- [ ] [Fifteens](https://en.wikipedia.org/wiki/Fifteens_(solitaire))
- [ ] [Five Piles](https://en.wikipedia.org/wiki/Five_Piles_(solitaire))
- [ ] [Florentine Patience](https://en.wikipedia.org/wiki/Florentine_Patience_(solitaire))
- [x] [Flower Garden](https://en.wikipedia.org/wiki/Flower_Garden_(solitaire))
- [ ] [Fly](https://en.wikipedia.org/wiki/Fly_(solitaire))
- [ ] [Following](https://en.wikipedia.org/wiki/Following_(solitaire))
- [ ] [Fortress](https://en.wikipedia.org/wiki/Fortress_(solitaire))
- [ ] [Fortune's Favor](https://en.wikipedia.org/wiki/Fortune%27s_Favor)
- [ ] [Forty Thieves](https://en.wikipedia.org/wiki/Napoleon_at_St_Helena)
- [ ] [Four Corners](https://en.wikipedia.org/wiki/Four_Corners_(patience))
- [ ] [Four Seasons](https://en.wikipedia.org/wiki/Four_Seasons_(solitaire))
- [ ] [Four Winds](https://en.wikipedia.org/wiki/Four_Winds_(solitaire))
- [ ] [Fourteen Out](https://en.wikipedia.org/wiki/Fourteen_Out)
- [ ] [Fourteens](https://en.wikipedia.org/wiki/Fourteens_(solitaire))
- [x] [FreeCell](https://en.wikipedia.org/wiki/FreeCell)
- [ ] [Frog](https://en.wikipedia.org/wiki/Frog_(patience))
- [ ] [Frustration](https://en.wikipedia.org/wiki/Frustration_(solitaire))
### G
- [x] [Gaps](https://en.wikipedia.org/wiki/Gaps)
- [ ] [Gargantua](https://en.wikipedia.org/wiki/Gargantua_(solitaire))
- [ ] [Gate](https://en.wikipedia.org/wiki/Gate_(solitaire))
- [ ] [Gavotte](https://en.wikipedia.org/wiki/Gavotte_(solitaire))
- [ ] [Gay Gordons](https://en.wikipedia.org/wiki/Gay_Gordons_(solitaire))
- [ ] [German Clock](https://en.wikipedia.org/wiki/German_Clock)
- [ ] [German Patience](https://en.wikipedia.org/wiki/German_Patience)
- [ ] [Giant](https://en.wikipedia.org/wiki/Giant_(solitaire))
- [ ] [Giza](https://en.wikipedia.org/wiki/Giza_(solitaire))
- [ ] [Glencoe](https://en.wikipedia.org/wiki/Glencoe_(solitaire))
- [ ] [Golf](https://en.wikipedia.org/wiki/Golf_(patience))
- [ ] [Good Measure](https://en.wikipedia.org/wiki/Good_Measure_(solitaire))
- [ ] [Good Thirteen](https://en.wikipedia.org/wiki/Good_Thirteen)
- [ ] [Grampus](https://en.wikipedia.org/wiki/Grampus_(solitaire))
- [ ] [Granada](https://en.wikipedia.org/wiki/Granada_(solitaire))
- [ ] [Grand Duchess](https://en.wikipedia.org/wiki/Grand_Duchess_(solitaire))
- [ ] [Grandfather's Clock](https://en.wikipedia.org/wiki/Grandfather%27s_Clock_(solitaire))
- [ ] [Grandfather's Patience](https://en.wikipedia.org/wiki/Grandfather%27s_Patience)
- [ ] [Grandmother's Patience](https://en.wikipedia.org/wiki/Grandmother%27s_Patience)
### H
- [ ] [Harp](https://en.wikipedia.org/wiki/Harp_(solitaire))
- [ ] [Heads and Tails](https://en.wikipedia.org/wiki/Heads_and_Tails_(solitaire))
- [ ] [Herring-Bone](https://en.wikipedia.org/wiki/Herring-Bone)
- [ ] [Herz zu Herz](https://en.wikipedia.org/wiki/Herz_zu_Herz)
- [ ] [Hide-and-Seek](https://en.wikipedia.org/wiki/Hide-and-Seek_(solitaire))
- [ ] [Hit or Miss](https://en.wikipedia.org/wiki/Hit_or_Miss_(solitaire))
- [ ] [House in the Woods](https://en.wikipedia.org/wiki/House_in_the_Woods_(solitaire))
- [ ] [House on the Hill](https://en.wikipedia.org/wiki/House_on_the_Hill_(solitaire))
### I
- [ ] [Idiot's Delight](https://en.wikipedia.org/wiki/Idiot%27s_Delight_(solitaire))
- [ ] [Imaginary Thirteen](https://en.wikipedia.org/wiki/Imaginary_Thirteen)
- [ ] [Imperial Guards](https://en.wikipedia.org/wiki/Imperial_Guards_(solitaire))
- [ ] [Indian](https://en.wikipedia.org/wiki/Indian_(solitaire))
- [ ] [Indian Carpet](https://en.wikipedia.org/wiki/Indian_Carpet_(solitaire))
- [ ] [Interregnum](https://en.wikipedia.org/wiki/Interregnum_(solitaire))
- [ ] [Intrigue](https://en.wikipedia.org/wiki/Intrigue_(solitaire))
### J
- [ ] [Josephine](https://en.wikipedia.org/wiki/Josephine_(solitaire))
- [ ] [Jubilee](https://en.wikipedia.org/wiki/Jubilee_(solitaire))
### K
- [ ] [King Albert](https://en.wikipedia.org/wiki/King_Albert_(solitaire))
- [ ] [King Tut](https://en.wikipedia.org/wiki/King_Tut_(solitaire))
- [ ] [Kings in the Corners](https://en.wikipedia.org/wiki/Kings_in_the_Corners)
- [ ] [King's Audience](https://en.wikipedia.org/wiki/King%27s_Audience)
- [ ] [Kingsdown Eights](https://en.wikipedia.org/wiki/Kingsdown_Eights)
- [x] [Klondike](https://en.wikipedia.org/wiki/Klondike_(solitaire))
- [ ] [Knockout](https://en.wikipedia.org/wiki/Knockout_(solitaire))
- [x] [Kuipers](https://en.wikipedia.org/wiki/Klondike_(solitaire)#Variants)
### L
- [ ] [La Belle Lucie](https://en.wikipedia.org/wiki/La_Belle_Lucie)
- [ ] [La Chatelaine](https://en.wikipedia.org/wiki/La_Chatelaine)
- [ ] [La Croix d'Honneur](https://en.wikipedia.org/wiki/La_Croix_d%27Honneur)
- [ ] [Labyrinth](https://en.wikipedia.org/wiki/Labyrinth_(solitaire))
- [ ] [Lady Betty](https://en.wikipedia.org/wiki/Lady_Betty_(solitaire))
- [ ] [Lady of the Manor](https://en.wikipedia.org/wiki/Lady_of_the_Manor_(solitaire))
- [ ] [Laggard Lady](https://en.wikipedia.org/wiki/Laggard_Lady)
- [ ] [Las Vegas Solitaire](https://en.wikipedia.org/wiki/Las_Vegas_Solitaire)
- [ ] [Last Chance](https://en.wikipedia.org/wiki/Last_Chance_(solitaire))
- [ ] [Laying Siege](https://en.wikipedia.org/wiki/Laying_Siege)
- [ ] [Leoni's Own](https://en.wikipedia.org/wiki/Leoni%27s_Own)
- [ ] [Limited](https://en.wikipedia.org/wiki/Limited_(solitaire))
- [ ] [Little Milligan](https://en.wikipedia.org/wiki/Little_Milligan)
- [ ] [Little Spider](https://en.wikipedia.org/wiki/Little_Spider)
- [ ] [Little Windmill](https://en.wikipedia.org/wiki/Little_Windmill)
- [ ] [Long Braid](https://en.wikipedia.org/wiki/Long_Braid)
- [ ] [Lovely Lucy](https://en.wikipedia.org/wiki/Lovely_Lucy)
- [ ] [Louis](https://en.wikipedia.org/wiki/Louis_(solitaire))
- [ ] [Lucas](https://en.wikipedia.org/wiki/Lucas_(solitaire))
### M
- [ ] [Maria](https://en.wikipedia.org/wiki/Maria_(solitaire))
- [ ] [Martha](https://en.wikipedia.org/wiki/Martha_(solitaire))
- [ ] [Matrimony](https://en.wikipedia.org/wiki/Matrimony_(solitaire))
- [ ] [Maze](https://en.wikipedia.org/wiki/Maze_(solitaire))
- [ ] [Memory](https://en.wikipedia.org/wiki/Memory_(card_game))
- [ ] [Millie](https://en.wikipedia.org/wiki/Millie_(solitaire))
- [ ] [Milligan Cell](https://en.wikipedia.org/wiki/Milligan_Cell)
- [ ] [Milligan Harp](https://en.wikipedia.org/wiki/Milligan_Harp)
- [ ] [Milligan Yukon](https://en.wikipedia.org/wiki/Milligan_Yukon)
- [ ] [Miss Milligan](https://en.wikipedia.org/wiki/Miss_Milligan)
- [ ] [Montana](https://en.wikipedia.org/wiki/Montana_(solitaire))
- [ ] [Monte Carlo](https://en.wikipedia.org/wiki/Monte_Carlo_(solitaire))
- [ ] [Moojub](https://en.wikipedia.org/wiki/Moojub)
- [ ] [Mount Olympus](https://en.wikipedia.org/wiki/Mount_Olympus_(solitaire))
- [ ] [Mrs. Mop](https://en.wikipedia.org/wiki/Mrs._Mop)
### N
- [ ] [Narcotic](https://en.wikipedia.org/wiki/Narcotic_(solitaire))
- [ ] [Napoleon at St Helena](https://en.wikipedia.org/wiki/Napoleon_at_St_Helena)
- [ ] [Napoleon's Favorite](https://en.wikipedia.org/wiki/Napoleon%27s_Favorite)
- [ ] [Napoleon's Square](https://en.wikipedia.org/wiki/Napoleon%27s_Square)
- [ ] [Nerts+](https://en.wikipedia.org/wiki/Nerts)
- [ ] [Nestor](https://en.wikipedia.org/wiki/Nestor_(solitaire))
- [ ] [Nine Across](https://en.wikipedia.org/wiki/Klondike_(solitaire)#Variants)
- [ ] [Ninety-One](https://en.wikipedia.org/wiki/Ninety-One_(solitaire))
- [ ] [Nivernaise (La Nivernaise)](https://en.wikipedia.org/wiki/Nivernaise_(solitaire))
- [ ] [Number Ten](https://en.wikipedia.org/wiki/Number_Ten_(solitaire))
- [ ] [Numerica](https://en.wikipedia.org/wiki/Numerica)
### O
- [ ] [Odd and Even](https://en.wikipedia.org/wiki/Odd_and_Even)
- [ ] [Old Fashioned](https://en.wikipedia.org/wiki/Old_Fashioned_(solitaire))
- [ ] [Old Mole](https://en.wikipedia.org/wiki/Old_Mole)
- [ ] [Old Patience](https://en.wikipedia.org/wiki/Old_Patience)
- [ ] [One234](https://en.wikipedia.org/wiki/One234)
- [ ] [Osmosis](https://en.wikipedia.org/wiki/Osmosis_(solitaire))
### P
- [ ] [Päckchen](https://en.wikipedia.org/wiki/P%C3%A4ckchen_(solitaire))
- [ ] [Pairs](https://en.wikipedia.org/wiki/Pairs_(solitaire))
- [ ] [Parallels](https://en.wikipedia.org/wiki/Parallels_(solitaire))
- [ ] [Parisienne (La Parisienne, Parisian)](https://en.wikipedia.org/wiki/Parisienne_(solitaire))
- [ ] [Parliament](https://en.wikipedia.org/wiki/Parliament_(solitaire))
- [ ] [Pas de Deux](https://en.wikipedia.org/wiki/Pas_de_Deux_(solitaire))
- [ ] [Patience](https://en.wikipedia.org/wiki/Patience_(game))
- [ ] [Patriarchs](https://en.wikipedia.org/wiki/Patriarchs_(solitaire))
- [ ] [Penguin](https://en.wikipedia.org/wiki/Penguin_(solitaire))
- [ ] [Perpetual Motion](https://en.wikipedia.org/wiki/Perpetual_Motion_(solitaire))
- [ ] [Perseverance](https://en.wikipedia.org/wiki/Perseverance_(solitaire))
- [ ] [Persian Patience](https://en.wikipedia.org/wiki/Persian_Patience)
- [ ] [Persian Rug](https://en.wikipedia.org/wiki/Persian_Rug_(solitaire))
- [ ] [Pharaoh′s Grave](https://en.wikipedia.org/wiki/Pharaoh%27s_Grave)
- [ ] [Picture Gallery](https://en.wikipedia.org/wiki/Picture_Gallery_(solitaire))
- [ ] [Picture Patience](https://en.wikipedia.org/wiki/Picture_Patience)
- [ ] [Pigtail](https://en.wikipedia.org/wiki/Pigtail_(solitaire))
- [ ] [Plait](https://en.wikipedia.org/wiki/Plait_(solitaire))
- [ ] [Poker Squares](https://en.wikipedia.org/wiki/Poker_Squares)
- [ ] [Portuguese Solitaire](https://en.wikipedia.org/wiki/Portuguese_Solitaire)
- [ ] [Precedence (Order of Precedence)](https://en.wikipedia.org/wiki/Precedence_(solitaire))
- [ ] [Propeller](https://en.wikipedia.org/wiki/Propeller_(solitaire))
- [ ] [Puss in the Corner](https://en.wikipedia.org/wiki/Puss_in_the_Corner_(solitaire))
- [ ] [Putt Putt](https://en.wikipedia.org/wiki/Putt_Putt_(solitaire))
- [ ] [Pyramid](https://en.wikipedia.org/wiki/Pyramid_(solitaire))
- [ ] [Pyramide](https://en.wikipedia.org/wiki/Pyramide_(solitaire))
- [ ] [Pyramid Golf](https://en.wikipedia.org/wiki/Pyramid_Golf)
### Q
- [ ] [Quadrille](https://en.wikipedia.org/wiki/Quadrille_(solitaire))
- [ ] [Queen of Italy](https://en.wikipedia.org/wiki/Queen_of_Italy)
- [ ] [Queen's Audience](https://en.wikipedia.org/wiki/Queen%27s_Audience)
### R
- [ ] [Racing Demon+](https://en.wikipedia.org/wiki/Racing_Demon)
- [ ] [Raglan](https://en.wikipedia.org/wiki/Raglan_(solitaire))
- [ ] [Rainbow Canfield](https://en.wikipedia.org/wiki/Rainbow_Canfield)
- [ ] [Rank and File](https://en.wikipedia.org/wiki/Rank_and_File_(solitaire))
- [ ] [Red and Black](https://en.wikipedia.org/wiki/Red_and_Black_(solitaire))
- [ ] [Roosevelt at San Juan](https://en.wikipedia.org/wiki/Roosevelt_at_San_Juan)
- [ ] [Rosamund's Bower](https://en.wikipedia.org/wiki/Rosamund%27s_Bower)
- [ ] [Rouge et Noir](https://en.wikipedia.org/wiki/Rouge_et_Noir_(solitaire))
- [ ] [Royal Cotillion](https://en.wikipedia.org/wiki/Royal_Cotillion)
- [ ] [Royal Flush](https://en.wikipedia.org/wiki/Royal_Flush_(solitaire))
- [ ] [Royal Marriage](https://en.wikipedia.org/wiki/Royal_Marriage)
- [ ] [Royal Parade](https://en.wikipedia.org/wiki/Royal_Parade)
- [ ] [Royal Rendezvous](https://en.wikipedia.org/wiki/Royal_Rendezvous)
- [ ] [Russian Bank+](https://en.wikipedia.org/wiki/Russian_Bank)
- [x] [Russian](https://en.wikipedia.org/wiki/Yukon_(solitaire))
### S
- [ ] [Salic Law](https://en.wikipedia.org/wiki/Salic_Law_(solitaire))
- [ ] [Scorpion](https://en.wikipedia.org/wiki/Scorpion_(solitaire))
- [ ] [Seahaven Towers](https://en.wikipedia.org/wiki/Seahaven_Towers)
- [ ] [Seven Devils](https://en.wikipedia.org/wiki/Seven_Devils_(solitaire))
- [ ] [Sham Battle](https://en.wikipedia.org/wiki/Sham_Battle)
- [ ] [Shamrocks](https://en.wikipedia.org/wiki/Shamrocks_(solitaire))
- [ ] [Simple Simon](https://en.wikipedia.org/wiki/Simple_Simon_(solitaire))
- [ ] [Simplicity](https://en.wikipedia.org/wiki/Simplicity_(solitaire))
- [x] [Sir Tommy](https://en.wikipedia.org/wiki/Sir_Tommy)
- [ ] [Six By Six](https://en.wikipedia.org/wiki/Six_By_Six)
- [ ] [Sixes and Sevens](https://en.wikipedia.org/wiki/Sixes_and_Sevens_(solitaire))
- [ ] [Sixty Thieves](https://en.wikipedia.org/wiki/Sixty_Thieves)
- [ ] [Sly Fox](https://en.wikipedia.org/wiki/Sly_Fox_(solitaire))
- [ ] [Solitaire](https://en.wikipedia.org/wiki/Solitaire)
- [ ] [Somerset](https://en.wikipedia.org/wiki/Klondike_(solitaire)#Variants)
- [ ] [Spaces](https://en.wikipedia.org/wiki/Spaces_(solitaire))
- [ ] [Spanish Patience](https://en.wikipedia.org/wiki/Spanish_Patience)
- [ ] [Spider](https://en.wikipedia.org/wiki/Spider_(solitaire))
- [ ] [Spiderette](https://en.wikipedia.org/wiki/Spiderette)
- [ ] [Spiderwort](https://en.wikipedia.org/wiki/Spiderwort_(solitaire))
- [ ] [Spit+](https://en.wikipedia.org/wiki/Spit_(card_game))
- [ ] [Square](https://en.wikipedia.org/wiki/Square_(solitaire))
- [ ] [St. Helena](https://en.wikipedia.org/wiki/St._Helena_(solitaire))
- [ ] [Stalactites](https://en.wikipedia.org/wiki/Stalactites_(solitaire))
- [ ] [Stonewall](https://en.wikipedia.org/wiki/Stonewall_(solitaire))
- [ ] [Storehouse](https://en.wikipedia.org/wiki/Storehouse_(solitaire))
- [ ] [Strategy](https://en.wikipedia.org/wiki/Strategy_(solitaire))
- [ ] [Streets](https://en.wikipedia.org/wiki/Streets_(solitaire))
- [ ] [Streets and Alleys](https://en.wikipedia.org/wiki/Streets_and_Alleys)
- [ ] [Stronghold](https://en.wikipedia.org/wiki/Stronghold_(solitaire))
- [ ] [Sultan](https://en.wikipedia.org/wiki/Sultan_(solitaire))
- [ ] [Super Flower Garden](https://en.wikipedia.org/wiki/Super_Flower_Garden)
- [ ] [Superior Canfield](https://en.wikipedia.org/wiki/Superior_Canfield)
### T
- [ ] [Tableau](https://en.wikipedia.org/wiki/Tableau_(solitaire))
- [ ] [Take Fourteen](https://en.wikipedia.org/wiki/Take_Fourteen)
- [ ] [Tam O'Shanter](https://en.wikipedia.org/wiki/Tam_O%27Shanter_(solitaire))
- [ ] [Tens](https://en.wikipedia.org/wiki/Tens_(solitaire))
- [ ] [Terrace](https://en.wikipedia.org/wiki/Terrace_(solitaire))
- [ ] [The Clock](https://en.wikipedia.org/wiki/The_Clock_(solitaire))
- [ ] [The Fan](https://en.wikipedia.org/wiki/The_Fan_(solitaire))
- [ ] [The Plot](https://en.wikipedia.org/wiki/The_Plot_(solitaire))
- [ ] [Thirteens](https://en.wikipedia.org/wiki/Thirteens_(solitaire))
- [ ] [Thirteen Up](https://en.wikipedia.org/wiki/Thirteen_Up)
- [ ] [Thirteen Down](https://en.wikipedia.org/wiki/Thirteen_Down)
- [ ] [Three Blind Mice](https://en.wikipedia.org/wiki/Three_Blind_Mice_(solitaire))
- [ ] [Three Shuffles and a Draw](https://en.wikipedia.org/wiki/Three_Shuffles_and_a_Draw)
- [ ] [Thumb and Pouch](https://en.wikipedia.org/wiki/Thumb_and_Pouch)
- [ ] [Tournament](https://en.wikipedia.org/wiki/Tournament_(solitaire))
- [ ] [Tower of Hanoy (Tower of Hanoi)](https://en.wikipedia.org/wiki/Tower_of_Hanoi_(solitaire))
- [ ] [Tower of Pisa](https://en.wikipedia.org/wiki/Tower_of_Pisa_(solitaire))
- [x] [Travellers](https://en.wikipedia.org/wiki/Travellers_(solitaire))
- [ ] [Trefoil](https://en.wikipedia.org/wiki/Trefoil_(solitaire))
- [ ] [Triangle](https://en.wikipedia.org/wiki/Triangle_(solitaire))
- [ ] [Tri Peaks](https://en.wikipedia.org/wiki/Tri_Peaks_(game))
- [ ] [Tut's Tomb](https://en.wikipedia.org/wiki/Tut%27s_Tomb)
- [ ] [Twenty](https://en.wikipedia.org/wiki/Twenty_(solitaire))
### V
- [ ] [Vanishing Cross](https://en.wikipedia.org/wiki/Vanishing_Cross)
- [ ] [Vertical](https://en.wikipedia.org/wiki/Vertical_(solitaire))
- [ ] [Virginia Reel](https://en.wikipedia.org/wiki/Virginia_Reel_(solitaire))
### W
- [ ] [Washington's Favorite](https://en.wikipedia.org/wiki/Washington%27s_Favorite)
- [ ] [Wasp](https://en.wikipedia.org/wiki/Wasp_(solitaire))
- [ ] [Watch](https://en.wikipedia.org/wiki/Watch_(solitaire))
- [ ] [Weavers](https://en.wikipedia.org/wiki/Weavers_(solitaire))
- [x] [Westcliff (American)](https://en.wikipedia.org/wiki/Westcliff_(card_game))
- [x] [Westcliff (Classic)](https://en.wikipedia.org/wiki/Westcliff_(card_game))
- [x] [Whitehead](https://en.wikipedia.org/wiki/Whitehead_(solitaire))
- [ ] [Wildflower](https://en.wikipedia.org/wiki/Wildflower_(solitaire))
- [ ] [Will o' the Wisp](https://en.wikipedia.org/wiki/Will_o%27_the_Wisp_(solitaire))
- [ ] [Windmill](https://en.wikipedia.org/wiki/Windmill_(solitaire))
### Y
- [x] [Yukon](https://en.wikipedia.org/wiki/Yukon_(solitaire))
