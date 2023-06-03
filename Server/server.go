package Server

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"spade-7/Game"
	"spade-7/Server/util"
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
	games  map[int]map[int]Game.Server
	router *httprouter.Router
	logger *log.Logger
}

func (s Server) list(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "listing games:\n")
}

func (s Server) newID(t int) int {
	id := rand.Int()
	for _, ok := s.games[t][id]; ok; {
		id += rand.Int()
	}
	return id
}

func (s Server) new(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/player":
		s.player(w, r)
		break
	default:
		s.game(w, r)
	}
}

func (s Server) socket(w http.ResponseWriter, r *http.Request) {
	p := httprouter.ParamsFromContext(r.Context())
	t := toTypeInt(p.ByName("type"))
	g, e := strconv.Atoi(p.ByName("id"))

	if e != nil || !s.check(g, t) {
		http.Error(w, "Game Not Supported or Game does not Exist", http.StatusBadRequest)
		return
	}

	player, e := strconv.Atoi(p.ByName("player"))
	if e != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	if !s.games[t][g].HasPlayer(Game.ID(player)) {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	s.games[t][g].Socket(Game.ID(player), w, r)
}

func (s Server) control(w http.ResponseWriter, r *http.Request) {
	p := httprouter.ParamsFromContext(r.Context())
	t := toTypeInt(p.ByName("type"))
	g, e := strconv.Atoi(p.ByName("id"))

	if e != nil || !s.check(g, t) {
		http.Error(w, "Game Not Supported or Game does not Exist", http.StatusBadRequest)
		return
	}

	s.games[t][g].Control(w, r)
}

func (s Server) check(g, t int) bool {
	_, ok := s.games[t]
	if !ok {
		return false
	}
	_, ok = s.games[t][g]
	return ok
}

func (s Server) join(w http.ResponseWriter, r *http.Request) {
	p := httprouter.ParamsFromContext(r.Context())
	t := toTypeInt(p.ByName("type"))
	g, e := strconv.Atoi(p.ByName("id"))

	if e != nil || !s.check(g, t) {
		http.Error(w, "Game Not Supported or Game does not Exist", http.StatusBadRequest)
		return
	}

	s.games[t][g].HandleJoin(w, r)
}

func (s Server) game(w http.ResponseWriter, r *http.Request) {
	p := httprouter.ParamsFromContext(r.Context())
	t := toTypeInt(p.ByName("type"))

	switch t {
	case spade7:
		if _, ok := s.games[spade7]; !ok {
			s.games[spade7] = make(map[int]Game.Server, 10)
		}
		id := s.newID(spade7)
		s.games[spade7][id] = Spade7.New(Game.ID(id), s.logger)
		s.logger.Println("Game Spade7:", id, "created")
		break
	default:
		http.Error(w, "Game Not Supported", http.StatusBadRequest)
		return
	}
}

func (Server) player(w http.ResponseWriter, r *http.Request) {
	i := r.Header.Get("PLAYER-ID")
	b, s := util.HashID(i)
	w.Header().Add("ID", s)
	w.Write(b)
}

func New() *Server {
	r := httprouter.New()
	s := Server{
		router: r,
		logger: log.Default(),
		games:  make(map[int]map[int]Game.Server, 3),
	}
	r.HandlerFunc("GET", "/:type", s.list)
	r.HandlerFunc("POST", "/:type/:id", s.join)
	r.HandlerFunc("PATCH", "/:type/:id", s.control)
	r.HandlerFunc("GET", "/:type/:id/:player", s.socket)
	r.HandlerFunc("POST", "/:type", s.new)
	// todo
	r.HandlerFunc("GET", "/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "ws.html")
	})
	return &s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
