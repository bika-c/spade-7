package util

import (
	"hash/fnv"
	"net/http"
	"spade-7/Game"
	"strconv"
	"unsafe"

	"github.com/julienschmidt/httprouter"
)

var algo = fnv.New32()

func HashID(s string) ([]byte, string) {
	algo.Reset()
	b := unsafe.StringData(s)
	algo.Write(unsafe.Slice(b, len(s)))

	var buffer [4]byte
	i := algo.Sum32()
	return algo.Sum(buffer[:]), strconv.FormatInt(int64(i), 10)
}

func ReadID(r *http.Request) (Game.ID, error) {
	id := r.Header.Get("ID")
	if id == "" {
		id = httprouter.ParamsFromContext(r.Context()).ByName("player")
	}
	i, e := strconv.ParseUint(id, 10, 32)
	return Game.ID(i), e
}
