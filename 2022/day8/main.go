package main

import (
	"fmt"
	"os"
	"strings"
	// "strconv"
	// "sort"
)

type board_t struct {
	values  [][]uint8
}

func getValue(board board_t, x int, y int) uint8 {
	return board.values[x][y]
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parseLines(lines []string) board_t {
	var board [][]uint8
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		board = append(board, [](uint8){})

		for j := 0; j < len(line); j++ {
			board[i] = append(board[i], line[j]-'0')
		}
	}
	return board_t{values: board}
}

func isTreeVisible(board board_t, x int, y int) bool {
	height := getValue(board, x, y)

	i := -100
	// up
	for i = x - 1; i >= 0; i-- {
		if getValue(board, i, y) >= height {
			break
		}
	}
	if i == -1 {
		return true
	}
	// down
	for i = x + 1; i < len(board.values); i++ {
		if getValue(board, i, y) >= height {
			break
		}
	}
	if i == len(board.values) {
		return true
	}
	// left
	for i = y - 1; i >= 0; i-- {
		if getValue(board, x, i) >= height {
			break
		}
	}
	if i == -1 {
		return true
	}
	// right
	for i = y + 1; i < len(board.values); i++ {
		if getValue(board, x, i) >= height {
			break
		}
	}
	if i == len(board.values) {
		return true
	}
	return false
}

func solveP1(input string, lines []string) string {
	board := parseLines(lines)
	innerTreesVisible := 0

	boardSizeX := len(board.values)
	boardSizeY := len(board.values[0])

	for x := 1; x < boardSizeX-1; x++ {
		for y := 1; y < boardSizeY-1; y++ {
			if isTreeVisible(board, x, y) {
				innerTreesVisible++
			}
		}
	}

	outerTrees := boardSizeX*2 + boardSizeY*2 - 4

	return fmt.Sprintf("%v", innerTreesVisible+outerTrees)
}

func scenicScore(board board_t, x int, y int) int {
	height := getValue(board, x, y)

	// up
	upScore := 0
	for i := x - 1; i >= 0; i-- {
		value := getValue(board, i, y)
		if height > value {
			upScore++
		} else if height <= value {
			upScore++
			break
		}
	}

	// left
	leftScore := 0
	for i := y - 1; i >= 0; i-- {
		value := getValue(board, x, i)
		if height > value {
			leftScore++
		} else if height <= value {
			leftScore++
			break
		}
	}

	// right
	rightScore := 0
	for i := x + 1; i < len(board.values); i++ {
		value := getValue(board, i, y)
		if height > value {
			rightScore++
		} else if height <= value {
			rightScore++
			break
		}
	}

	// down
	downScore := 0
	for i := y + 1; i < len(board.values); i++ {
		value := getValue(board, x, i)
		if height > value {
			downScore++
		} else if height <= value {
			downScore++
			break
		}
	}

	return upScore * leftScore * downScore * rightScore
}

func solveP2(input string, lines []string) string {
    board := parseLines(lines)
    topScenicScore := 0

    for x := 1; x < len(board.values)-1; x++ {
        for y := 1; y < len(board.values)-1; y++ {
            currentScenicScore := scenicScore(board, x, y)
            if topScenicScore < currentScenicScore {
                topScenicScore = currentScenicScore
            }
        }
    }
	return fmt.Sprintf("%v", topScenicScore)
}

func main() {
	dat, err := os.ReadFile("input.txt")
	check(err)

	input := string(dat)
	lines := strings.Split(input, "\n")

	fmt.Println("Solution P1:")
	fmt.Println(solveP1(input, lines))
	fmt.Println("")
	fmt.Println("Solution P2:")
	fmt.Println(solveP2(input, lines))
}
