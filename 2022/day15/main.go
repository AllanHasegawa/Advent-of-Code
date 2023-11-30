package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"golang.org/x/exp/maps"
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
	start point_t
	end   point_t
}

func toInt(s string) int {
	i, err := strconv.Atoi(s)
	check(err)
	return i
}

func parseLine(rawLine string) line_t {
	const prefix = "Sensor at x="
	middle := strings.Split(rawLine, ": closest beacon is at x=")
	sensorSplit := strings.Split(middle[0][len(prefix):], ", y=")
	beaconSplit := strings.Split(middle[1], ", y=")

	return line_t{
		start: point_t{x: toInt(sensorSplit[0]), y: toInt(sensorSplit[1])},
		end:   point_t{x: toInt(beaconSplit[0]), y: toInt(beaconSplit[1])},
	}
}

func parseLines(lines []string) []line_t {
	ls := []line_t{}

	for _, l := range lines {
		ls = append(ls, parseLine(l))
	}

	return ls
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func dist(a point_t, b point_t) int {
	return abs(a.x-b.x) + abs(a.y-b.y)
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a int, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func mergeLines(lines []line_t) []line_t {
	sort.Slice(lines, func(i, j int) bool {
		return lines[i].start.x < lines[j].start.x
	})
	linesP := []line_t{}
	current := lines[0]

	for i := 1; i < len(lines); i++ {
		looking := lines[i]
		if looking.start.x-1 <= current.end.x {
			if looking.end.x > current.end.x {
				current.end.x = looking.end.x
			}
            if i == len(lines)-1 {
                linesP = append(linesP, current)
            }
		} else {
			linesP = append(linesP, current)
			current = looking
		}
	}

	return linesP
}

func countInvalidBeaconPosInRowLimited(lines []line_t, row int, xLimitMin int, xLimitMax int) (int, []line_t) {
	signals := []line_t{}

	for _, l := range lines {
		strength := dist(l.start, l.end)
		if abs(l.start.y-row) <= strength {
			strengthY := strength - abs(l.start.y-row)
			xMin := max(l.start.x-strengthY, xLimitMin)
			xMax := min(l.start.x+strengthY, xLimitMax)

			lineStart := point_t{x: xMin, y: row}
			lineEnd := point_t{x: xMax, y: row}

			signals = append(signals, line_t{start: lineStart, end: lineEnd})
		}
	}

    beaconsX := map[int]bool{}
    for _, l := range lines {
        if l.end.y == row {
            beaconsX[l.end.x] = true
        }
    }

	signals = mergeLines(signals)
	countUnavailablePos := 0
	for _, s := range signals {
		countUnavailablePos += s.end.x - s.start.x
        for _, bX := range maps.Keys(beaconsX) {
            if bX >= s.start.x && bX <= s.end.x {
                countUnavailablePos--
            }
        }
	}

	return countUnavailablePos + 1, signals
}

func solveP1(input string, lines []string) string {
	ls := parseLines(lines)
	// row := 2000000
	row := 10
	invalidPos, _ := countInvalidBeaconPosInRowLimited(ls, row, math.MinInt, math.MaxInt)

	return fmt.Sprintf("%v", invalidPos)
}

func solveP2(input string, lines []string) string {
	ls := parseLines(lines)
    limit := 4000000
    // limit := 20
    sol := -1

    for row := 0; row <= limit; row++ {
	    _, signals := countInvalidBeaconPosInRowLimited(ls, row, 0, limit)
        if len(signals) > 1 {
            sol = (signals[0].end.x+1) * 4000000 + row
        }
    }

	return fmt.Sprintf("%v", sol)
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
