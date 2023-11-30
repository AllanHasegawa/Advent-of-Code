package main

import (
    "fmt"
    "os"
    "strings"
    // "strconv"
    // "sort"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func solveP1(input string, lines []string) string {
    return "TODO"
}

func solveP2(input string, lines []string) string {
    return "TODO"
}

func main() {
    dat, err := os.ReadFile("tinput.txt")
    check(err)

    input := string(dat)
    lines := strings.Split(input, "\n")

    fmt.Println("Solution P1:")
    fmt.Println(solveP1(input, lines))
    fmt.Println("")
    fmt.Println("Solution P2:")
    fmt.Println(solveP2(input, lines))
}
