package main

import (
	"fmt"
	"os"
	"strings"
	"github.com/emirpasic/gods/sets/hashset"
	// "strconv"
	// "sort"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func windowIsStarterPack(input string, idx int, size int) bool {
    // idx = 2, size = 4, true
    // idx = 3, size = 4, false
    // idx = 4, size = 4, false
    if (idx + 1) < size {
         return false
    }

    set := hashset.New()

    for i := 0; i < size; i++ {
        set.Add(input[idx - i])

        if set.Size() != (i + 1) {
            break
        }
    }

    return set.Size() == size
}

func findStartPackIdx(input string, packSize int) int {
    for i := packSize; i < len(input); i++ {
        if windowIsStarterPack(input, i, packSize) {
            return i
        }
    }
    return -1
}

func solveP1(input string, lines []string) string {
    return fmt.Sprintf("%d", findStartPackIdx(input, 4) + 1)
}

func solveP2(input string, lines []string) string {
    return fmt.Sprintf("%d", findStartPackIdx(input, 14) + 1)
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
