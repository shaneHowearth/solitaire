package game

type Category int

const (
	CatUnknown Category = iota
	CatKlondike
	CatSpider
	CatPairing
	CatFoundation
	CatSpecialty
)

// String returns the display name for the UI
func (c Category) String() string {
	switch c {
	case CatKlondike:
		return "Klondike Family"
	case CatSpider:
		return "Spider & Yukon"
	case CatPairing:
		return "Adding & Pairing"
	case CatFoundation:
		return "Foundation Builders"
	case CatSpecialty:
		return "Specialty Games"
	default:
		return "Other Variants"
	}
}
