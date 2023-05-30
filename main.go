package main

import (
	"net/http"
	"spade-7/Server"
)

func main() {
	s := Server.New()
	http.ListenAndServe(":8080", s)
}
