# Spade7

The game starts with a board contains a single card, spade 7.

## ADDON

```json
{
   "decks": 1,
   "options": [], 
   "challenge": 12321
}
```

The addon has an array of cards representing all possible cards to play in the player's perspective given the current board. If the options is empty, then the `response.index` could be omitted to let the server to draw a card. Otherwise it must be a valid index: from 0 to len(previous player's card) -1. Challenge will only appear at when it's player's turn and will be updated each time.

## Player Response

The server accepts the response in following format:

```json
{
    "index": 0-53,
    "challenge": 100
}
```

- `index`: In response, `"index"` is the card index with respect to the server order that the player chose to play or draw.  Index can be omitted to let the server to pick one.
- `challenge`: In each broadcast, server will attach a challenge if it is the current player's turn. To verify a valid response, the client must reply to the server with the challenge. Example: if the server sends an addon:

```json
{
    "deck": 1,
    "options": [],
    "challenge": 9081
}
```

The client should response:

```json
{
    "challenge": 9081    
}
```

To indicate let the server to pick a card

### IMPORTANT

The order of the cards is not persevered upon removal but persevered upon addition. The algo:

> Given `[1,2,3,4,5]`, to remove number 3 at index 2: first, `[1,2,5,4,5]` then `[1,2,5,4]`  
> Given `[1,2,3,5]`, to add number 4: `[1,2,3,5,4]`
