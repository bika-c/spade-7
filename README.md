# API Doc

------------------------------------------------------------------------------------------

## Player

### `POST:` `/player`

> Register a player given an unique identification to the player. Returns a server id

- Header: `Player-ID: <string>` (must be unique)
- Parameters: NONE
- Body: NONE

#### Response

- Status: `200`
- Header:
    - `Content-Type`: `application/octet-stream`
    - `ID`: `"ABCK`
- Body: a big endian `uint32` (subject to change)

    ```binary
    {uint32} 
    ```

## Game

> `{type}` denotes the name of the game i.e. spade7

### `GET:` `/{type}`

Not yet implemented

> List all games

- Parameters: NONE
- Body: NONE

#### Response

- Status: `200`
- Header:
    - `Content-Type`: `application/json`
- Body:

    ```json
    {
        "games": [
            {
                "id": 1234,
                "players": 4,
                "expected": 5,
                "status": "started"|"ended"|"pending"
            }
        ]
    }
    ```

### `POST` `/{type}`

> Create a game returns the id of the game created

- Parameters: NONE
- Header: NONE
- Body: A json object defined by game

#### Response

- Status: `200`
- Header:
    - `Content-Type`: `application/json`
    - `ID`: `1234`

### `POST` `/{type}/{id}`

> Join the game.  Must carry the server id of the player
> To use websocket, this method must be invoked before dialing

- Parameters: NONE
- Body: NONE

#### Response

- Status: `200`

### `PATCH` `/{type}/{id}`

> Modify the game

- Parameters: NONE
- Header: `ID`: 12321
- Body: A json object defined by game

#### Response

- Status: `200`

## Card Game Websocket Protocol

- All messages are communicated using JSON (client to server and server to json)

- Game stat is defined as

    ```json
    {
        "id": 123,
        "status": "ended"|"started"|"pending",
        "players": [
            {
                "id": 123,
                "name": "abc",
                "cards": 54, 
            } // player def
        ],
        "current": players.index,
        "player": {
            "id": 123,
            "name": "abc",
            "cards": [
                {1,1} // card def
            ]
        },
        "board": [
            // cards
        ],
        "addon": any // any other supportive info goes here. Refer to game implementation
    }
    ```

- At the event of new player joining the game, a broadcast of game stat will be sent

### Running

- At the event of game start, a game stat will be broadcasted
- At each player's turn, the protocol accepts the first valid json object (Refer to game implementation). An update of game stat response will be broadcasted.
- If it is not the player's turn, all sent messages are undefined.
- The order of the `stat.players` is implementation defined but should stay still. `stat.current` is the index of the array
