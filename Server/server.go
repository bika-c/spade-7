package server

import (
	"fmt"
	"net/http"
	"spade-7/Game"
)

type Server struct {
	games map[Game.ID]Game.Game
}

func (s *Server) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "spade-7":
		fmt.Println("Handle spade7")
		break
	}
}
