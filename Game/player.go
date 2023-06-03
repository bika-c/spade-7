package Game

type Player interface {
	// Play accepts the player's choice of cards to play
	// Play(...Deck.Card)

	// Connect(net.Conn, Game)
	Accept(data []byte)

	// Broadcast(data[]byte)

	// Cards() Deck.Deck
	// AddCards(...Deck.Card)
	// RemoveCards(...Deck.Card)
	// Remove(int)

	// ID() ID
	// Name() string
}
