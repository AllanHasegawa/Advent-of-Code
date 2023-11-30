package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func solveP1(input string, lines []string) string {
	var currentElfCalories int = 0
	var highestElfCalories int = 0

	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if len(line) == 0 {
			if currentElfCalories > highestElfCalories {
				highestElfCalories = currentElfCalories
			}
			currentElfCalories = 0
		} else {
			calories, err := strconv.Atoi(line)
			check(err)
			currentElfCalories += calories
		}
	}

	return fmt.Sprintf("%d", highestElfCalories)
}

func solveP2(input string, lines []string) string {
	var elvesCalories []int = []int{0}
	elfIndex := 0

	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if len(line) == 0 {
			elfIndex++
			elvesCalories = append(elvesCalories, 0)
		} else {
			calories, err := strconv.Atoi(line)
			check(err)
			elvesCalories[elfIndex] += calories
		}
	}

	sort.Sort(sort.Reverse(sort.IntSlice(elvesCalories)))

	top3Calories := 0

	for i := 0; i < 3; i++ {
		top3Calories += elvesCalories[i]
	}

	return fmt.Sprintf("%d", top3Calories)
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
