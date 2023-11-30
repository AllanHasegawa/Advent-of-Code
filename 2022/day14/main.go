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

type point_t struct {
    x, y int
}

type line_t struct {
    start,end point_t
}

type world_t struct {
    things map[point_t](thing_t)
    lowestPointDepth int
    hasFloor bool
}

type thing_t int

const (
    Air thing_t = 0
    Rock thing_t = 1
    Sand thing_t = 2
)

func getThing(w world_t, p point_t) thing_t {
    if w.hasFloor && p.y == getLowestPointDepth(w) {
        return Rock
    }
    return w.things[p]
}

func getLowestPointDepth(w world_t) int {
    if w.hasFloor {
        return w.lowestPointDepth + 2
    }
    return w.lowestPointDepth
}

func parsePoint(rawPoint string) point_t {
    splits := strings.Split(rawPoint, ",")

    x, err := strconv.Atoi(splits[0])
    check(err)
    y, err := strconv.Atoi(splits[1])
    check(err)

    return point_t{x: x, y: y}
}

func parseLine(raw string) []line_t {
    lines := []line_t{}
    splits := strings.Split(raw, " -> ")

    for i := 1; i < len(splits); i++ {
        line := line_t{start: parsePoint(splits[i-1]), end: parsePoint(splits[i])}
        lines = append(lines, line)
    }

    return lines
}

func parseLines(rawLines []string) []line_t {
    lines := []line_t{}

    for _, l := range rawLines {
        lines = append(lines, parseLine(l)...)
    }

    return lines
}

func Abs(i int) int {
    if i < 0 {
        return -i
    }
    return i
}

func lineInterpolation(line line_t) []point_t {
    xDirection := line.end.x - line.start.x
    if xDirection != 0 {
        xDirection = xDirection / Abs(xDirection)
    }
    yDirection := line.end.y - line.start.y
    if yDirection != 0 {
        yDirection = yDirection / Abs(yDirection)
    }
    points := []point_t{}
    current := line.start

    for ;; {
        points = append(points, current)
        if current.x == line.end.x && current.y == line.end.y {
            break
        }

        current.x += xDirection
        current.y += yDirection
    }

    return points
}

func addRocksToWorld(w world_t, lines []line_t) world_t {
    for _, l := range lines {
        for _, p := range lineInterpolation(l) {
            w.things[p] = Rock
            if p.y > w.lowestPointDepth {
                w.lowestPointDepth = p.y
            }
        }
    }
    return w
}

func pourSand(w world_t, source point_t) (world_t, bool) {
    current := source
    moved := true

    for ;; {
        if moved && current.y > getLowestPointDepth(w) {
            return w, false
        }

        if !moved {
            w.things[current] = Sand

            if current.x == source.x && current.y == source.y {
                return w, false
            }
            return w, true
        }
        moved = false

        // try down
        next := current
        next.y++
        dst := getThing(w, next)

        if dst == Air {
            current = next
            moved = true
            continue
        }

        // try down-left
        next = current
        next.x--
        next.y++
        dst = getThing(w, next)

        if dst == Air {
            current = next
            moved = true
            continue
        }

        // try down-right
        next = current
        next.x++
        next.y++
        dst = getThing(w, next)

        if dst == Air {
            current = next
            moved = true
            continue
        }
    }
}

func solveP1(input string, lines []string) string {
    rockLines := parseLines(lines)
    w := world_t{things: make(map[point_t]thing_t)}
    w.hasFloor = false
    w = addRocksToWorld(w, rockLines)

    sandSource := point_t{x: 500, y: 0}
    sandsDropped := 0

    for {
        ok := false
        w, ok = pourSand(w, sandSource)
        if !ok {
            break
        }
        sandsDropped++
    }

    return fmt.Sprintf("%v", sandsDropped)
}

func solveP2(input string, lines []string) string {
    rockLines := parseLines(lines)
    w := world_t{things: make(map[point_t]thing_t)}
    w.hasFloor = true
    w = addRocksToWorld(w, rockLines)

    sandSource := point_t{x: 500, y: 0}
    sandsDropped := 0

    for {
        ok := false
        w, ok = pourSand(w, sandSource)
        sandsDropped++
        if !ok {
            break
        }
    }

    return fmt.Sprintf("%v", sandsDropped)
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
