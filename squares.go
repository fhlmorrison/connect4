package main

import (
	"crypto/rand"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"sync"
)

type Game struct {
	Id            string
	Board         Board
	CurrentPlayer Tile
}

func (g *Game) GetBoard() Board {
	return g.Board
}

func (g *Game) Reset() {
	g.Board.Reset()
	g.CurrentPlayer = Red
}

func (g *Game) AddTile(col int, player Tile) (Tile, error) {
	if player != g.CurrentPlayer {
		return Empty, fmt.Errorf("not your turn")
	}

	row, err := g.Board.PlaceTile(col, g.CurrentPlayer)

	if err != nil {
		return Empty, err
	}

	// Check for win
	if g.Board.CheckWin(col, row, player) {
		return g.CurrentPlayer, nil
	}

	// Switch player
	if g.CurrentPlayer == Red {
		g.CurrentPlayer = Yellow
	} else if g.CurrentPlayer == Yellow {
		g.CurrentPlayer = Red
	}

	// Send new board
	return Empty, nil
}

func NewGame(id string) Game {
	return Game{
		Id:            id,
		Board:         NewBoard(),
		CurrentPlayer: Red,
	}
}

type PlayerSession struct {
	Player Tile
	Board  Board
	IsTurn bool
}

func new_game_id() (string, error) {
	b := make([]byte, 4)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}

func UniqueGameId(games *sync.Map, attempts uint) (string, error) {

	var err error = nil
	var uuid string

	var i uint
	for i = 0; i < attempts; i++ {
		// Try 10 times to generate a unique game ID
		// If it fails 10 times, return an error
		uuid, err = new_game_id()

		if err != nil {
			continue
		}

		if _, ok := games.Load(uuid); !ok {
			break
		}

		err = fmt.Errorf("Game ID already exists")
	}
	return uuid, err
}

func main() {
	fmt.Println("Initializing server...")

	templates, err := template.ParseFiles("templates/index.html", "templates/game.html")

	if err != nil {
		fmt.Println(err)
		return
	}

	// TODO: Figure out how to clean up inactive games

	// Investigate how go handles map size (auto downsizing?)

	var games sync.Map

	var mux = http.NewServeMux()

	mux.HandleFunc("/new", func(w http.ResponseWriter, r *http.Request) {

		// Try to generate a unique game ID
		uuid, err := UniqueGameId(&games, 10)

		if err != nil {
			http.Error(w, "Could not generate unique game ID after 10 attempts", http.StatusInternalServerError)
			return
		}

		gameState := NewGame(uuid)

		games.Store(uuid, gameState)

		http.Header.Add(w.Header(), "HX-Push-Url", "/game/"+uuid)

		err = templates.ExecuteTemplate(w, "board", gameState)
		if err != nil {
			fmt.Println(err)
		}
	})

	mux.HandleFunc("/game/{id}", func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		v, ok := games.Load(idString)

		if !ok {
			// TODO: Send 404 page
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}

		gameState := v.(Game)

		err = templates.ExecuteTemplate(w, "game", gameState)

	})

	mux.HandleFunc("/game/{id}/add", func(w http.ResponseWriter, r *http.Request) {

		idString := r.PathValue("id")

		v, ok := games.Load(idString)

		if !ok {
			// TODO: Send 404 page
			http.Error(w, "Game not found", http.StatusNotFound)
			return
		}

		gameState := v.(Game)

		params := r.URL.Query()
		column := params.Get("c")
		col, err := strconv.Atoi(column)

		if err != nil {
			http.Error(w, "Column number invalid", http.StatusBadRequest)
			return
		}

		playerString := params.Get("p")
		playerNum, err := strconv.Atoi(playerString)

		if err != nil {
			http.Error(w, "Player number could not be parsed", http.StatusBadRequest)
			return
		}

		if playerNum != int(Red) && playerNum != int(Yellow) {
			http.Error(w, "Invalid player number", http.StatusBadRequest)
		}
		player := Tile(playerNum)

		// Inputs validated by this point
		// What follows is all game logic

		winner, err := gameState.AddTile(col, player)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if winner != Empty {
			if winner == Draw {
				fmt.Println("Draw")
			} else {
				fmt.Printf("Player %d Win", winner)
			}
			// TODO: Send game over message instead of resetting
			gameState.Reset()
		}

		// Save game state
		games.Store(idString, gameState)

		// Send new board
		err = templates.ExecuteTemplate(w, "board", gameState)

		if err != nil {
			fmt.Println(err)
		}

	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err = templates.ExecuteTemplate(w, "index", nil)
		if err != nil {
			fmt.Println(err)
		}
	})

	fmt.Println("Listening on port 8080")
	fmt.Println("Serving at http://localhost:8080/")

	err = http.ListenAndServe(":8080", mux)

	if err != nil {
		fmt.Println(err)
		return
	}

}
