package Game

import (
	"net/http"
)

type GameServer interface {
	Game
	http.Handler

	Broadcast([]byte)
	JSONStat(Player) []byte
}
