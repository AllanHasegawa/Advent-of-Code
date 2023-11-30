package main

import (
	"fmt"
	"github.com/juliangruber/go-intersect/v2"
	"os"
	"strings"
	// "strconv"
	// "sort"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func findEqualItem(line string) byte {
	var pivotIdx = len(line) / 2
	var firstCompartmentSlice = []byte(line[0:pivotIdx])
	var secondCompartmentSlice = []byte(line[pivotIdx:])

	var equalItem byte = intersect.HashGeneric(firstCompartmentSlice, secondCompartmentSlice)[0]

	return equalItem
}

func itemToPriority(item byte) int {
    if item > 'a' {
        return int(item - 'a') + 1
    } else {
        return int(item - 'A') + 27
    }
}

func findSecurityBadge(line0 string, line1 string, line2 string) byte {
    var firstRepeatsInTwoRugsacks = intersect.HashGeneric([]byte(line0), []byte(line1))
    var remainingRepeats = intersect.HashGeneric(firstRepeatsInTwoRugsacks, []byte(line2))

    return remainingRepeats[0]
}

func solveP1(input string, lines []string) string {
    totalPriorities := 0

    for i := 0; i < len(lines); i++ {
	    equalItem := findEqualItem(lines[i])
        totalPriorities += itemToPriority(equalItem)
    }
	return fmt.Sprintf("%d", totalPriorities)
}

func solveP2(input string, lines []string) string {
    totalPriorities := 0

    for i := 0; i < len(lines); i += 3 {
	    equalItem := findSecurityBadge(lines[i], lines[i+1], lines[i+2])
        totalPriorities += itemToPriority(equalItem)
    }
	return fmt.Sprintf("%d", totalPriorities)
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
