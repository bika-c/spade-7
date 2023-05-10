package Spade7

import (
	"spade-7/Deck"
	"testing"
)

func BenchmarkSpade7_Options(b *testing.B) {
	p := player{
		name:  "abc",
		cards: Deck.Deck{},
	}
	s := Spade7{
		options: Deck.Deck{
			{Rank: Deck.SEVEN, Suit: Deck.SPADE},
			{Rank: Deck.EIGHT, Suit: Deck.SPADE},
		},
	}
	has := player{
		name: "abc",
		cards: Deck.Deck{
			{Rank: Deck.SEVEN, Suit: Deck.SPADE},
		},
	}
	for n := 0; n < b.N; n++ {
		if n%2 == 0 {
			a := s.Options(&p)
			_ = len(a)
		} else {
			a := s.Options(&has)
			_ = len(a)
		}
	}
}
