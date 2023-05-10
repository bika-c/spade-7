package Game

import (
	"spade-7/Deck"
)

type ID uint64

type Handle func(...Deck.Card) bool

func (h Handle) Handle(c ...Deck.Card) bool {
	return h(c...)
}

// func (h Handle) String() string {
// 	str := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
// 	if len(str) > 3 && str[len(str)-3:] == "-fm" {
// 		return str[strings.LastIndex(str, ".")+1 : len(str)-3]
// 	}
// 	return str[strings.LastIndex(str, ".")+1:]
// }

type Handler interface {
	Handle(...Deck.Card) bool
}

// Game is an abstraction of card games
type Game interface {
	// Reset. Equivalent to restart a game
	Reset()
	// Returns true when the game is ended
	Ended() bool
	/*
		Advance the game.

		Handler is performed by the current player.
		Return the next player to play the round
	*/
	Next(Handler, ...Deck.Card) Player
	/*
		Gives all the possible actions for the player
		i.e. pass or take card
	*/
	Options(Player) []Handler

	Current() Player
	Players() []Player
	AddPlayers(p ...Player)
	RemovePlayers(p ...Player)

	Broadcast()
}
