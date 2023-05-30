package Spade7

import (
	"encoding/json"
	"io"
	"net/http"
	"spade-7/Deck"
	"spade-7/Game"
	"sync/atomic"
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
	Player  *player   `json:"player"`
	Addon   addon     `json:"addon"`
}

func (g gameStat) MarshalJSON() ([]byte, error) {
	return json.Marshal(g)
}

func (s *Spade7) JSONStat(p Game.Player) []byte {
	r, e := gameStat{
		s.players,
		s.board,
		int(atomic.LoadInt32(&s.current)),
		s.id,
		s.Status(),
		p.(*player),
		addon{
			s.numOfDeck,
			s.options,
		}}.MarshalJSON()

	if e != nil {
		return []byte{}
	}

	return r
}

func (s *Spade7) Broadcast(data []byte) {
	for i := range s.players {
		go s.players[i].Accept(data)
	}
}

func (s *Spade7) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "\nSpade7 here\n")
}
