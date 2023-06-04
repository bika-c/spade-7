package Game

import (
	"net/http"
)

type Server interface {
	Game
	HandleJoin(http.ResponseWriter, *http.Request)
	Socket(ID, http.ResponseWriter, *http.Request)
	Control(http.ResponseWriter, *http.Request)
	HasPlayer(ID) bool
}
