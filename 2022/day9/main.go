package main

import (
    "fmt"
    "os"
    "strings"
    "strconv"
    // "sort"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type pos_t struct {
    x int
    y int
}

func add(a pos_t, b pos_t) pos_t {
    return pos_t{x: a.x + b.x, y: a.y + b.y}
}

type instruction_t struct {
    direction pos_t
    count int
}

type state_t struct {
    head pos_t
    tail pos_t
    visited map[pos_t]bool
}

func parseLines(lines []string) []instruction_t {
    var instructions []instruction_t

    for i := 0; i < len(lines); i++ {
        line := lines[i]
        splits := strings.Split(line, " ")
        
        var direction pos_t
        if splits[0] == "R" {
            direction = pos_t{x: 0, y: 1}
        } else if splits[0] == "L" {
            direction = pos_t{x: 0, y: -1}
        } else if splits[0] == "U" {
            direction = pos_t{x: 1, y: 0}
        } else if splits[0] == "D" {
            direction = pos_t{x: -1, y: 0}
        } else {
            panic("Wat!? Don't know the direction: " + splits[0])
        }

        count, err := strconv.Atoi(splits[1])
        check(err)

        instructions = append(instructions, instruction_t{direction: direction, count: count})
    }

    return instructions
}

// This includes diagonally
func areTheyTouching(a pos_t, b pos_t) bool {
    for i := -1; i < 2; i++ {
        for j := -1; j < 2; j++ {
            if a.x+i == b.x && a.y+j == b.y {
                return true
            }
        }
    }
    return false
}

func moveDirection(a int, b int) int {
    diff := b - a
    if diff < 0 {
        return -1
    }
    if diff > 0 {
        return 1
    }
    return 0
}

func newTailPos(tail pos_t, head pos_t) pos_t {
    if areTheyTouching(tail, head) {
        return tail
    }

    horizontalMove := moveDirection(tail.x, head.x)
    verticalMove := moveDirection(tail.y, head.y)
    move := pos_t{x: horizontalMove, y: verticalMove}

    return add(tail, move)
}

func step(instruction instruction_t, s state_t) state_t {
    headPrime := s.head
    tailPrime := s.tail
    for i := 0; i < instruction.count; i++ {
        headPrime = add(headPrime, instruction.direction)
        tailPrime = newTailPos(tailPrime, headPrime)
        s.visited[tailPrime] = true
    }
    s.head = headPrime
    s.tail = tailPrime
    return s
}

func solveP1(input string, lines []string) string {
    instructions := parseLines(lines)

    initial_head := pos_t{x: 0, y: 0}
    initial_tail := pos_t{x: 0, y: 0}
    visited := make(map[pos_t]bool)
    visited[initial_tail] = true
    state := state_t{head: initial_head, tail: initial_tail, visited: visited}

    for i := 0; i < len(instructions); i++ {
        state = step(instructions[i], state)
    }

    return fmt.Sprintf("%v", len(state.visited))
}

type state_p2_t struct {
    knots [10]pos_t
    visited map[pos_t]bool
}

func stepP2(instruction instruction_t, s state_p2_t) state_p2_t {
    for i := 0; i < instruction.count; i++ {
        // Move the head as per the instructions
        s.knots[0] = add(s.knots[0], instruction.direction)

        // Following knots will move following the previous one
        for k := 1; k < 10; k++ {
            leadPos := s.knots[k-1]
            s.knots[k] = newTailPos(s.knots[k], leadPos)
        }

        // We mark the position the last tail visited
        s.visited[s.knots[9]] = true
    }
    return s
}

func solveP2(input string, lines []string) string {
    instructions := parseLines(lines)

    visited := make(map[pos_t]bool)
    visited[pos_t{}] = true
    state := state_p2_t{knots: [10](pos_t){}, visited: visited}

    for i := 0; i < len(instructions); i++ {
        state = stepP2(instructions[i], state)
    }

    return fmt.Sprintf("%v", len(state.visited))
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
