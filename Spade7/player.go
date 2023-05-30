package Spade7

import (
	"encoding/json"
	"fmt"
	"net"
	"spade-7/Deck"
	"spade-7/Game"
)

type player struct {
	info
	Deck Deck.Deck `json:"Cards"`
	con  net.Conn
}

type info struct {
	Name string  `json:"name"`
	ID         Game.ID `json:"id"`
}

type digest struct {
	info
	Cards byte `json:"cards"`
}

func (d digest) MarshalJSON() ([]byte, error) {
	return json.Marshal(d)
}

func (p *player) Accept(d []byte) {
	p.con.Write(d)
}

func (p *player) Connect(c net.Conn, g Game.Game) {

	var s *Spade7 = g.(*Spade7)
	p.disConnect()
	p.con = c

	go func(p *player, g *Spade7) {
		defer p.disConnect()
		b := make([]byte, 0, 4)
		g.AddPlayers(p)
		defer g.RemovePlayers(p)

		for n, e := c.Read(b); e != nil; n, e = c.Read(b) {
			fmt.Printf("n: %v\n", n)
		}

	}(p, s)
}

func (p *player) Name() string {
	return p.info.Name
}

func (p *player) ID() Game.ID {
	return p.info.ID
}

func (p *player) Cards() Deck.Deck {
	return p.Deck
}

func (p *player) AddCards(c ...Deck.Card) {
	p.Deck.Add(c...)
}

func (p *player) RemoveCards(remove ...Deck.Card) {
	p.Deck.Remove(remove...)
}

func (p *player) Remove(i int) {
	p.Deck[i] = p.Deck[len(p.Deck)-1]
	p.Deck[len(p.Deck)-1] = Deck.Card{}
	p.Deck = p.Deck[:len(p.Deck)-1]
}

func (p *player) json(cards bool) json.Marshaler {
	if cards {
		return p
	}
	return digest{
		p.info,
		byte(p.Deck.Len()),
	}
}

func (p *player) MarshalJSON() ([]byte, error) {
	return json.Marshal(p)
}

func (p *player) disConnect() {
	if p.con != nil {
		p.con.Close()
		p.con = nil
	}
}

type players []player

func (p players) hasPlayer(id Game.ID) bool {
	for i := 0; i < len(p); i++ {
		if p[i].info.ID == id {
			return true
		}
	}
	return false
}

func (p players) MarshalJSON() ([]byte, error) {
	r := make([]json.Marshaler, len(p))
	for i := 0; i < len(p); i++ {
		r = append(r, p[i].json(false))
	}
	return json.Marshal(r)
}