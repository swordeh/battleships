package game

import (
	"bytes"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type Board struct {
	Tiles []int
}

func NewEmptyBoard() Board {
	board := &Board{Tiles: make([]int, 100)}
	for i := 0; i < 100; i++ {
		board.Tiles[i] = 0
	}
	return *board
}

func (b *Board) Draw() []byte {
	// Store all output in a bytes.Buffer
	out := bytes.Buffer{}
	// Map int values to the ascii representation of the cell
	mapping := map[int]string{0: "-", 1: "S", 2: "*", 3: "x"} // 0 empty 1 Ship 2 Hit 3 Miss
	header := []string{" ", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "\n"}
	cols := map[int]string{0: "A", 10: "B", 20: "C", 30: "D", 40: "E", 50: "F", 60: "G", 70: "H", 80: "I", 90: "J"}

	// Write the header first
	for i, _ := range header {
		out.Write([]byte(header[i]))
	}
	for i, _ := range b.Tiles {
		if colValue, ok := cols[i]; ok {
			out.Write([]byte(colValue))
		}
		out.Write([]byte(mapping[b.Tiles[i]]))
		if (i+1)%10 == 0 {
			out.Write([]byte("\n"))
		}
	}
	return out.Bytes()
}

func gridToSliceIndex(ref string) (int, error) {
	// check reference is two-part
	if len(ref) < 2 {
		return -1, errors.New("ref too short")
	}

	// check reference is not longer than 3 characters (regex [a-jA-J]10[1-9])
	if len(ref) > 3 {
		return -1, errors.New("ref too long")
	}

	// Calculate the difference between the letter provided and the ascii code for 'a'
	row := int(strings.ToLower(ref)[:1][0] - 'a')
	if row < 0 || row > 9 {
		return -1, errors.New("grid reference out of bounds")
	}

	// Convert the string column to an integer
	col, err := strconv.Atoi(ref[1:])
	if err != nil {
		return -1, errors.New("invalid column number")
	}

	if col < 1 || col > 10 {
		return -1, errors.New("column number out of bounds")
	}

	return row*10 + col - 1, nil
}

func (b *Board) Update(ref string, value int) error {
	i, err := gridToSliceIndex(ref)
	if err != nil {
		errString := fmt.Sprintf("invalid reference %s, %s", ref, err)
		return errors.New(errString)
	}
	b.Tiles[i] = value
	return nil
}

// Fire will check the grid reference and check for a ship. If a ship found, it will be marked as a hit
// it returns true for a hit and returns false for a miss
// Fire will return an error if you attempt to hit a cell already hit
func (b *Board) Fire(ref string) (bool, error) {
	i, err := gridToSliceIndex(ref)
	if err != nil {
		errString := fmt.Sprintf("invalid reference %s, %s", ref, err)
		return false, errors.New(errString)
	}
	if b.Tiles[i] == 2 {
		return false, errors.New("already hit")
	}
	if b.Tiles[i] == 1 {
		// Still count even if there was a hit already
		b.Tiles[i] = 2
		return true, nil
	} else {
		b.Tiles[i] = 3
		return false, nil
	}
}

// Ship can either be horizontal or vertical, and can stack up against, top and tail and
// otherwise touch but may never intersect.

// - - - S S S - - - -
// - - - S - - - - - -
// - - - S - - - - S -
// - - - S - - - - S -
// - - - - - - - - S -
// - S S S S S - - S -
// - - - - - - - - - -
// - - - - - - S S - -
// - - - - - - - - - -

// PlaceShip will position one fleet of 5 ships on the board
// Ships that can be placed go from 2 cells in length to 5, with 2 being 3 cells long
func (b *Board) PlaceShipsRandomly() {
	ships := []int{2, 3, 3, 4, 5}
	for _, ship := range ships {
		placed := false
		for !placed {
			// Generate a random starting position
			start := rand.Intn(100)
			// Generate a random direction (0 = horizontal, 1 = vertical)
			direction := rand.Intn(2)
			// Check if the ship can be placed in the chosen direction
			if canPlaceShip(b, start, ship, direction) {
				placeShip(b, start, ship, direction)
				placed = true
			}
		}
	}
}

// canPlaceShip checks that the placement of a ship doesn't extend horizontally or vertically
// out of bounds of the 10x10 grid. If it does, it returns false.
func canPlaceShip(board *Board, start, ship, direction int) bool {
	// horizontal placement
	if direction == 0 {
		// iterate through every ship tile
		for i := 0; i < ship; i++ {
			// if the board tile the ship tile is to be placed on extends to a new line (goes above 10)
			// or is already occupied, return false
			if start%10+i >= 10 || board.Tiles[start+i] != 0 {
				return false
			}
		}
	} else {
		// vertical placement
		for i := 0; i < ship; i++ {
			// check that the ship tile doesn't extend into a line that doesn't exist below
			// or is already occupied
			if start/10+i >= 10 || board.Tiles[start+i*10] != 0 {
				return false
			}
		}
	}
	return true
}

// placeShip assigns the value 1 to the tiles on the board to indicate the placement of a ship
func placeShip(board *Board, start, ship, direction int) {
	if direction == 0 {
		for i := 0; i < ship; i++ {
			board.Tiles[start+i] = 1
		}
	} else {
		for i := 0; i < ship; i++ {
			board.Tiles[start+i*10] = 1
		}
	}
}
