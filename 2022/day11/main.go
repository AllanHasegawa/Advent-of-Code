package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
    "math/big"
	// "sort"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type monkey_t struct {
    items []*big.Int
    operation func(*big.Int) *big.Int
    testDivisible *big.Int
    testTrueTargetIdx int
    testFalseTargetIdx int
    inspections int
}

type state_t []monkey_t

func parseLines(lines []string) state_t {
    var state state_t

    for i := 0; i < len(lines); i += 7 {
        monkey := monkey_t{}
        // 0 line is useless
        monkey.items = parseStartingItems(lines[i+1])
        monkey.operation = parseOperation(lines[i+2])
        monkey.testDivisible = parseTestDivisible(lines[i+3])
        monkey.testTrueTargetIdx = parseTestTrueTargetIdx(lines[i+4])
        monkey.testFalseTargetIdx = parseTestFalseTargetIdx(lines[i+5])
        // 6 line is empty

        state = append(state, monkey)
    }

    return state
}

func parseStartingItems(line string) []*big.Int {
    prefix := "  Starting items: "
    dataRaw := line[len(prefix):]
    split := strings.Split(dataRaw, ", ")
    items := [](*big.Int){}
    
    for i := 0; i < len(split); i++ {
        worryLevel, err := strconv.Atoi(split[i])
        check(err)
        items = append(items, big.NewInt(int64(worryLevel)))
    }

    return items
}

func parseOperand(raw string) func(*big.Int) *big.Int {
    if raw == "old" {
        return func(i *big.Int) *big.Int { return i }
    }
    return func(*big.Int) *big.Int { 
        operand, err := strconv.Atoi(raw)
        check(err)
        return big.NewInt(int64(operand))
    }
}

func parseOperation(line string) func(*big.Int) *big.Int {
    prefix := "  Operation: new = "
    dataRaw := line[len(prefix):]
    split := strings.Split(dataRaw, " ")
    aOperand := parseOperand(split[0])
    bOperand := parseOperand(split[2])
    return func(old *big.Int) *big.Int {
        if split[1] == "*" {
            a := aOperand(old)
            return a.Mul(a, bOperand(old))
        } else {
            a := aOperand(old)
            return a.Add(a, bOperand(old))
        }
    }
}

func parseTestDivisible(line string) *big.Int {
    prefix := "  Test: divisible by "
    dataRaw := line[len(prefix):]
    divisible, err := strconv.Atoi(dataRaw)
    check(err)
    return big.NewInt(int64(divisible))
}

func parseTestTrueTargetIdx(line string) int {
    prefix := "    If true: throw to monkey "
    dataRaw := line[len(prefix):]
    idx, err := strconv.Atoi(dataRaw)
    check(err)
    return idx
}

func parseTestFalseTargetIdx(line string) int {
    prefix := "    If false: throw to monkey "
    dataRaw := line[len(prefix):]
    idx, err := strconv.Atoi(dataRaw)
    check(err)
    return idx
}

func monkeyTurn(state state_t, monkeyIdx int, divideBy *big.Int, modBy *big.Int) state_t {
    monkey := &state[monkeyIdx]
    monkey.inspections += len(monkey.items)
    modResult := big.NewInt(0)

    for i := 0; i < len(monkey.items); i++ {
        newWorry := monkey.operation(monkey.items[i])
        if divideBy != nil {
            newWorry.Div(newWorry, divideBy)
        }
        if modBy != nil {
            newWorry.Mod(newWorry, modBy)
        }
        var monkeyTargetIdx int
        if modResult.Mod(newWorry, monkey.testDivisible).Int64()  == 0 {
            // fmt.Println(monkeyIdx, "WAS TRUE", newWorry)
            monkeyTargetIdx = monkey.testTrueTargetIdx
        } else {
            // fmt.Println(monkeyIdx, "WAS FALSE", newWorry)
            monkeyTargetIdx = monkey.testFalseTargetIdx
        }
        state[monkeyTargetIdx].items = append(state[monkeyTargetIdx].items, newWorry)
    }
    monkey.items = [](*big.Int){}
    return state
}

func monkeyRound(state state_t, divideBy *big.Int, modBy *big.Int) state_t {
    for m := 0; m < len(state); m++ {
        state = monkeyTurn(state, m, divideBy, modBy)
    }
    return state
}

func solveP1(input string, lines []string) string {
    state := parseLines(lines)
    monkeyInspections := [](int){}

    for i := 0; i < 20; i++ {
        state = monkeyRound(state, big.NewInt(3), nil)
    }

    for i := 0; i < len(state); i++ {
        monkeyInspections = append(monkeyInspections, state[i].inspections)
    }

    sort.Slice(monkeyInspections, func(i, j int) bool {
        return monkeyInspections[i] > monkeyInspections[j]
    })

    return fmt.Sprintf("%v", monkeyInspections[0]*monkeyInspections[1])
}

func solveP2(input string, lines []string) string {
    state := parseLines(lines)
    monkeyInspections := [](int){}
    maxDivisor := big.NewInt(1)

    for i := 0; i < len(state); i++ {
        maxDivisor.Mul(state[i].testDivisible, maxDivisor)
    }

    for i := 0; i < 10_000; i++ {
        state = monkeyRound(state, nil, maxDivisor)
    }

    for i := 0; i < len(state); i++ {
        monkeyInspections = append(monkeyInspections, state[i].inspections)
    }

    sort.Slice(monkeyInspections, func(i, j int) bool {
        return monkeyInspections[i] > monkeyInspections[j]
    })

    return fmt.Sprintf("%v", monkeyInspections[0]*monkeyInspections[1])
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
