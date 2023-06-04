package Spade7

import (
	"log"
	"math/rand"
	"spade-7/Deck"
	"spade-7/Game"
	"sync"

	"time"

	"github.com/bytedance/sonic"

	"github.com/gorilla/websocket"
)

var ast = sonic.ConfigStd

type player struct {
	Name      string    `json:"name"`
	ID        Game.ID   `json:"id"`
	Deck      Deck.Deck `json:"cards"`
	challenge int
	con       *websocket.Conn
	logger    *log.Logger
	l         *sync.RWMutex
}

type response struct {
	Index     int `json:"index"`
	Challenge int `json:"challenge"`
}

func newPlayer(l *log.Logger) player {
	p := player{
		con:       nil,
		Deck:      nil,
		challenge: rand.Intn(100000),
		logger:    l,
		l:         &sync.RWMutex{},
	}
	return p
}

func (p *player) readResponse(opt Deck.Deck, pre Deck.Deck, t time.Duration) Deck.Card {
	if p.con == nil {
		return p.autoPick(opt, pre)
	}

	c := make(chan Deck.Card)
	defer close(c)
	go p.read(opt, pre, c)
	select {
	case <-time.After(t):
		p.logger.Println("Player", p.Name, p.ID, "timeout")
		p.disConnect()
	case card := <-c:
		return card
	}
	return p.autoPick(opt, pre)
}

func (p *player) setKey() int {
	p.challenge = rand.Intn(100000)
	return p.challenge
}

func (p *player) autoPick(opt Deck.Deck, pre Deck.Deck) Deck.Card {
	if opt.Len() == 0 {
		c, _ := pre.Random()
		return c
	}
	return opt[0]
}

func (p *player) valid(opt Deck.Deck, pre Deck.Deck, i int) bool {
	if i < 0 {
		return false
	}
	if opt.Len() > 0 && i >= opt.Len() {
		return false
	}
	if opt.Len() == 0 && i >= pre.Len() {
		return false
	}
	return true
}

func (p *player) read(opt Deck.Deck, pre Deck.Deck, c chan Deck.Card) {
	var r response
	for r.Challenge != p.challenge || !p.valid(opt, pre, r.Index) {
		e := p.con.ReadJSON(&r)
		p.disConnect()
		if e != nil {
			return
		}
	}
	if opt.Len() == 0 {
		c <- pre[r.Index]
	} else {
		c <- opt[r.Index]
	}
}

// func (d digest) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(d)
// }

func (p *player) Accept(d []byte) {
	if p.con == nil {
		return
	}
	p.l.Lock()
	defer p.l.Unlock()
	p.con.WriteMessage(websocket.TextMessage, d)
}

// func (p *player) Connect(c net.Conn, g Game.Game) {

// 	var s *Spade7 = g.(*Spade7)
// 	p.disConnect()

// 	go func(p *player, g *Spade7) {
// 		defer p.disConnect()
// 		b := make([]byte, 0, 4)
// 		g.AddPlayers(p)
// 		defer g.RemovePlayers(p)

// 		for n, e := c.Read(b); e != nil; n, e = c.Read(b) {
// 			fmt.Printf("n: %v\n", n)
// 		}

// 	}(p, s)
// }

// func (p *player) Name() string {
// 	return p.info.Name
// }

// func (p *player) ID() Game.ID {
// 	return p.info.ID
// }

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

func (p *player) json(cards bool) ([]byte, error) {
	type a player
	var x a = a(*p)
	s, e := ast.MarshalToString(x)
	if e != nil {
		return []byte{}, e
	}
	t, _ := sonic.GetFromString(s)
	t.SetAny("online", p.con != nil)
	if !cards {
		t.SetAny("cards", p.Deck.Len())
	}
	return t.MarshalJSON()
}

func (p *player) MarshalJSON() ([]byte, error) {
	return p.json(false)
}

func (p *player) disConnect() {
	if p.con != nil {
		// p.con.WriteMessage(websocket.CloseMessage, nil)
		// p.con.Close()
		// p.con = nil
		p.logger.Println("Player", p.Name, p.ID, "disconnected")
	}
}

type players []player

func (p players) Len() int {
	return len(p)
}

func (p players) get(id Game.ID) *player {

	for i := 0; i < len(p); i++ {
		if p[i].ID == id {
			return &p[i]
		}
	}

	return nil
}

func (p players) hasPlayer(id Game.ID) bool {
	for i := 0; i < len(p); i++ {
		if p[i].ID == id {
			return true
		}
	}
	return false
}

// func (p players) MarshalJSON() ([]byte, error) {
// 	return ast.Marshal(p)
// }
