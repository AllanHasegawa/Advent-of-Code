package main

import (
	"fmt"
	"os"
	"strings"
	// "strconv"
	// "sort"
)

const ARCHIVE_SIZE = 100
const TARGET_P1 int64 = 2022
const TARGET_P2 int64 = 1000000000000

const (
	Left  byte = '<'
	Right byte = '>'
)

type point_t struct {
	x, y int64
}

type shape_t struct {
	points []point_t
	height int64
	width  int64
}

type falling_shape_t struct {
	shape  shape_t
	offset point_t
}

type state_t struct {
	rocks          map[point_t]bool
	numberOfShapes int64
	fallingShape   *falling_shape_t
	width          int64
	peak           int64
	input          []byte
	inputIdx       int
	archives       map[archive_key_t]([]archive_t)
	skipped        bool
	target         int64
}

type archive_key_t struct {
	fallingShapeIdx int
	inputIdx        int
}

type archive_t struct {
	restingRocks   [ARCHIVE_SIZE][7]bool
	key            archive_key_t
	peak           int64
	numberOfShapes int64
}

func archive(state state_t) *archive_t {
	key := archive_key_t{
		fallingShapeIdx: int(state.numberOfShapes % 5),
		inputIdx:        state.inputIdx,
	}
	restingRocks := [ARCHIVE_SIZE][7]bool{}
	for x := 0; x < ARCHIVE_SIZE; x++ {
		for y := 0; y < 7; y++ {
			p := point_t{x: state.peak + int64(x), y: int64(y)}
			restingRocks[x][y] = hasRock(p, state)
		}
	}
	archive := archive_t{
		key:            key,
		restingRocks:   restingRocks,
		peak:           state.peak,
		numberOfShapes: state.numberOfShapes,
	}
	arr, has := state.archives[key]
	if has {
		for _, old := range state.archives[key] {
			if areArchiveTheSame(old, archive) {
				return &old
			}
		}
		state.archives[key] = append(arr, archive)
	} else {
		arr = []archive_t{archive}
		state.archives[key] = arr
	}
	return nil
}

func areArchiveTheSame(a archive_t, b archive_t) bool {
	if a.key.fallingShapeIdx != b.key.fallingShapeIdx || a.key.inputIdx != b.key.inputIdx {
		return false
	}
	for x := 0; x < ARCHIVE_SIZE; x++ {
		for y := 0; y < 7; y++ {
			if a.restingRocks[x][y] != b.restingRocks[x][y] {
				return false
			}
		}
	}

	return true
}

func shapeLine() shape_t {
	return shape_t{
		points: []point_t{
			{x: 0, y: 0},
			{x: 0, y: 1},
			{x: 0, y: 2},
			{x: 0, y: 3},
		},
		height: 1,
		width:  4,
	}
}

func shapeCross() shape_t {
	return shape_t{
		points: []point_t{
			{x: 0, y: 1},
			{x: 1, y: 0},
			{x: 1, y: 1},
			{x: 1, y: 2},
			{x: 2, y: 1},
		},
		height: 3,
		width:  3,
	}
}

func shapeL() shape_t {
	return shape_t{
		points: []point_t{
			{x: 0, y: 2},
			{x: 1, y: 2},
			{x: 2, y: 2},
			{x: 2, y: 1},
			{x: 2, y: 0},
		},
		height: 3,
		width:  3,
	}
}

func shapeBar() shape_t {
	return shape_t{
		points: []point_t{
			{x: 0, y: 0},
			{x: 1, y: 0},
			{x: 2, y: 0},
			{x: 3, y: 0},
		},
		height: 4,
		width:  1,
	}
}

func shapeSquare() shape_t {
	return shape_t{
		points: []point_t{
			{x: 0, y: 0},
			{x: 0, y: 1},
			{x: 1, y: 0},
			{x: 1, y: 1},
		},
		height: 2,
		width:  2,
	}
}

func getShapeInOrder(idx int64) shape_t {
	idx = idx % 5

	var shape shape_t

	switch idx {
	case 0:
		shape = shapeLine()
	case 1:
		shape = shapeCross()
	case 2:
		shape = shapeL()
	case 3:
		shape = shapeBar()
	case 4:
		shape = shapeSquare()
	default:
		panic("Shape unknown order")
	}

	return shape
}

func add(a point_t, b point_t) point_t {
	return point_t{x: a.x + b.x, y: a.y + b.y}
}

func getStartingOffset(state state_t, shape shape_t) point_t {
	return point_t{x: state.peak - 3 - shape.height, y: 2}
}

func directionToOffset(direction byte) point_t {
	if direction == '<' {
		return point_t{x: 0, y: -1}
	} else {
		return point_t{x: 0, y: +1}
	}
}

func hasRock(p point_t, state state_t) bool {
	return state.rocks[p]
}

func canShapeFitNextOffset(state state_t, offset point_t) bool {
	newOffset := add(state.fallingShape.offset, offset)
	for _, p := range state.fallingShape.shape.points {
		movedP := add(p, newOffset)

		// out of bounds
		if movedP.x > 0 || movedP.y < 0 || movedP.y >= state.width {
			return false
		}

		// already has rocks
		if hasRock(movedP, state) {
			return false
		}
	}
	return true
}

func addRock(p point_t, state state_t) {
	state.rocks[p] = true
}

func addFallingShapeRocksToMap(state state_t) state_t {
	offset := state.fallingShape.offset
	for _, p := range state.fallingShape.shape.points {
		movedP := add(p, offset)
		addRock(movedP, state)
	}
	return state
}

func findNewHeight(state state_t) int64 {
	offset := state.fallingShape.offset
	highPoint := state.fallingShape.shape.points[0] // first point should be the highest
	offsetedPoint := add(offset, highPoint)
	if state.peak < offsetedPoint.x {
		return state.peak
	}
	return offsetedPoint.x
}

func copyRocksFromArchive(a archive_t, state state_t) {
	for x := 0; x < ARCHIVE_SIZE; x++ {
		for y := 0; y < 7; y++ {
			p := point_t{x: state.peak + int64(x), y: int64(y)}
			if a.restingRocks[x][y] {
				addRock(p, state)
			}
		}
	}
}

func step(state state_t, allowSkip bool) state_t {
	if state.fallingShape == nil {
		if allowSkip && !state.skipped {
			archive := archive(state)
			if archive != nil {
				targetRocks := state.target - state.numberOfShapes
				diffPeak := state.peak - archive.peak
				diffShapes := state.numberOfShapes - archive.numberOfShapes
				cycles := targetRocks / int64(diffShapes)
				state.peak = state.peak + diffPeak*cycles
				state.numberOfShapes = state.numberOfShapes + diffShapes*cycles
				copyRocksFromArchive(*archive, state)

				state.archives = map[archive_key_t][]archive_t{}
				state.skipped = true
			}
		}

		shape := getShapeInOrder(state.numberOfShapes)
		state.fallingShape = &falling_shape_t{
			shape:  shape,
			offset: getStartingOffset(state, shape),
		}
		state.numberOfShapes++
		return state
	}
	nextDirection := state.input[state.inputIdx]
	state.inputIdx = (state.inputIdx + 1) % len(state.input)

	var shapeNextOffset point_t
	directionOffset := directionToOffset(nextDirection)
	directionIsOk := false
	if canShapeFitNextOffset(state, directionOffset) {
		shapeNextOffset = directionOffset
		directionIsOk = true
	} else {
		shapeNextOffset = point_t{x: 0, y: 0}
	}

	gravityOffset := point_t{x: 1, y: 0}
	shapeNextOffset = add(shapeNextOffset, gravityOffset)
	if canShapeFitNextOffset(state, shapeNextOffset) {
		state.fallingShape.offset = add(state.fallingShape.offset, shapeNextOffset)
	} else {
		if directionIsOk {
			state.fallingShape.offset = add(state.fallingShape.offset, directionOffset)
		}
		state = addFallingShapeRocksToMap(state)
		state.peak = findNewHeight(state)
		state.fallingShape = nil
	}

	return state
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func pointInShape(point point_t, shape falling_shape_t) bool {
	for _, p := range shape.shape.points {
		movedP := add(p, shape.offset)
		if movedP.x == point.x && movedP.y == point.y {
			return true
		}
	}
	return false
}

func printState(state state_t, limit int64) {
	fmt.Println("STATE:")
	for x := state.peak - 10; x <= state.peak-10+limit; x++ {
		for y := int64(-1); y <= state.width; y++ {
			p := point_t{x: x, y: y}
			if state.rocks[p] {
				fmt.Print("#")
				continue
			}
			if state.fallingShape != nil && pointInShape(p, *state.fallingShape) {
				fmt.Print("@")
				continue
			}
			if y == -1 {
				fmt.Print("|")
				continue
			}
			if y == state.width {
				fmt.Print("|")
				continue
			}
			if x == 1 {
				fmt.Print("-")
				continue
			}
			fmt.Print(".")
		}
		fmt.Println()
	}
}

func parseInput(input string) []byte {
	return []byte(input)
}

func stepForP1(state state_t) state_t {
	for {
		state = step(state, true)
		if state.numberOfShapes == state.target+1 {
			return state
		}
	}
}

func solveP1(input string, lines []string) string {
	initState := state_t{
		rocks:          map[point_t]bool{},
		numberOfShapes: 0,
		fallingShape:   nil,
		width:          7,
		peak:           1,
		input:          parseInput(input),
		inputIdx:       0,
		archives:       map[archive_key_t][]archive_t{},
		skipped:        false,
		target:         TARGET_P1,
	}
	state := stepForP1(initState)
	return fmt.Sprintf("%v", state.peak*(-1)+1)
}

func stepForP2(state state_t) state_t {
	for {
		state = step(state, true)
		if state.numberOfShapes == state.target+1 {
			return state
		}
	}
}

func solveP2(input string, lines []string) string {
	initState := state_t{
		rocks:          map[point_t]bool{},
		numberOfShapes: 0,
		fallingShape:   nil,
		width:          7,
		peak:           1,
		input:          parseInput(input),
		inputIdx:       0,
		archives:       map[archive_key_t][]archive_t{},
		skipped:        false,
		target:         TARGET_P2,
	}
	state := stepForP2(initState)
	return fmt.Sprintf("%v", state.peak*(-1)+1)
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
