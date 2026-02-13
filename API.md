
# API Documentation: Gomoku Game Server (42 Edition)

This API manages a single instance of a Gomoku game (Five in a Row) based on the **42 ruleset**, which includes stone captures. It supports local play, AI matches, or online multiplayer via an invitation system.

## 🔑 Authentication System

The API uses two types of tokens:

1. **Invitation Token (4 chars):** Generated for Player 2 during multiplayer game creation. It is used in the `POST /join` request.
2. **Session Token (16 chars):** Used to authenticate `POST /move`, `POST /ai-suggest`, and `POST /giveUp` requests.

---

## Table of Contents

* [GET /status](#get-status)
* [GET /board](#get-board)
* [POST /create](#post-create)
* [POST /giveUp](#post-giveup)
* [POST /join](#post-join)
* [POST /move](#post-move)
* [POST /ai-suggest](#post-ai-suggest)
* [POST /debug](#post-debug)

---

### `GET /status`

Checks if a game is already running or if an invitation is pending before attempting to create a new session.

**Response (200 OK):**

```json
{
  "goban_free": true,
  "pending_invitation": false
}

```

---

### `GET /board`

Returns the full state of the current game.

**Response (200 OK):**

* `board`: 19x19 matrix of integers (0: empty, 1: black, 2: white).
* `captured_b` / `captured_w`: Number of stones captured by each player.
* `turn`: Current turn number.
* `goban_free`: Server availability status.

```json
{
  "board": [
    [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ],
    [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ],
    [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ],
    [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ],
    [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ],
    [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ],
    [ 0, 0, 0, 0, 0, 0, 2, 0, 1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0 ],
    [ 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 2, 0, 1, 0, 0, 0, 0, 0, 0 ],
    [ 0, 0, 0, 0, 0, 0, 2, 0, 1, 0, 0, 0, 0, 0, 2, 0, 1, 0, 0 ],
    [ 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 2, 0, 1, 0, 0, 0, 0, 0, 0 ],
    [ 0, 0, 0, 0, 0, 0, 2, 0, 1, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0 ],
    [ 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 2, 0, 1, 0, 0, 0, 0, 0, 0 ],
    [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ],
    [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ],
    [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ],
    [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ],
    [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ],
    [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ],
    [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ]
  ],
  "captured_b": 3,
  "captured_w": 1,
  "goban_free": false,
  "turn": 34
}

```

---

### `POST /create`

Creates a new game session if the server is available.

**Request Body:**

* `ai_mode`: (bool) If `true`, `local_mode` is automatically set to `false`.
* `local_mode`: (bool) Toggle between local (same machine) or online play.

```json
{
  "ai_mode": false,
  "local_mode": false
}

```

**Response (200 OK):**

* `player_one`: Your session token (16 chars).
* `player_two`: Invitation token (4 chars) for your opponent, **OR** empty if AI/Local.

```json
{
  "ai_mode": false,
  "local_mode": false,
  "player_one": "b859573857a65df0",
  "player_two": "74e7"
}

```

---

### `POST /move`

Submits a move to the server. If AI mode is active, the AI will respond immediately.

**Coordinates:**

* The board coordinates start at **(0, 0)** and end at **(18, 18)**.

**Request Body:**

```json
{
  "x": 10,
  "y": 5,
  "token": "your_16_char_token"
}

```

**Success Response (200 OK):**

* `winner`: (Optional) Appears only when the game is over. Returns the player ID (**1** or **2**).
* `time_us`: Execution time in microseconds (provided when the AI plays).

```json
{
  "board": [[...]],
  "captured_b": 0,
  "captured_w": 0,
  "time_us": 22,
  "turn": 3,
  "winner": 2
}

```

**Error Responses:**

* **400 Bad Request:** Move is out of bounds (e.g., 19, 19) or spot is occupied.
```json
{ "error": "Illegal move" }

```


* **401 Unauthorized:** Invalid token or it's not your turn.
```json
{ "error": "Invalid token" }

```



---

### `POST /ai-suggest`

Asks the AI for the best move suggestion without executing it.

**Request Body:** `{"token": "your_16_char_token"}`

**Response (200 OK):**

```json
{
  "x": 5,
  "y": 5,
  "time_us": 1250
}

```

---

### `POST /giveUp`

Ends the game immediately, grants victory to the opponent, and frees the server.

**Request Body:** `{"token": "your_16_char_token"}`

**Response (200 OK):**

```json
{
  "message": "Game over."
}

```

---

### `POST /join`

Exchanges a 4-character invitation code for a 16-character session token.

**Request Body:** `{"token": "a1b2"}`

**Response (200 OK):**

```json
{
  "token": "5a51beedfede8cbe"
}

```

---

### `POST /debug`

Modifies the internal server state for testing scenarios.

**Request Body:**

* `reset_board`: (bool) Resets everything and frees the server.
* `board`: (Optional) Injects a custom 19x19 matrix.
* `captured_b` / `captured_w`: (Optional) Overrides capture counts.

**Response:** Returns the current board, turn, and both player tokens.

```json
{
  "board": [[...]],
  "captured_b": 0,
  "captured_w": 0,
  "token_player_one": "...",
  "token_player_two": "...",
  "turn": 1
}

```
