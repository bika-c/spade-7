package Spade7

import (
	"math/rand"
	"spade-7/Deck"
	"spade-7/Game"
	"sync/atomic"
)

type Spade7 struct {
	players players
	options Deck.Deck
	board   Deck.Deck

	drawCardHandler []Game.Handler
	playCardHandler []Game.Handler

	// todo get rid off this and merge with d Deck.Deck
	numOfDeck int
	current   int32
	id        Game.ID
}

func New(id Game.ID) *Spade7 {
	s := Spade7{
		id: id,
	}
	s.drawCardHandler = []Game.Handler{Game.Handle(s.drawCard)}
	s.playCardHandler = []Game.Handler{Game.Handle(s.playCard)}
	s.Reset()
	return &s
}

func (s *Spade7) playCard(c ...Deck.Card) bool {
	if len(c) != 0 {
		panic("One card at a time")
	}
	card := c[0]
	s.Current().RemoveCards(card)
	s.board = append(s.board, card)

	s.options.Remove(card)

	if card.Rank <= Deck.SEVEN && card.Rank > Deck.ACE {
		s.options.Add(Deck.Card{Rank: card.Rank - 1, Suit: card.Suit})
	} else if card.Rank > Deck.SEVEN && card.Rank < Deck.KING {
		s.options.Add(Deck.Card{Rank: card.Rank + 1, Suit: card.Suit})
	}
	return true
}

func (s *Spade7) drawCard(c ...Deck.Card) bool {
	p := (int(s.current) + len(s.players) - 1) % len(s.players)
	r := rand.Intn(len(s.players[p].Cards()))
	s.Current().AddCards(s.players[p].Cards()[r])
	return len(s.Current().Cards().Intersection(s.options)) != 0
}

func (s *Spade7) Ended() bool {
	return s.board.Len() == s.numOfDeck*52
}

func (s *Spade7) Reset() {
	if len(s.players) == 0 {
		return
	}
	// todo optimize this to a local var
	d := make(Deck.Deck, 0, s.numOfDeck*52)
	for i := 0; i < s.numOfDeck; i++ {
		d = append(d, Deck.Get(Deck.Config{Jokers: false, Faces: true})...)
	}
	d.Shuffle()
	count := d.Len() / len(s.players)

	// todo check each player gets a new copy of the slices of cards
	for i, p := range s.players {
		p.Deck = d[(count * i) : (count*i)+count]
	}
	for i := 0; i < d.Len()%len(s.players); i++ {
		s.players[i].Deck = append(s.players[i].Deck, d[count*len(s.players)+1])
	}
	s.options = s.options[:0]
	s.board = s.board[:0]
	s.current = 0
}

func (s *Spade7) Current() Game.Player {
	return &s.players[atomic.LoadInt32(&s.current)]
}

func (s *Spade7) String() string {
	return "Spade 7"
}

func (s *Spade7) Next(m Game.Handler, c ...Deck.Card) Game.Player {
	if m.Handle(c...) {
		atomic.SwapInt32(&s.current, int32(int(atomic.LoadInt32(&s.current))+1%len(s.players)))
	}

	return &s.players[s.current]
}

func (s *Spade7) Options(p Game.Player) []Game.Handler {
	if len(p.Cards().Intersection(s.options)) == 0 {
		return s.drawCardHandler
	}
	return s.playCardHandler
}

func (s *Spade7) Players() []Game.Player {
	r := make([]Game.Player, 0, len(s.players))
	for i := 0; i < len(s.players); i++ {
		r = append(r, &s.players[i])
	}
	return r
}

func (s *Spade7) AddPlayers(players ...Game.Player) {
	for _, p := range players {
		s.players = append(s.players, *p.(*player))
	}
}

func (s *Spade7) RemovePlayers(remove ...Game.Player) {
	count := len(remove)
	for i := len(s.players); i >= 0; i-- {
		if count == 0 {
			return
		}
		c := s.players[i]
		for _, rem := range remove {
			if c.ID() != rem.ID() {
				continue
			}
			count--
			s.players[i] = s.players[len(s.players)-1]
			s.players[len(s.players)-1] = player{}
			s.players = s.players[:len(s.players)-1]
		}
	}
}

func (s *Spade7) Status() string {
	var status string = "ready"
	if s.board.Len() == 1 {
		status = "started"
	} else if s.Ended() {
		status = "ended"
	} else {
		status = "running"
	}
	return status
}
