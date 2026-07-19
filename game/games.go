package game

// AllVariants returns a slice of all registered game variants.
func AllVariants() []Variant {
	return []Variant{
		&Klondike{},
		&KlondikeVegas{},
		&Accordion{},
		&Addiction{},
		&AcesAndKings{},
		&AcesSquare{},
		&AcesUp{},
		&Acme{},
		&Agnes{},
		&Algerian{},
		&AmericanToad{},
		&Appreciate{},
		&Calculation{},
		&Canberra{},
		&Easthaven{},
		&FlowerGarden{},
		&FreeCell{},
		&Gaps{},
		&Kuipers{},
		&Russian{},
		&SirTommy{},
		&Somerset{},
		&Tasmanian{},
		&Travellers{},
		&WestcliffAmerican{},
		&WestcliffClassic{},
		&Whitehead{},
		&Yukon{},
	}
}
