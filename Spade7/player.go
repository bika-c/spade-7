package Spade7

import (
	"fmt"
	"net"
	"spade-7/Deck"
	"spade-7/Game"
)

type player struct {
	name  string
	id    Game.ID
	cards Deck.Deck
}

func (p *player) Connect(c net.Conn, g Game.Game) {

	go func(p *player, c net.Conn, g Game.Game) {
		defer c.Close()
		b := make([]byte, 0, 4)
		g.AddPlayers(p)
		defer g.RemovePlayers(p)

		for n, e := c.Read(b); e != nil; n, e = c.Read(b) {
			fmt.Printf("n: %v\n", n)
		}

	}(p, c, g)
}

func (p *player) Name() string {
	return p.name
}

func (p *player) ID() Game.ID {
	return p.id
}

// func (p *player) Play(c ...Deck.Card) {
// 	p.selected = c
// }

func (p *player) Cards() Deck.Deck {
	return p.cards
}

func (p *player) AddCards(c ...Deck.Card) {
	p.cards.Add(c...)
}

func (p *player) RemoveCards(remove ...Deck.Card) {
	p.cards.Remove(remove...)
}

func (p *player) Remove(i int) {
	p.cards[i] = p.cards[len(p.cards)-1]
	p.cards[len(p.cards)-1] = Deck.Card{}
	p.cards = p.cards[:len(p.cards)-1]
}
