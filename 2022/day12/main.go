package main

import (
	"fmt"
	"math"
	"os"
	"strings"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	// "strconv"
	// "sort"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type world_t struct {
	start     *tile_t
	end       *tile_t
	heightMap map[int]map[int]*tile_t
}

type tile_t struct {
	x, y   int
	height int
	world  *world_t
}

func parseWorld(lines []string) world_t {
	world := world_t{}
	var heightMap = map[int]map[int]*tile_t{}

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if heightMap[i] == nil {
			heightMap[i] = map[int]*tile_t{}
		}

		for j := 0; j < len(line); j++ {
			height := line[j]
			tile := &tile_t{x: i, y: j, world: &world, height: int(height)}

			if height == 'S' {
				world.start = tile
				tile.height = 'a'
			} else if height == 'E' {
				world.end = tile
				tile.height = 'z'
			}

			heightMap[i][j] = tile
		}
	}
	world.heightMap = heightMap
	return world
}

func findLowestFScoreNode(openSet map[*tile_t]bool, fScore map[*tile_t]int) *tile_t {
	lowestFScore := math.MaxInt
	var lowestFScoreNode *tile_t = nil

	for t := range openSet {
		score := fScore[t]
		if score < lowestFScore {
			lowestFScoreNode = t
            lowestFScore = score
		}
	}

	return lowestFScoreNode
}

func getXScore(scoreMap map[*tile_t]int, key *tile_t) int {
	if v, ok := scoreMap[key]; ok {
		return v
	} else {
		scoreMap[key] = math.MaxInt
		return math.MaxInt
	}
}

func rebuildPath(cameFrom map[*tile_t]*tile_t, current *tile_t) []*tile_t {
	totalPath := []*tile_t{current}

	for {
		if !slices.Contains(maps.Keys(cameFrom), current) {
			break
		}
		current = cameFrom[current]
		totalPath = append(totalPath, current)
	}

	return totalPath
}

func aStarSeach(start *tile_t,
	goal func(*tile_t) bool,
	heuristic func(*tile_t) int,
	neighbours func(*tile_t) []*tile_t,
	moveCost func(*tile_t, *tile_t) int) []*tile_t {
	openSet := map[*tile_t]bool{start: true}
	cameFrom := map[*tile_t]*tile_t{}

	gScore := map[*tile_t]int{start: 0}
	fScore := map[*tile_t]int{start: heuristic(start)}

	for {
		if len(openSet) == 0 {
			break
		}
		current := findLowestFScoreNode(openSet, fScore)
		if goal(current) {
			return rebuildPath(cameFrom, current)
		}
		delete(openSet, current)

		for _, neighbour := range neighbours(current) {
			currentGScore := getXScore(gScore, current)
			tentativeGScore := currentGScore + moveCost(current, neighbour)
			if tentativeGScore < getXScore(gScore, neighbour) {
				cameFrom[neighbour] = current
				gScore[neighbour] = tentativeGScore
				fScore[neighbour] = tentativeGScore + heuristic(neighbour)
				openSet[neighbour] = true
			}
		}
	}

	return []*tile_t{}
}

func pathNeighbors(t *tile_t, goingUp bool) []*tile_t {
	neighbors := []*tile_t{}
	for _, offset := range [][]int{
		{-1, 0},
		{1, 0},
		{0, -1},
		{0, 1},
	} {
        tx := t.x+offset[0]
        if t.world.heightMap[tx] == nil {
            continue
        }
        ty := t.y+offset[1]
		target := t.world.heightMap[tx][ty]
        if target == nil {
            continue
        }

		isAllowed := (goingUp && target.height <= t.height+1) ||
			(!goingUp && target.height >= t.height-1)

		if isAllowed {
			neighbors = append(neighbors, target)
		}
	}
	return neighbors
}

// Heuristic function is manhattan distance
func pathEstimatedCost(from *tile_t, to *tile_t) int {
	absX := to.x - from.x
	if absX < 0 {
		absX = -absX
	}
	absY := to.y - from.y
	if absY < 0 {
		absY = -absY
	}
	return absX + absY
}

func solveP1(input string, lines []string) string {
	world := parseWorld(lines)
	path := aStarSeach(
		world.start,
		func(t *tile_t) bool { return t == world.end },
		func(t *tile_t) int { return pathEstimatedCost(t, world.end) },
		func(t *tile_t) []*tile_t { return pathNeighbors(t, true) },
		func(t *tile_t, to *tile_t) int { return 1 },
	)

	return fmt.Sprintf("%v", len(path)-1)
}

func solveP2(input string, lines []string) string {
	world := parseWorld(lines)
	path := aStarSeach(
		world.end,
		func(t *tile_t) bool { return t.height == int('a') },
        // we make the heuristic function admissible (but not optimal) to return the shortest path (but slowly)
		func(t *tile_t) int { return 1 },
		func(t *tile_t) []*tile_t { return pathNeighbors(t, false) },
		func(t *tile_t, to *tile_t) int { return 1 },
	)

	return fmt.Sprintf("%v", len(path)-1)
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
