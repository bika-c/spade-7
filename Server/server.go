package Server

import (
	"io"
	"math/rand"
	"net/http"
	"spade-7/Game"
	"spade-7/Spade7"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

const (
	spade7 int = iota
)

func toTypeInt(s string) int {
	switch s {
	case "spade7":
		return spade7
	default:
		return -1
	}
}

type Server struct {
	games  map[int]map[int]Game.GameServer
	router *httprouter.Router
}

func (s Server) list(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "listing games")
}

func (s Server) new(w http.ResponseWriter, r *http.Request) {
	p := httprouter.ParamsFromContext(r.Context())

	t := toTypeInt(p.ByName("type"))

	switch t {
	case spade7:
		if _, ok := s.games[spade7]; !ok {
			s.games[spade7] = make(map[int]Game.GameServer, 10)
		}
		id := rand.Int()
		s.games[spade7][id] = Spade7.New(Game.ID(id))
		io.WriteString(w, "Created with ID: "+strconv.Itoa(id)) //xxx
		break
	}
}

func (s Server) join(w http.ResponseWriter, r *http.Request) {
	p := httprouter.ParamsFromContext(r.Context())

	t := toTypeInt(p.ByName("type"))
	id, e := strconv.Atoi(p.ByName("id"))
	if e != nil {
		return
	}

	io.WriteString(w, "joined the game: "+strconv.Itoa(id)) // todo
	s.games[t][id].ServeHTTP(w, r)
}

func New() *Server {
	r := httprouter.New()
	s := Server{
		router: r,
		games:  make(map[int]map[int]Game.GameServer, 3),
	}
	r.HandlerFunc("GET", "/:type", s.list)
	r.HandlerFunc("PATCH", "/:type/:id", s.join)
	r.HandlerFunc("POST", "/:type/", s.new)
	return &s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
