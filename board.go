package main

import "fmt"

// Tile represents a square on the board.
type Tile int

const (
	Empty Tile = iota
	Red
	Yellow
	Draw
)

// Board represents a 7x6 grid of tiles.
type Board [7][6]Tile

// PlaceTile places a tile in the given column and returns the row it was placed in.
// If the column is full or invalid, an error is returned.
func (b *Board) PlaceTile(column int, tile Tile) (int, error) {
	// TODO

	// Check if column is valid
	if column < 0 || column >= len(b) {
		return -1, fmt.Errorf("invalid column: %d", column)
	}

	// Insert in the first empty row
	for i := len(b[column]) - 1; i >= 0; i-- {
		if b[column][i] == Empty {
			b[column][i] = tile
			return i, nil
		}
	}
	return -1, fmt.Errorf("column %d is full", column)
}

func (b *Board) checkHorizontalWin(row int, tile Tile) bool {
	// Check for horizontal win
	count := 0
	for _, v := range b {
		if v[row] == tile {
			count++
			if count == 4 {
				return true
			}
		} else {
			count = 0
		}
	}
	return false
}

func (b *Board) checkVerticalWin(column int, tile Tile) bool {
	// Check for vertical win
	count := 0
	for _, v := range b[column] {
		if v == tile {
			count++
			if count == 4 {
				return true
			}
		} else {
			count = 0
		}
	}
	return false
}

func (b *Board) checkDiagonalWin(column int, row int, tile Tile) bool {

	count := 0
	// Top left to bottom right
	start := column - row
	for i := 0; i < len(b[0]) || i < len(b) || start+i > len(b); i++ {
		if start+i < 0 {
			continue
		}

		if i >= len(b[0]) || start+i >= len(b) {
			break
		}

		if b[start+i][i] == tile {
			count++
			if count == 4 {
				return true
			}
		} else {
			count = 0
		}
	}

	count = 0
	// Bottom left to top right
	start = column + row
	for i := 0; i < len(b); i++ {
		if start-i >= len(b) {
			continue
		}

		if i >= len(b[0]) || start-i < 0 {
			break
		}

		if b[start-i][i] == tile {
			count++
			if count == 4 {
				return true
			}
		} else {
			count = 0
		}
	}
	return false
}

// Check for win or draw given last placed tile
func (b *Board) CheckWin(column int, row int, tile Tile) bool {
	win := b.checkHorizontalWin(row, tile) || b.checkVerticalWin(column, tile) || b.checkDiagonalWin(column, row, tile)
	return win
}

func (b *Board) CheckDraw() bool {
	// Check for draw
	for _, v := range b {
		for _, t := range v {
			if t == Empty {
				return false
			}
		}
	}
	return true
}

func (b *Board) Reset() {
	*b = [7][6]Tile{}
}

func NewBoard() Board {
	var b Board
	b.Reset()
	return b
}
