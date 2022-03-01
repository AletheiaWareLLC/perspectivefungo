perspectivefungo
================

https://perspective.fun

Perspective is a casual puzzle game in which players navigate a 3-dimensional maze controlling only their orientation and gravity.

- Swipe to rotate maze.
- Tap to activate gravity.

## Generate Puzzle

```sh
go run ./cmd/generator -size 11 -moves 7 -blocks 8 -portals 2 puzzle.json
```

## Play Puzzle

```sh
go run ./cmd/glfwplayer puzzle.json
```

## Build Web Player

```sh
tinygo build -o ./cmd/server/assets/static/player.wasm ./cmd/wasmplayer
```

## Run Web Server

```sh
go run ./cmd/server
```

Navigate to `localhost`
