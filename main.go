package main

import (
	"fmt"
	"math/rand"
)

const (
	height = 32
	width  = 64
)

var boardState [height][width]int
var nextState [height][width]int

func main() {
	deadState()
	randomState()
	renderBoard()
}

func deadState() {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			boardState[i][j] = 0
		}
	}
}
func randomState() {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			randomInt := rand.Intn(2)
			boardState[i][j] = randomInt
		}
	}
}

func renderBoard() {
	//assign some colors to the dead and alive cells
	//print out the board to terminal for now row after row
	//next scope is to use ebiten and display it via GUI
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			fmt.Printf("%d ", boardState[i][j])
		}
		fmt.Println()
	}
}

func nextBoardState(){
	initialState := boardState
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			cell := initialState[i][j]
			if cell == 0 && i == 0{
				//edge cases along the top row, then it has only 5 neighbours as opposed to normal 8 neighbours
				if j == 0

			} else if cell == 1{

			}
		}
	}
}
