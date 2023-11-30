package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	// "strconv"
	// "sort"
)

const GRID_SIZE = 22

type thing_t int

const (
	NOTHING thing_t = 0
	LAVA    thing_t = 1
	STEAM   thing_t = 2
)

type voxels_t [GRID_SIZE][GRID_SIZE][GRID_SIZE]thing_t

type coord_t struct {
	x, y, z int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func add(a coord_t, b coord_t) coord_t {
	return coord_t{x: a.x + b.x, y: a.y + b.y, z: a.z + b.z}
}

func parseInput(lines []string) []coord_t {
	coords := []coord_t{}
	for _, l := range lines {
		splits := strings.Split(l, ",")
		x, ok := strconv.Atoi(splits[0])
		check(ok)
		y, ok := strconv.Atoi(splits[1])
		check(ok)
		z, ok := strconv.Atoi(splits[2])
		check(ok)
		coords = append(coords, coord_t{x: x, y: y, z: z})
	}
	return coords
}

func fillGrid(coords []coord_t) voxels_t {
	voxels := voxels_t{}
	for _, c := range coords {
		voxels[c.x][c.y][c.z] = LAVA
	}
	return voxels
}

var CUBE_SIDES []coord_t = []coord_t{
	coord_t{x: -1, y: 0, z: 0}, // left side
	coord_t{x: 0, y: 0, z: 1},  // front
	coord_t{x: 1, y: 0, z: 0},  // right side
	coord_t{x: 0, y: 0, z: -1}, // back side
	coord_t{x: 0, y: 1, z: 0},  // top side
	coord_t{x: 0, y: -1, z: 0}, // lower side
}

func isCoordValid(c coord_t) bool {
	return c.x >= 0 && c.x < GRID_SIZE &&
		c.y >= 0 && c.y < GRID_SIZE &&
		c.z >= 0 && c.z < GRID_SIZE
}

func get(grid voxels_t, c coord_t) thing_t {
	return grid[c.x][c.y][c.z]
}

func countSidesVoxel(coord coord_t, grid voxels_t, countAirPockets bool) int {
	voxel := get(grid, coord)
	if voxel == NOTHING || voxel == STEAM {
		return 0
	}

	totalSides := 6
	for _, side := range CUBE_SIDES {
		c := add(side, coord)

		var adjacentVoxel thing_t
		if !isCoordValid(c) {
			adjacentVoxel = STEAM
		} else {
			adjacentVoxel = get(grid, c)
		}
		if !countAirPockets && adjacentVoxel == NOTHING {
			totalSides--
		}

		if adjacentVoxel == LAVA {
			totalSides--
		}
	}
	return totalSides
}

func countSides(grid voxels_t, countAirPockets bool) int {
	totalSides := 0
	for x := 0; x < GRID_SIZE; x++ {
		for y := 0; y < GRID_SIZE; y++ {
			for z := 0; z < GRID_SIZE; z++ {
				c := coord_t{x: x, y: y, z: z}
				totalSides += countSidesVoxel(c, grid, countAirPockets)
			}
		}
	}
	return totalSides
}

func spreadSteam(c coord_t, grid *voxels_t) {
	if !isCoordValid(c) {
		return
	}

	voxel := get(*grid, c)
	if voxel != NOTHING {
		return
	}
	grid[c.x][c.y][c.z] = STEAM

	for _, side := range CUBE_SIDES {
		adjacentVoxel := add(side, c)
		spreadSteam(adjacentVoxel, grid)
	}
}

func fillGridWithSteam(grid *voxels_t) {
	// Front and Back
	for x := 0; x < GRID_SIZE; x++ {
		for y := 0; y < GRID_SIZE; y++ {
			z := 0
			c := coord_t{x: x, y: y, z: z}
			spreadSteam(c, grid)
			z = GRID_SIZE - 1
			c.z = z
			spreadSteam(c, grid)
		}
	}

	// Left and Right
	for z := 0; z < GRID_SIZE; z++ {
		for y := 0; y < GRID_SIZE; y++ {
			x := 0
			c := coord_t{x: x, y: y, z: z}
			spreadSteam(c, grid)
			x = GRID_SIZE - 1
			c.x = x
			spreadSteam(c, grid)
		}
	}

	// Top and Bottom
	for x := 0; x < GRID_SIZE; x++ {
		for z := 0; z < GRID_SIZE; z++ {
			y := 0
			c := coord_t{x: x, y: y, z: z}
			spreadSteam(c, grid)
			y = GRID_SIZE - 1
			c.y = y
			spreadSteam(c, grid)
		}
	}
}

func solveP1(input string, lines []string) string {
	coords := parseInput(lines)
	grid := fillGrid(coords)
	sides := countSides(grid, true)
	return fmt.Sprintf("%v", sides)
}

func solveP2(input string, lines []string) string {
	coords := parseInput(lines)
	grid := fillGrid(coords)
	fillGridWithSteam(&grid)

	sides := countSides(grid, false)
	return fmt.Sprintf("%v", sides)
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
