# Gomoku AI

School project aiming to implement an AI capable of playing Gomoku (19x19) with:

- Full rule support (captures, forbidden double-threes, win by alignment or by captures),
- A Minimax-style algorithm (or variant) with a custom evaluation heuristic,
- A user interface allowing matches:
  - human vs AI,
  - or human vs human.

The long-term goal is to cleanly separate the game engine + AI from the user interface, so the same core can be reused with different frontends (terminal UI, web UI, or other integrations).
