package Deck

type Suit uint8
type Rank uint8

const (
	UNDEFINED Suit = iota
	SPADE
	CLUB
	HEART
	DIAMOND
)

func (s Suit) String() string {
	switch s {
	case SPADE:
		return "Spade"
	case CLUB:
		return "Club"
	case HEART:
		return "Heart"
	case DIAMOND:
		return "Diamond"
	default:
		return ""
	}
}

const (
	ACE Rank = iota + 1
	TWO
	THREE
	FOUR
	FIVE
	SIX
	SEVEN
	EIGHT
	NINE
	TEN
	JACK
	QUEEN
	KING
	BLACKJOKER
	REDJOKER
)

func (r Rank) String() string {
	switch r {
	case ACE:
		return "A"
	case TWO:
		return "2"
	case THREE:
		return "3"
	case FOUR:
		return "4"
	case FIVE:
		return "5"
	case SIX:
		return "6"
	case SEVEN:
		return "7"
	case EIGHT:
		return "8"
	case NINE:
		return "9"
	case TEN:
		return "10"
	case JACK:
		return "J"
	case QUEEN:
		return "Q"
	case KING:
		return "K"
	case BLACKJOKER:
		return "Black Joker"
	case REDJOKER:
		return "Red Joker"
	default:
		return ""
	}
}

var Invalid = Card{
	99, 100,
}

func (c Card) Valid() bool {
	return c.Rank >= 1 && c.Rank <= 4 && c.Suit >= 1 && c.Suit <= 15
}

type Card struct {
	Rank Rank `json:"rank"`
	Suit Suit `json:"suit"`
}

func (a Card) Less(b Card) bool {
	if a.Suit == UNDEFINED || a.Rank > 13 {
		return false
	}
	return a.Rank < b.Rank || a.Suit < b.Suit
}
