package Spade7

import (
	"errors"
	"log"
	"spade-7/Deck"
	"spade-7/Game"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

type handle func(Deck.Card) bool

type Spade7 struct {
	id              Game.ID
	players         players
	expectedPlayers int
	curr            atomic.Int32

	cards Deck.Deck
	opts  Deck.Deck
	board Deck.Deck

	upgrader websocket.Upgrader
	l        sync.RWMutex
	logger   *log.Logger
}

func New(id Game.ID, l *log.Logger) *Spade7 {
	s := Spade7{
		id:              id,
		expectedPlayers: 0,
		board:           make(Deck.Deck, 0),
		upgrader: websocket.Upgrader{
			WriteBufferPool:   &sync.Pool{},
			EnableCompression: true,
		},
		l:      sync.RWMutex{},
		logger: l,
	}
	return &s
}

func (s *Spade7) Ended() bool {
	return s.board.Len() == s.cards.Len()
}

func (s *Spade7) start() error {
	if s.expectedPlayers >= 0 && s.players.Len() < s.expectedPlayers {
		return errors.New("Not Enough Player")
	}
	if s.board.Len() > 0 {
		return errors.New("Already Started")
	}
	s.Reset()
	s.BroadcastStat()
	go s.manager()
	return nil
}

// func (s *Spade7) check(t string) {
// 	for _, v := range s.players {
// 		if (v.Cards().Has(Deck.Card{Rank: 0, Suit: 0})) {
// 			panic("IMPOSSIBLE: " + t)
// 		}
// 	}
// }

func (s *Spade7) manager() {
	for s.Status() != "ended" {
		opt := s.opts.Intersection(s.current().Cards())
		c := s.current().readResponse(opt, s.previous().Cards(), 90*time.Second)
		s.next(s.options(*s.current()), c)
		s.BroadcastStat()
	}
}

var spade7 = Deck.Card{Rank: Deck.SEVEN, Suit: Deck.SPADE}
var initOpt = Deck.Deck{
	Deck.Card{Rank: Deck.SEVEN, Suit: Deck.HEART},
	Deck.Card{Rank: Deck.SEVEN, Suit: Deck.DIAMOND},
	Deck.Card{Rank: Deck.SEVEN, Suit: Deck.CLUB},
}

func (s *Spade7) checkSpade7(p int, d *Deck.Deck) {
	if s.board.Len() != 0 {
		return
	}
	for i := 0; i < d.Len(); i++ {
		if (*d)[i] == spade7 {
			d.RemoveAt(i)
			s.board.Add(spade7)
			s.curr.Store(int32(p))
		}
	}
}

func (s *Spade7) Reset() {
	s.l.Lock()
	defer s.l.Unlock()

	if s.players.Len() == 0 {
		return
	}

	s.cards.Shuffle()
	count := s.cards.Len() / s.players.Len()

	for i := range s.players {
		d := s.cards[(count * i) : (count*i)+count]
		s.players[i].Deck = make(Deck.Deck, len(d))
		copy(s.players[i].Deck, d)
	}
	for i := 0; i < s.cards.Len()%s.players.Len(); i++ {
		s.players[i].Deck = append(s.players[i].Deck, s.cards[count*len(s.players)+1])
	}
	s.initBoard()
}

func (s *Spade7) initBoard() {
	s.opts = s.opts[:0]
	s.opts = append(s.opts, Deck.Card{
		Rank: Deck.EIGHT,
		Suit: Deck.SPADE,
	}, Deck.Card{
		Rank: Deck.SIX,
		Suit: Deck.SPADE,
	})
	for i := 0; i < s.cards.Len()/52; i++ {
		s.opts = append(s.opts, initOpt...)
	}
	s.board.Add(spade7)
	for i := 0; i < s.players.Len(); i++ {
		if s.players[i].Deck.Has(spade7) {
			s.players[i].RemoveCards(spade7)
			s.curr.Store(int32(i))
		}
	}
}

func (s *Spade7) current() *player {
	return &s.players[s.curr.Load()]
}

func (s *Spade7) previous() *player {
	p := (int(s.curr.Load()) + (s.players.Len() - 1)) % len(s.players)
	for i := 2; s.players[p].Deck.Len() == 0; i++ {
		p = (int(s.curr.Load()) + (s.players.Len() - i)) % len(s.players)
	}
	return &s.players[p]
}

func (s *Spade7) String() string {
	return "Spade 7"
}

func (s *Spade7) next(h handle, c Deck.Card) {
	if h(c) {
		s.curr.Store(int32(int((s.curr.Load())+1) % s.players.Len()))
	}

	// return &s.players[s.curr.Load()]
}

func (s *Spade7) options(p player) handle {
	if len(p.Cards().Intersection(s.opts)) == 0 {
		return s.drawCard
	}
	return s.playCard
}

func (s *Spade7) Status() string {
	s.l.RLock()
	defer s.l.RUnlock()

	if s.Ended() {
		return "ended"
	} else if s.board.Len() > 0 {
		return "started"
	}
	return "pending"
}

func (s *Spade7) Players() (int, int) {
	return s.players.Len(), 0
}

func (s *Spade7) playCard(card Deck.Card) bool {
	s.current().RemoveCards(card)
	s.board = append(s.board, card)

	s.opts.Remove(card)

	if card.Rank < Deck.SEVEN && card.Rank > Deck.ACE {
		s.opts.Add(Deck.Card{Rank: card.Rank - 1, Suit: card.Suit})
	} else if card.Rank > Deck.SEVEN && card.Rank < Deck.KING {
		s.opts.Add(Deck.Card{Rank: card.Rank + 1, Suit: card.Suit})
	} else {
		s.opts.Add(Deck.Card{Rank: card.Rank + 1, Suit: card.Suit})
		s.opts.Add(Deck.Card{Rank: card.Rank - 1, Suit: card.Suit})
	}
	return true
}

func (s *Spade7) drawCard(c Deck.Card) bool {
	p := s.previous()
	if !c.Valid() {
		r, i := p.Deck.Random()
		s.current().AddCards(r)
		p.Deck.RemoveAt(i)
		return !s.opts.Has(c)
	}

	p.RemoveCards(c)
	s.current().AddCards(c)
	return !s.opts.Has(c)
}
