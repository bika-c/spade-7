package Game

import (
	"net/http"
)

type Server interface {
	Game
	http.Handler

	HandleJoin(http.ResponseWriter, *http.Request)
	Socket(ID, http.ResponseWriter, *http.Request)
	Control(http.ResponseWriter, *http.Request)
	HasPlayer(ID) bool

	// Broadcast([]byte)
	// JSONStat(Player) []byte

	// AddPlayers(p ...Player)
	// RemovePlayers(p ...Player)
}
