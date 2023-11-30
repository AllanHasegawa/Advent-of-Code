package main

import (
    "fmt"
    "os"
    "strings"
    "strconv"
    // "sort"
)

type sections struct {
    low int
    high int
}

type input struct {
    first sections
    second sections
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func parseLines(lines []string) []input {
    parseIndividualElf := func(raw string) sections {
        sectionsRaw := strings.Split(raw, "-")
        lowParsed, err := strconv.Atoi(sectionsRaw[0])
        check(err)
        highParsed, err := strconv.Atoi(sectionsRaw[1])
        check(err)
        return sections{low: lowParsed, high: highParsed}
    }

    var inputs []input
    
    for i := 0; i < len(lines); i++ {
        line := lines[i]
        individualElves := strings.Split(line, ",")
        firstElf := parseIndividualElf(individualElves[0])
        secondElf := parseIndividualElf(individualElves[1])
        inputs = append(inputs, input{first: firstElf, second: secondElf})
    }

    return inputs
}

func doesSectionAFullyOverlapsB(a sections, b sections) bool {
    return a.low <= b.low && a.high >= b.high
}

func solveP1(input string, lines []string) string {
    var inputs = parseLines(lines)
    var fullyOverlaps = 0

    for i := 0; i < len(inputs); i++ {
        input := inputs[i]
        if doesSectionAFullyOverlapsB(input.first, input.second) {
            fullyOverlaps++
            continue
        } else if doesSectionAFullyOverlapsB(input.second, input.first) {
            fullyOverlaps++
        }
    }

    return fmt.Sprintf("%d", fullyOverlaps)
}


/**
A overlaps to the right.
         3   7 
a :      |---|
b :     |---|
        2   6
*/
func doesSectionAOverlapsWithB(a sections, b sections) bool {
    return b.low <= a.high && b.high >= a.low
}

func solveP2(input string, lines []string) string {
    var inputs = parseLines(lines)
    var overlaps = 0

    for i := 0; i < len(inputs); i++ {
        var input = inputs[i]
        if doesSectionAOverlapsWithB(input.first, input.second) {
            overlaps++
        }
    }

    return fmt.Sprintf("%d", overlaps)
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
