package core

import (
	"fmt"
	"math/rand"
	"sync"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var board [height][width]int
var nextState [height][width]int
var count int
var cellState int
var decider int
var genCount int
var displayBuffer [height][width]int

var wg sync.WaitGroup

const (
	height int     = 64
	width  int     = 128
	winX   float64 = 512
	winY   float64 = 364
)

func Logic(win *pixelgl.Window) {
	initState()
	randomState()
	genCount = 0
	for {
		clearBoard()
		renderBoard(board)
		displayBuffer = nextBoardState()
		genCount += 1
		nextGen(genCount)
		// wg.Add(1)
		DrawScreen(&displayBuffer, *win)
		wg.Wait()
	}
}

func Gui() {
	cfg := pixelgl.WindowConfig{
		Title:  "LifeG",
		Bounds: pixel.R(0, 0, winX, winY),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	Logic(win)
	// for !win.Closed() {
	// 	win.Update()
	// }
}

func DrawScreen(buffer *[height][width]int, win pixelgl.Window) {
	win.Clear(colornames.Black)
	drew := imdraw.New(nil)
	drew.Color = pixel.RGB(1, 1, 1)
	dispW, dispH := winX/float64(width), winY/float64(height)
	fmt.Println("dispW:", dispW, "dispH:", dispH)

	for y := 0; y < int(height); y++ {
		for x := 0; x < int(width); x++ {
			if buffer[y][x] == 0 {
				// drew.Color = pixel.RGB(0, 0, 0)
				continue
			}
			drew.Push(pixel.V(dispW*float64(y), dispH*float64(x)))
			drew.Push(pixel.V(dispW*float64(y)+dispW, dispH*float64(x)+dispH))
			drew.Rectangle(0)
		}
	}
	drew.Draw(&win)
	win.Update()
	// wg.Done()
}

// Move the cursor to the upper-left corner of the screen.
// Delete everything from the cursor to the end of the line.
func clearBoard() {
	for i := 1; i < height+1; i++ {
		fmt.Print("\033[H\033[K")
	}
}

func initState() {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			board[i][j] = 0
		}
	}
}
func randomState() {
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			randomInt := rand.Intn(2)
			board[i][j] = randomInt
		}
	}
}

func renderBoard(b [height][width]int) {
	// fmt.Print("\033[?25l") //Hide the cursor.
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if b[i][j] == 0 {
				fmt.Printf("%d \033[31m", b[i][j])
			} else if b[i][j] == 1 {
				fmt.Printf("%d \033[36m", b[i][j])
			}
		}
		fmt.Println()
	}

}

// Any live cell with 0 or 1 live neighbors becomes dead, because of underpopulation
// Any live cell with 2 or 3 live neighbors stays alive, because its neighborhood is just right
// Any live cell with more than 3 live neighbors becomes dead, because of overpopulation
// Any dead cell with exactly 3 live neighbors becomes alive, by reproduction
func nextBoardState() [height][width]int {
	// initialState := boardState
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			cellState = board[i][j]
			// fmt.Println("i:", i, "j:", j, "cell state is:", cellState)
			if (i-1 > -1 && i+1 < height) && (j-1 > -1 && j+1 < width) {
				decider = normalNeighbours(i, j, board)
			} else if (i == 0 || i == (height-1)) && (j > 0 && j < width-1) {
				decider = edgeRows(i, j, board)
			} else if (j == 0 || j == (width-1)) && (i > 0 && i < height-1) {
				decider = edgeColumns(i, j, board)
			} else if (i == 0 && j == 0) || (i == 0 && j == width-1) || (i == height-1 && j == 0) || (i == height-1 && j == width-1) {
				decider = cornerGrids(i, j, board)
			}

			if cellState == 1 && decider > 3 {
				cellState = 0
			} else if cellState == 1 && decider < 2 {
				cellState = 0
			} else if cellState == 1 && (decider == 2 || decider == 3) {
				cellState = 1
			} else if cellState == 0 && decider == 3 {
				cellState = 1
			}
			board[i][j] = cellState
		}
	}
	return board
}
func normalNeighbours(i, j int, iState [height][width]int) int {
	count = 0
	if iState[i][j-1] == 1 {
		count += 1
	}
	if iState[i][j+1] == 1 {
		count += 1
	}
	if iState[i-1][j-1] == 1 {
		count += 1
	}
	if iState[i-1][j] == 1 {
		count += 1
	}
	if iState[i-1][j+1] == 1 {
		count += 1
	}
	if iState[i+1][j] == 1 {
		count += 1
	}
	if iState[i+1][j-1] == 1 {
		count += 1
	}
	if iState[i+1][j+1] == 1 {
		count += 1
	}
	return count
}

func edgeColumns(i, j int, iState [height][width]int) int {
	count = 0
	switch j {
	case 0:
		if iState[i-1][j] == 1 {
			count += 1
		}
		if iState[i+1][j] == 1 {
			count += 1
		}
		if iState[i][j+1] == 1 {
			count += 1
		}
		if iState[i-1][j+1] == 1 {
			count += 1
		}
		if iState[i+1][j+1] == 1 {
			count += 1
		}
	case width - 1:
		if iState[i-1][j] == 1 {
			count += 1
		}
		if iState[i+1][j] == 1 {
			count += 1
		}
		if iState[i][j-1] == 1 {
			count += 1
		}
		if iState[i-1][j-1] == 1 {
			count += 1
		}
		if iState[i+1][j-1] == 1 {
			count += 1
		}
	}
	return count
}
func edgeRows(i, j int, iState [height][width]int) int {
	count = 0
	switch i {
	case 0:
		if iState[i+1][j] == 1 {
			count += 1
		}
		if iState[i+1][j+1] == 1 {
			count += 1
		}
		if iState[i+1][j-1] == 1 {
			count += 1
		}
		if iState[i][j+1] == 1 {
			count += 1
		}
		if iState[i][j-1] == 1 {
			count += 1
		}
	case height - 1:
		if iState[i-1][j] == 1 {
			count += 1
		}
		if iState[i-1][j+1] == 1 {
			count += 1
		}
		if iState[i-1][j-1] == 1 {
			count += 1
		}
		if iState[i][j-1] == 1 {
			count += 1
		}
		if iState[i][j+1] == 1 {
			count += 1
		}
	}
	return count
}

func cornerGrids(i, j int, iState [height][width]int) int {
	count = 0
	switch i {
	case 0:
		if j == 0 {
			if iState[i][j+1] == 1 {
				count += 1
			}
			if iState[i+1][j] == 1 {
				count += 1
			}
			if iState[i+1][j+1] == 1 {
				count += 1
			}
		} else if j == (width - 1) {
			if iState[i][j-1] == 1 {
				count += 1
			}
			if iState[i+1][j-1] == 1 {
				count += 1
			}
			if iState[i+1][j] == 1 {
				count += 1
			}
		}
	case height - 1:
		if j == 0 {
			if iState[i-1][j] == 1 {
				count += 1
			}
			if iState[i-1][j+1] == 1 {
				count += 1
			}
			if iState[i][j+1] == 1 {
				count += 1
			}
		} else if j == width-1 {
			if iState[i-1][j] == 1 {
				count += 1
			}
			if iState[i-1][j-1] == 1 {
				count += 1
			}
			if iState[i][j-1] == 1 {
				count += 1
			}
		}
	}
	return count
}

func nextGen(int) {
	fmt.Println("-----------------------------------------------")
	fmt.Println("Life's generation:", genCount)
	fmt.Println("-----------------------------------------------")
}
