package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
    "github.com/golang-collections/collections/stack"
	// "sort"
)

type instruction struct {
	source int
	target int
	amount int
}

type crate byte

type pile []crate

type input struct {
	piles        []pile
	instructions []instruction
}

type state struct {
    piles []*stack.Stack
    instructionsToDo []instruction
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func parsePilesInput(lines []string) []pile {
    countPiles := (len(lines[0]) + 1)/4
    var piles []pile = make([]pile, countPiles)

    pileIdxToRawIdx := func(pIdx int) int {
        offset := pIdx * 4
        return offset + 1
    }

    getCrate := func(pIdx int, line string) crate {
        hasCrate := line[pileIdxToRawIdx(pIdx) - 1] == '['

        if hasCrate {
            return crate(line[pileIdxToRawIdx(pIdx)])
        } else {
            return crate(0)
        }
    }

    for i := 0; i < len(lines); i++ {
        line := lines[i]
        for p := 0; p < countPiles; p++ {
            crate := getCrate(p, line)
            if crate > 0 {
                piles[p] = append(piles[p], crate)
            }
        }
    }

	return piles
}

func parseInstructions(lines []string) []instruction {
	var instructions []instruction

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		splits := strings.Split(line, " ")
		amount, err := strconv.Atoi(splits[1])
		check(err)
		source, err := strconv.Atoi(splits[3])
		check(err)
		target, err := strconv.Atoi(splits[5])
		check(err)

		instruction := instruction{
			amount: amount,
			source: source - 1,
			target: target - 1,
		}
		instructions = append(instructions, instruction)
	}

	return instructions
}

func parseInput(lines []string) input {
	var breakInputTypeIdx int

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if len(line) == 0 {
			breakInputTypeIdx = i
			break
		}
	}

	return input{
		piles:        parsePilesInput(lines[0 : breakInputTypeIdx-1]),
		instructions: parseInstructions(lines[breakInputTypeIdx+1:]),
	}
}

func initialState(input input) state {
    var stacks []*stack.Stack

    for i := 0; i < len(input.piles); i++ {
        pile := input.piles[i]

        stack := stack.New()
        stacks = append(stacks, stack)

        for p := len(pile)-1; p >= 0; p-- {
            stack.Push(pile[p])
        }
    }

    return state{piles: stacks, instructionsToDo: input.instructions}
}

func doNextInstruction(current state) state {
    instruction := current.instructionsToDo[0]

    for i := 0; i < instruction.amount; i++ {
        crate := current.piles[instruction.source].Pop().(crate)
        current.piles[instruction.target].Push(crate)
    }

    return state{
        piles: current.piles,
        instructionsToDo: current.instructionsToDo[1:],
    }
}

func printState(current state) {
    for i := 0; i < len(current.piles); i++ {
        top := current.piles[i].Peek()
        if top == nil {
            fmt.Println(i, top)
        } else {
            fmt.Println(i, string(top.(crate)))
        }
    }

    fmt.Println("TODO:")
    for i := 0; i < len(current.instructionsToDo); i++ {
        fmt.Println(current.instructionsToDo[i])
    }
}

func solveP1(input string, lines []string) string {
    state := initialState(parseInput(lines))
    result := ""

    for ;; {
        // printState(state)
        if len(state.instructionsToDo) == 0 {
            break
        }
        state = doNextInstruction(state)
    }

    for i := 0; i < len(state.piles); i++ {
        result += string(state.piles[i].Peek().(crate))
    }

	return result
}

func doNextInstructionWith9001(current state) state {
    instruction := current.instructionsToDo[0]
    stack9001 := stack.New()

    for i := 0; i < instruction.amount; i++ {
        crate := current.piles[instruction.source].Pop().(crate)
        stack9001.Push(crate)
    }
    for i := 0; i < instruction.amount; i++ {
        current.piles[instruction.target].Push(stack9001.Pop())
    }

    return state{
        piles: current.piles,
        instructionsToDo: current.instructionsToDo[1:],
    }
}

func solveP2(input string, lines []string) string {
    state := initialState(parseInput(lines))
    result := ""

    for ;; {
        // printState(state)
        if len(state.instructionsToDo) == 0 {
            break
        }
        state = doNextInstructionWith9001(state)
    }

    for i := 0; i < len(state.piles); i++ {
        result += string(state.piles[i].Peek().(crate))
    }

	return result
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
