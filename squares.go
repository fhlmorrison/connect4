package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
)

type Game struct {
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

type PlayerSession struct {
	Player Tile
	Board  Board
	IsTurn bool
}

func main() {
	fmt.Println("Hello, world.")

	templates, err := template.ParseFiles("templates/index.html", "templates/game.html")

	if err != nil {
		fmt.Println(err)
		return
	}

	gameState := Game{
		Board:         NewBoard(),
		CurrentPlayer: Red,
	}

	var mux = http.NewServeMux()

	mux.HandleFunc("/add", func(w http.ResponseWriter, r *http.Request) {
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

		// Send new board
		err = templates.ExecuteTemplate(w, "board", gameState)

		if err != nil {
			fmt.Println(err)
		}

	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err = templates.ExecuteTemplate(w, "index", gameState)
		if err != nil {
			fmt.Println(err)
		}
	})

	http.ListenAndServe(":8080", mux)
	fmt.Println("Listening on port 8080")
	fmt.Println("http://localhost:8080/")
}
