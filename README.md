# Gomoku
Gomoku is a strategic board game project where you compete against an AI or other players by placing stones on a Goban (board).

## Setup
This project is built using a Client-Server architecture. To play, you must first set up the server.
```bash
make nix # Setup the terminal with the software of the version we want
make gomoku # Launch the server
```
Once the server is running, you can connect via the API (see the [API Documentation](/API.md))

If you want to use the Ncurses client, run:

```bash
make ncurses
```

## Rules

The game is turn-based. The goal is to align 5 stones or capture 10 stones. However, this implementation includes specific variation rules:

### Capture
You can capture (remove) opponent stones by flanking a pair of them.

* **Mechanism:** If the opponent has two stones together (`O O`) and you surround them with yours (`X O O X`), the opponent's stones are removed.
* **Restriction:** You can only capture **pairs**. You cannot capture single stones or lines of 3 or more.

### Double-threes

The move known as a `double-three` is forbidden.

* **Free-three:** An alignment of three stones that, if not immediately blocked, allows for an indefensible alignment of four stones.
* **Restriction:** You cannot place a stone that creates **two simultaneous free-three alignments**. This is because a double-three is impossible to block effectively.


## Bonus Features

We implemented several advanced features for this project:

1. [Client-Server Architecture](#client-server)
2. [AI Assistance](#ai-suggest)
3. [Remote Play](#remote-play)
4. [Local Play](#local-play)
5. [Multiple Client](#multiple-client)


### Client Server

The project is coded as a Client-Server application, allowing gameplay via a API.

### AI Suggest

The API includes an `ai-suggest` endpoint that analyzes the board and recommends the best next move.

### Remote Play
You can play on multiplayer, when creating a remote room, it give you a code to invite your friend. (see the [API Docs](/API.md))

### Local Play

Support for "Hot-seat" multiplayer, allowing two players to play on the same machine (see the [API Docs](/API.md))

### Multiple Client

We developed three distinct front-ends:

1. **[Debug Client](/debug-client/):** A bare-bones terminal client. It is functional, primarily used for testing, and provides a raw view of the game state.

  ![debug-client image](/images/debug-client.png)

2. **[Ncurses Client](/client-ncurses/):** A terminal-based UI using `ncurses`. It features a proper windowed interface, supporting both arrow keys and mouse clicks to place stones.

  ![ncurses client image](/images/ncurses.png)

3. **[Minecraft Client](/MC-Client/):** A fully immersive client-side Minecraft mod. It connects to any standard Gomoku server and allows you to play the game within Minecraft.

  ![minecraft client image](/images/mc-client.png)