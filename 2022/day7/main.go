package main

import (
	"fmt"
	"os"
	"sort"

	// "reflect"
	"strconv"
	"strings"
	// "sort"
)

type state struct {
	root        *file_dir
	currentPath *file_dir
}

type file struct {
	size int
	name string
}

type file_dir struct {
	name     string
	files    []file
	parent   *file_dir
	children []*file_dir
}

type action interface {
	process(*state)
}

type cd_cmd struct {
	param string
}

type ls_cmd struct{}

type dir_list struct {
	dirName string
}

type file_list struct {
	fileSize int
	fileName string
}

type line struct {
	cmd_or_list action
}

func (c cd_cmd) process(s *state) {
	if c.param == "/" {
		s.currentPath = s.root
		return
	}
	if c.param == ".." {
		s.currentPath = s.currentPath.parent
		return
	}

	children := s.currentPath.children
	for i := 0; i < len(children); i++ {
		child := children[i]
		if child.name == c.param {
			s.currentPath = child
			return
		}
	}
	panic(fmt.Sprintf("Command not processed: %s", c))
}

func (c ls_cmd) process(s *state) {
	// Nothing to do
}

func (c dir_list) process(s *state) {
	newDir := file_dir{name: c.dirName, files: [](file){}, parent: s.currentPath, children: [](*file_dir){}}
	s.currentPath.children = append(s.currentPath.children, &newDir)
}

func (c file_list) process(s *state) {
	newFile := file{name: c.fileName, size: c.fileSize}
	s.currentPath.files = append(s.currentPath.files, newFile)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func processLines(lines []line) state {
	rootDir := file_dir{name: "root", files: [](file){}, parent: nil, children: [](*file_dir){}}
	state := state{currentPath: &rootDir, root: &rootDir}

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		line.cmd_or_list.process(&state)
	}

	return state
}

func printState(s *state) {
	printDir(s.root, 0)
}

func printDir(d *file_dir, depth int) {
	fmt.Println(d.name)

	for i := 0; i < len(d.files); i++ {
		file := d.files[i]
		printDepth(depth, 0)
		fmt.Print("|  ")
		fmt.Println(file.name, file.size)
	}

	for i := 0; i < len(d.children); i++ {
		dir := d.children[i]
		printDepth(depth, 0)
		fmt.Print("|> ")
		printDir(dir, depth+1)
	}
}

func printDepth(depth int, offset int) {
	for i := 0; i < depth*4+offset; i++ {
		fmt.Print(" ")
	}
}

func dirSize(d *file_dir) int {
	allFilesSize := 0
	allChildrenSize := 0

	for i := 0; i < len(d.files); i++ {
		allFilesSize += d.files[i].size
	}

	for i := 0; i < len(d.children); i++ {
		allChildrenSize += dirSize(d.children[i])
	}

	return allFilesSize + allChildrenSize
}

func findAllDirLessThanACertainSize(current *file_dir, dirsSoFar *[](*file_dir), size int) {
	if dirSize(current) <= size {
		*dirsSoFar = append(*dirsSoFar, current)
	}

	for i := 0; i < len(current.children); i++ {
		findAllDirLessThanACertainSize(current.children[i], dirsSoFar, size)
	}
}

func findAllDirMoreThanACertainSize(current *file_dir, dirsSoFar *[](*file_dir), size int) {
	if dirSize(current) >= size {
		*dirsSoFar = append(*dirsSoFar, current)
	}

	for i := 0; i < len(current.children); i++ {
		findAllDirMoreThanACertainSize(current.children[i], dirsSoFar, size)
	}
}

func parseLine(rawLine string) line {
	var cmd_or_list action
	split := strings.Split(rawLine, " ")
	if split[0] == "$" {
		if split[1] == "cd" {
			cmd_or_list = cd_cmd{param: split[2]}
		} else if split[1] == "ls" {
			cmd_or_list = ls_cmd{}
		} else {
			panic("Unknown command: " + split[1])
		}
	} else {
		if split[0] == "dir" {
			cmd_or_list = dir_list{dirName: split[1]}
		} else if len(split) == 2 {
			fileSize, err := strconv.Atoi(split[0])
			check(err)
			cmd_or_list = file_list{fileSize: fileSize, fileName: split[1]}
		} else {
			panic("Unknown line: " + rawLine)
		}
	}

	return line{cmd_or_list: cmd_or_list}
}

func parseInput(rawLines []string) []line {
	var lines []line
	for i := 0; i < len(rawLines); i++ {
		lines = append(lines, parseLine(rawLines[i]))
	}
	return lines
}

func solveP1(input string, lines []string) string {
	parsedLines := parseInput(lines)
	state := processLines(parsedLines)
	size := 100001
	var dirsBelowSize [](*file_dir)
	findAllDirLessThanACertainSize(state.root, &dirsBelowSize, size)

	sum := 0
	for i := 0; i < len(dirsBelowSize); i++ {
		sum += dirSize(dirsBelowSize[i])
	}

	return fmt.Sprintf("%d", sum)
}

func solveP2(input string, lines []string) string {
	parsedLines := parseInput(lines)
	state := processLines(parsedLines)
    totalDiskSize := 70000000
    freeSpaceRequired := 30000000
    usedSize := dirSize(state.root)
    freeSpace := totalDiskSize - usedSize
    spaceThatNeedsToGo := freeSpaceRequired - freeSpace

    var dirsAboveSize [](*file_dir)
    findAllDirMoreThanACertainSize(state.root, &dirsAboveSize, spaceThatNeedsToGo)

    var dirsAboveSizeSize []int
    for i := 0; i < len(dirsAboveSize); i++ {
        dirsAboveSizeSize = append(dirsAboveSizeSize, dirSize(dirsAboveSize[i]))
    }
    sort.Ints(dirsAboveSizeSize)

	return fmt.Sprintf("%d", dirsAboveSizeSize[0])
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
