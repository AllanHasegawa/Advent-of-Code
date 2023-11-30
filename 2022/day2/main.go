package main

import (
	"fmt"
	"os"
	"strings"
	// "strconv"
	// "sort"
)

type SymbolP1 byte

const (
	ROCK_IN     SymbolP1 = 'A'
	PAPER_IN             = 'B'
	SCISSOR_IN           = 'C'
	ROCK_OUT             = 'X'
	PAPER_OUT            = 'Y'
	SCISSOR_OUT          = 'Z'
)

type SymbolP2 byte

const (
	ROCK    SymbolP2 = 'A'
	PAPER            = 'B'
	SCISSOR          = 'C'
	DO_LOSE          = 'X'
	DO_DRAW          = 'Y'
	DO_WIN           = 'Z'
)

type RoundResult uint8

const (
	LOST RoundResult = 0
	DRAW             = 1
	WON              = 2
)

type round_input_p1 struct {
	in  SymbolP1
	out SymbolP1
}

type round_input_p2 struct {
	in  SymbolP2
	out SymbolP2
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func roundResultP1(input round_input_p1) RoundResult {
	var in = input.in
	var out = input.out
	switch in {
	case ROCK_IN:
		switch out {
		case ROCK_OUT:
			return DRAW
		case PAPER_OUT:
			return WON
		case SCISSOR_OUT:
			return LOST
		}
	case PAPER_IN:
		switch out {
		case ROCK_OUT:
			return LOST
		case PAPER_OUT:
			return DRAW
		case SCISSOR_OUT:
			return WON
		}
	case SCISSOR_IN:
		switch out {
		case ROCK_OUT:
			return WON
		case PAPER_OUT:
			return LOST
		case SCISSOR_OUT:
			return DRAW
		}
	}
	return LOST
}

func roundScoreP1(input round_input_p1) int {
	var symbolSelectedScore int
	var resultScore int

	switch input.out {
	case ROCK_OUT:
		symbolSelectedScore = 1
	case PAPER_OUT:
		symbolSelectedScore = 2
	case SCISSOR_OUT:
		symbolSelectedScore = 3
	default:
		symbolSelectedScore = 99999999999999
	}

	switch roundResultP1(input) {
	case LOST:
		resultScore = 0
	case DRAW:
		resultScore = 3
	case WON:
		resultScore = 6
	default:
		resultScore = 99999999999999
	}

	return symbolSelectedScore + resultScore
}

func parseLines[R round_input_p1 | round_input_p2](lines []string, mapper func(byte, byte) R) []R {
	var roundInputs []R
	for i := 0; i < len(lines); i++ {
		var line = lines[i]
		roundInputs = append(roundInputs, mapper(line[0], line[2]))
	}
	return roundInputs
}

func solveP1(input string, lines []string) string {
	var totalScore int = 0
	var parsedLines []round_input_p1 = parseLines(
        lines, 
        func(sin byte, sout byte) round_input_p1 { return round_input_p1{in: SymbolP1(sin), out: SymbolP1(sout)} },
    )

	for i := 0; i < len(parsedLines); i++ {
		totalScore += roundScoreP1(parsedLines[i])
	}

	return fmt.Sprintf("%d", totalScore)
}

func selectSymbolP2(input round_input_p2) SymbolP2 {
    switch (input.out) {
    case DO_LOSE: switch input.in {
    case ROCK: return SCISSOR
    case PAPER: return ROCK
    case SCISSOR: return PAPER
    }
    case DO_DRAW: return input.in
    case DO_WIN: switch input.in {
    case ROCK: return PAPER
    case PAPER: return SCISSOR
    case SCISSOR: return ROCK
    }
    }
    return DO_DRAW
}

func roundScoreP2(input round_input_p2) int {
	var symbolSelectedScore int
	var resultScore int

	switch selectSymbolP2(input) {
	case ROCK:
		symbolSelectedScore = 1
	case PAPER:
		symbolSelectedScore = 2
	case SCISSOR:
		symbolSelectedScore = 3
	default:
		symbolSelectedScore = 99999999999999
	}

    switch input.out {
    case DO_LOSE: resultScore = 0
    case DO_DRAW: resultScore = 3
    case DO_WIN: resultScore = 6
    default: resultScore = 99999999999999
    }

    return symbolSelectedScore + resultScore
}

func solveP2(input string, lines []string) string {
	var totalScore int = 0
	var parsedLines []round_input_p2 = parseLines(
        lines, 
        func(sin byte, sout byte) round_input_p2 { return round_input_p2{in: SymbolP2(sin), out: SymbolP2(sout)} },
    )

	for i := 0; i < len(parsedLines); i++ {
		totalScore += roundScoreP2(parsedLines[i])
	}

	return fmt.Sprintf("%d", totalScore)
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
