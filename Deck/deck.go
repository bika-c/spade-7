package Deck

import (
	"math/rand"
	"sort"
)

type Deck []Card

type Config struct {
	Jokers bool
	Faces  bool
}

var cards Deck = Deck{
	{ACE, SPADE},
	{TWO, SPADE},
	{THREE, SPADE},
	{FOUR, SPADE},
	{FIVE, SPADE},
	{SIX, SPADE},
	{SEVEN, SPADE},
	{EIGHT, SPADE},
	{NINE, SPADE},
	{TEN, SPADE},
	{JACK, SPADE},
	{QUEEN, SPADE},
	{KING, SPADE},

	{ACE, CLUB},
	{TWO, CLUB},
	{THREE, CLUB},
	{FOUR, CLUB},
	{FIVE, CLUB},
	{SIX, CLUB},
	{SEVEN, CLUB},
	{EIGHT, CLUB},
	{NINE, CLUB},
	{TEN, CLUB},
	{JACK, CLUB},
	{QUEEN, CLUB},
	{KING, CLUB},

	{ACE, HEART},
	{TWO, HEART},
	{THREE, HEART},
	{FOUR, HEART},
	{FIVE, HEART},
	{SIX, HEART},
	{SEVEN, HEART},
	{EIGHT, HEART},
	{NINE, HEART},
	{TEN, HEART},
	{JACK, HEART},
	{QUEEN, HEART},
	{KING, HEART},

	{ACE, DIAMOND},
	{TWO, DIAMOND},
	{THREE, DIAMOND},
	{FOUR, DIAMOND},
	{FIVE, DIAMOND},
	{SIX, DIAMOND},
	{SEVEN, DIAMOND},
	{EIGHT, DIAMOND},
	{NINE, DIAMOND},
	{TEN, DIAMOND},
	{JACK, DIAMOND},
	{QUEEN, DIAMOND},
	{KING, DIAMOND},

	{BLACKJOKER, UNDEFINED},
	{REDJOKER, UNDEFINED},
}

var cardsNoJokers Deck = Deck{
	{ACE, SPADE},
	{TWO, SPADE},
	{THREE, SPADE},
	{FOUR, SPADE},
	{FIVE, SPADE},
	{SIX, SPADE},
	{SEVEN, SPADE},
	{EIGHT, SPADE},
	{NINE, SPADE},
	{TEN, SPADE},
	{JACK, SPADE},
	{QUEEN, SPADE},
	{KING, SPADE},

	{ACE, CLUB},
	{TWO, CLUB},
	{THREE, CLUB},
	{FOUR, CLUB},
	{FIVE, CLUB},
	{SIX, CLUB},
	{SEVEN, CLUB},
	{EIGHT, CLUB},
	{NINE, CLUB},
	{TEN, CLUB},
	{JACK, CLUB},
	{QUEEN, CLUB},
	{KING, CLUB},

	{ACE, HEART},
	{TWO, HEART},
	{THREE, HEART},
	{FOUR, HEART},
	{FIVE, HEART},
	{SIX, HEART},
	{SEVEN, HEART},
	{EIGHT, HEART},
	{NINE, HEART},
	{TEN, HEART},
	{JACK, HEART},
	{QUEEN, HEART},
	{KING, HEART},

	{ACE, DIAMOND},
	{TWO, DIAMOND},
	{THREE, DIAMOND},
	{FOUR, DIAMOND},
	{FIVE, DIAMOND},
	{SIX, DIAMOND},
	{SEVEN, DIAMOND},
	{EIGHT, DIAMOND},
	{NINE, DIAMOND},
	{TEN, DIAMOND},
	{JACK, DIAMOND},
	{QUEEN, DIAMOND},
	{KING, DIAMOND},
}

// var pool = sync.Pool{
// 	New: func() any {
// 		return cards
// 	},
// }

/*
Get a deck of cards.
*/
func Get(c Config) Deck {
	if c.Jokers && c.Faces {
		return cards
	} else if c.Faces && !c.Jokers {
		return cardsNoJokers
	}
	var a Deck = make(Deck, 0, 54)
	for i := 1; i <= 4; i++ {
		for j := 1; j <= 13; j++ {
			if !c.Faces && j >= int(JACK) {
				continue
			}
			a.Add(Card{Rank(j), Suit(i)})
		}
	}
	if c.Jokers {
		a = append(a, Card{REDJOKER, UNDEFINED})
		a = append(a, Card{BLACKJOKER, UNDEFINED})
	}
	return a

}

// Put it back
// func (d *Deck) Free() {
// 	pool.Put(*d)
// }

func (d Deck) Shuffle() {
	rand.Shuffle(d.Len(), d.Swap)
}

func (a Deck) Has(c Card) bool {
	for i := 0; i < a.Len(); i++ {
		if a[i] == c {
			return true
		}
	}
	return false
}

func (a Deck) Intersection(b Deck) Deck {
	hash := make(map[Card]bool, a.Len())
	len := a.Len()
	if a.Len() > b.Len() {
		len = b.Len()
	}
	i := make([]Card, 0, len)
	for _, c := range a {
		hash[c] = true
	}

	for _, c := range b {
		if hash[c] {
			i = append(i, c)
		}
	}
	return i
}

func (d *Deck) Add(c ...Card) {
	*d = append(*d, c...)
}

func (d *Deck) RemoveAt(i int) {
	if i < 0 || i >= d.Len() {
		return
	}

	(*d)[i] = (*d)[len(*d)-1]
	(*d)[len(*d)-1] = Card{}
	*d = (*d)[:len(*d)-1]

}

func (d *Deck) Remove(remove ...Card) {
	count := len(remove)
	for i := d.Len() - 1; i >= 0; i-- {
		if count == 0 {
			return
		}
		c := (*d)[i]
		for _, rem := range remove {
			if c != rem {
				continue
			}
			count--
			d.RemoveAt(i)
		}
	}
}

func (d Deck) Random() (Card, int){
	r := rand.Intn(d.Len())
	return d[r], r
}

func (d Deck) Len() int {
	return len(d)
}

// Sort according to the rank and suit
func (d Deck) Less(i, j int) bool {
	return d[i].Less(d[j])
}

func (d Deck) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func (d Deck) Sort() {
	sort.Sort(d)
}
