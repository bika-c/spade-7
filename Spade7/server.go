package Spade7

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"spade-7/Deck"
	"spade-7/Game"
	"spade-7/Server/util"

	"github.com/bytedance/sonic"
)

type addon struct {
	NumberOfDeck int       `json:"decks"`
	Options      Deck.Deck `json:"options"`
}

type gameStat struct {
	P       players   `json:"players"`
	Board   Deck.Deck `json:"board"`
	Current int       `json:"current"`
	Id      Game.ID   `json:"id"`
	Status  string    `json:"status"`
	Player  player    `json:"player"`
	Addon   addon     `json:"addon"`
}

type config struct {
	NumberOfDeck    int    `json:"decks"`
	ExpectedPlayers int    `json:"expectedPlayers"`
	Status          string `json:"status"`
}

func (s *Spade7) AddPlayers(players ...player) {
	s.l.Lock()
	defer s.l.Unlock()

	for _, p := range players {
		s.players = append(s.players, p)
	}
}

func (s *Spade7) RemovePlayers(remove ...player) {
	s.l.Lock()
	defer s.l.Unlock()

	count := len(remove)
	for i := len(s.players); i >= 0; i-- {
		if count == 0 {
			return
		}
		c := s.players[i]
		for _, rem := range remove {
			if c.ID != rem.ID {
				continue
			}
			count--
			s.players[i] = s.players[len(s.players)-1]
			s.players[len(s.players)-1] = player{}
			s.players = s.players[:len(s.players)-1]
		}
	}
}

func (s *Spade7) newGameStat(p player) gameStat {
	return gameStat{
		s.players,
		s.board,
		int(s.curr.Load()),
		s.id,
		s.Status(),
		p,
		addon{
			s.cards.Len() / 52,
			s.opts.Intersection(p.Cards()),
		}}
}

func injectChallenge(r []byte, p player) []byte {
	root, e := sonic.Get(r)
	if e != nil {
		return []byte{}
	}

	player := root.Get("addon")
	player.SetAny("challenge", p.setKey())
	b, e := root.MarshalJSON()
	if e != nil {
		return []byte{}
	}
	return b
}

func (s *Spade7) JSONStat(p player) []byte {
	r, e := json.Marshal(s.newGameStat(p))
	if e != nil {
		return []byte{}
	}
	if s.current().ID != p.ID {
		return r
	}
	return injectChallenge(r, p)
}

func (s *Spade7) Broadcast(data []byte) {
	for i := range s.players {
		go s.players[i].Accept(data)
	}
}

func (s *Spade7) BroadcastStat() {
	for i := range s.players {
		go s.players[i].Accept(s.JSONStat(s.players[i]))
	}
}

func (s *Spade7) HasPlayer(i Game.ID) bool {
	return s.players.hasPlayer(i)
}

func (s *Spade7) parsePlayer(r *http.Request) (player, error) {
	b, e := io.ReadAll(r.Body)
	if e != nil {
		return player{}, e
	}
	var p player = newPlayer(s.logger)

	e = json.Unmarshal(b, &p)
	if e != nil {
		return player{}, e
	}
	if p.ID == 0 {
		return player{}, errors.New("Invalid")
	}
	return p, nil
}

func (s *Spade7) HandleJoin(w http.ResponseWriter, r *http.Request) {
	p, e := s.parsePlayer(r)
	if e != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	s.AddPlayers(p)
	s.logger.Println("Player", p.Name, ":", p.ID, "joined Spade7", s.id)
	s.BroadcastStat()
}

func (s *Spade7) Socket(id Game.ID, w http.ResponseWriter, r *http.Request) {
	p := s.players.get(id)
	if p.con != nil {
		p.disConnect()
	}

	c, e := s.upgrader.Upgrade(w, r, nil)
	if e != nil {
		s.logger.Println(id, "Upgrade failed", e.Error())
		return
	}
	p.con = c
	s.BroadcastStat()
	s.logger.Println("Player", p.Name, ":", p.ID, "connected Spade7", s.id)
}

func (s *Spade7) Control(w http.ResponseWriter, r *http.Request) {
	i, err := util.ReadID(r)
	if !s.players.hasPlayer(i) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if s.Status() == "started" {
		http.Error(w, "Game Started", http.StatusBadRequest)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var c config
	err = json.Unmarshal(b, &c)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if c.ExpectedPlayers > 0 && c.ExpectedPlayers >= s.players.Len() {
		s.expectedPlayers = c.ExpectedPlayers
	}
	if c.NumberOfDeck > 0 {
		s.cards = make(Deck.Deck, 0, c.NumberOfDeck*52)
		for i := 0; i < c.NumberOfDeck; i++ {
			s.cards.Add(Deck.Get(Deck.Config{Jokers: false, Faces: true})...)
		}
	}
	switch c.Status {
	case "start":
		err = s.start()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		break
	}
	s.logger.Println("Game: ", s.id, "config updated with ", s.cards.Len(), "cards")
	s.BroadcastStat()
}

func (s *Spade7) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// switch r.Method {
	// case http.MethodPatch:
	// 	s.join(w, r)
	// 	break
	// default:
	// 	s.socket(w, r)
	// }
}
