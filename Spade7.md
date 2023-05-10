# Spade7

The game starts with a board contains a single card, spade 7.

## ADDON

The addon is an array of cards representing all possible cards to play in the player's perspective given the current board. If the addon is empty array or nil, the response must be `"move": 1` and `"move": 0` otherwise.

## Player Response

The server accepts the response in following format:

```json
{
    "move": 0|1,
    "index": 0-53
}
```

### Definition

In response, `0` represents the 'play a card' move and the `"index"` is the card index with respect to the server order that the player chose to play. `1` represents the `draw a card` move and if `index not in the range [0, len(pre.cards)]` or it is omitted then the server will pick one.

### IMPORTANT

The order of the cards is not persevered upon removal but persevered upon addition. The algo:

> Given `[1,2,3,4,5]`, to remove number 3 at index 2: first, `[1,2,5,4,5]` then `[1,2,5,4]`  
> Given `[1,2,3,5]`, to add number 4: `[1,2,3,5,4]`
