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
			http.Error(w, "Player number invalid", http.StatusBadRequest)
			return
		}

		if playerNum != int(Red) && playerNum != int(Yellow) {
			http.Error(w, "Invalid player", http.StatusBadRequest)
		}
		player := Tile(playerNum)

		if player != gameState.CurrentPlayer {
			http.Error(w, "Not your turn", http.StatusBadRequest)
		}

		row, err := gameState.Board.PlaceTile(col, gameState.CurrentPlayer)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Check for win
		if gameState.Board.CheckWin(col, row, player) {
			fmt.Printf("Player %d Win", gameState.CurrentPlayer)
			gameState.Board.Reset()
			gameState.CurrentPlayer = Red
			err = templates.ExecuteTemplate(w, "board", gameState)
			if err != nil {
				fmt.Println(err)
			}
			return
		}

		// Switch player
		if gameState.CurrentPlayer == Red {
			gameState.CurrentPlayer = Yellow
		} else if gameState.CurrentPlayer == Yellow {
			gameState.CurrentPlayer = Red
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
