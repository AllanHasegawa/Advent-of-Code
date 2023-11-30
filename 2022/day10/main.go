package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	// "sort"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type instruction_t interface {
	decCycles() bool
	process(state_t) state_t
}

type state_t struct {
	registerX        int
	cycles           int
	instructionsToDo []instruction_t
}

func newState(instructions []instruction_t) state_t {
	return state_t{registerX: 1, cycles: 0, instructionsToDo: instructions}
}

func moveToNextInstruction(s state_t) state_t {
	s.instructionsToDo = s.instructionsToDo[1:]
	return s
}

type noop_t struct {
	cycles int
}

func newNoop() noop_t {
	return noop_t{cycles: 1}
}

type addx_t struct {
	value  int
	cycles int
}

func newAddx(value int) addx_t {
	return addx_t{cycles: 2, value: value}
}

func (i *addx_t) decCycles() bool {
	i.cycles--
	return i.cycles == 0
}

func (i *addx_t) process(s state_t) state_t {
	s.registerX += i.value
	return s
}

func (i *noop_t) decCycles() bool {
	i.cycles--
	return i.cycles == 0
}

func (i *noop_t) process(s state_t) state_t {
	return s
}

func step(s state_t) (state_t,func(state_t)state_t) {
    var endCycleAction func(state_t) state_t
	currentInstruction := s.instructionsToDo[0]
	s.cycles++
	if currentInstruction.decCycles() {
        endCycleAction = func(endS state_t) state_t { 
		    endS = currentInstruction.process(endS)
            return endS
        }
		s = moveToNextInstruction(s)
	} else {
        endCycleAction = func(endS state_t) state_t { return endS }
    }
	return s, endCycleAction
}

func parseLines(lines []string) []instruction_t {
	instructions := []instruction_t{}

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		var newInstruction instruction_t
		if line == "noop" {
			noop := newNoop()
			newInstruction = &noop
		} else {
			split := strings.Split(line, " ")
			if split[0] == "addx" {
				value, err := strconv.Atoi(split[1])
				check(err)
				addx := newAddx(value)
				newInstruction = &addx
			} else {
				panic("Dont know instruction: " + line)
			}
		}
		instructions = append(instructions, newInstruction)
	}
	return instructions
}

func solveP1(input string, lines []string) string {
	instructions := parseLines(lines)
	state := newState(instructions)
    var endCycleAction func(state_t) state_t
	signalStrength := 0

	for {
		if len(state.instructionsToDo) == 0 {
			break
		}

		state, endCycleAction = step(state)

		if state.cycles == 20 ||
			state.cycles == 60 ||
			state.cycles == 100 ||
			state.cycles == 140 ||
			state.cycles == 180 ||
			state.cycles == 220 {
			signalStrength += state.cycles * state.registerX
		}

        state = endCycleAction(state)
	}

	return fmt.Sprintf("%v", signalStrength)
}

func stepCrtDisplay(s state_t, crtDisplay [6][40]string) [6][40]string {
    lineIdx := (s.cycles-1) / 40
    pixelIdx := (s.cycles-1) % 40

    var pixel string
    if s.registerX >= pixelIdx-1 && s.registerX <= pixelIdx+1 {
        pixel = "#"
    } else {
        pixel = "."
    }
    crtDisplay[lineIdx][pixelIdx] = pixel

    return crtDisplay
}

func printCrtDisplay(display [6][40]string) {
    for x := 0; x < len(display); x++ {
        line := display[x]
        for y := 0; y < len(line); y++ {
            fmt.Print(line[y])
        }
        fmt.Println()
    }
}

func solveP2(input string, lines []string) string {
	instructions := parseLines(lines)
	state := newState(instructions)
    crtDisplay := [6][40]string{}
    var endCycleAction func(state_t) state_t

	for {
		if len(state.instructionsToDo) == 0 {
			break
		}

		state, endCycleAction = step(state)
        crtDisplay = stepCrtDisplay(state, crtDisplay)
        state = endCycleAction(state)
	}

    printCrtDisplay(crtDisplay)

	return "Look up :)"
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
