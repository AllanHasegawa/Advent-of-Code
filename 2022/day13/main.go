package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

type operand_t interface {}

type pair_t struct {
    left operand_t
    right operand_t
}

type int_operand_t struct {
    value int
}

type list_operand_t struct {
    values []operand_t
}

// Returns the operand and the current lookahead index
func parseInt(raw string, idx int) (operand_t, int) {
    nextNonIntCharIdx := -1

    for idx, c := range raw[idx:] {
        if c == ']' || c == ',' {
            nextNonIntCharIdx = idx
            break
        }
    }

    if nextNonIntCharIdx == -1 {
        panic("Unknown state, couldn't find end of int for: " + raw)
    }

    intValue, err := strconv.Atoi(raw[idx:idx+nextNonIntCharIdx])
    check(err)

    return int_operand_t{value: intValue}, idx+nextNonIntCharIdx
}

// Returns the operand and the current lookahead index
func parseList(raw string, idx int) (operand_t, int) {
    opsList := []operand_t{}
    currentIndex := idx

    out:
    for ;; {
        if currentIndex >= len(raw) {
            break
        }
        switch raw[currentIndex] {
        case ']': {
            currentIndex++
            break out
        }
        case ',': {
            currentIndex++
            op, idxP := parseNextChar(raw, currentIndex)
            currentIndex = idxP
            opsList = append(opsList, op)
            continue out
        }
        case '[': {
            currentIndex++
            listOp, idxP := parseList(raw, currentIndex)
            currentIndex = idxP
            opsList = append(opsList, listOp)
            continue out
        }
        default: {
            intOp, idxP := parseInt(raw, currentIndex)
            currentIndex = idxP
            opsList = append(opsList, intOp)
            continue out
        }
        }
    }

    return list_operand_t{values: opsList}, currentIndex
}

// Returns the operand and the current lookahead index
func parseNextChar(raw string, idx int) (operand_t, int) {
    if len(raw) == 0 {
        panic("Unknown state, raw is zero.")
    }
    if len(raw) <= idx {
        panic(fmt.Sprintf("Raw len is too big: %v/%v", len(raw), idx))
    }

    if raw[idx] == '[' {
        return parseList(raw, idx+1)
    } else {
        return parseInt(raw, idx)
    }
}

func parsePacket(line string) operand_t {
    op, _ := parseNextChar(line, 0)
    return op
}

func parsePacketPairs(lines []string) []pair_t {
    pairs := []pair_t{}
    for i := 0; i < len(lines); i+=3 {
        packetA := parsePacket(lines[i])
        packetB := parsePacket(lines[i+1])
        pair := pair_t{left: packetA, right: packetB}
        pairs = append(pairs, pair)
    }
    return pairs
}

func Max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func compare(left operand_t, right operand_t) int {
    switch left.(type) {
    case int_operand_t: {
        leftValue := left.(int_operand_t).value
        switch right.(type) {
        case int_operand_t: {
            rightValue := right.(int_operand_t).value
            if leftValue == rightValue {
                return 0
            }
            if leftValue < rightValue {
                return -1
            }
            return 1
        }
        case list_operand_t: {
            leftOpP := list_operand_t{values: []operand_t{int_operand_t{value: leftValue}}}
            return compare(leftOpP, right)
        }
        }
    }
    case list_operand_t: {
        switch right.(type) {
        case int_operand_t: {
            rightOpP := list_operand_t{values: []operand_t{int_operand_t{value: right.(int_operand_t).value}}}
            return compare(left, rightOpP)
        }
        case list_operand_t: {
            leftList := left.(list_operand_t).values
            rightList := right.(list_operand_t).values
            maxIdx := Max(len(leftList), len(rightList))
            for i := 0; i < maxIdx; i++ {
                if i == len(leftList) {
                    return -1
                }
                if i == len(rightList) {
                    return 1
                }
                compareResult := compare(leftList[i], rightList[i])
                if compareResult == 0 {
                    continue
                }
                return compareResult
            }
            return 0
        }
        }
    }
    }
    fmt.Printf("%T, %T\n", left, right)
    panic(fmt.Sprintf("Comparison with unknown types?\n%v\n%v", left, right))
}

func solveP1(input string, lines []string) string {
    pairs := parsePacketPairs(lines)

    sumIdx := 0
    for i, p := range pairs {
        compareResult := compare(p.left, p.right)
        if compareResult <= 0 {
            sumIdx += i+1
        }
    }
    return fmt.Sprintf("%v", sumIdx)
}

func solveP2(input string, lines []string) string {
    pairs := parsePacketPairs(lines)
    dividerA := parsePacket("[[2]]")
    dividerB := parsePacket("[[6]]")
    idxDividerA := -1
    idxDividerB := -1
    allOps := []operand_t{}

    for _, p := range pairs {
        allOps = append(allOps, p.left)
        allOps = append(allOps, p.right)
    }
    allOps = append(allOps, dividerA)
    allOps = append(allOps, dividerB)

    sort.Slice(allOps, func(i, j int) bool {
        return compare(allOps[i], allOps[j]) < 0
    })

    // This is VERY inefficient and manual comparing the "deep" values of packets.
    // Wouldn't scale for bigger comparisons and test cases.
    // Ideally would compare references, but, Golang was hard enough to get here :(
    for idx, ops := range allOps {
        switch ops.(type) {
        case list_operand_t: {
            opsValues := ops.(list_operand_t).values
            if len(opsValues) == 1 {
                listOp, ok := opsValues[0].(list_operand_t)
                if ok {
                    listOpsValues := listOp.values
                    if len(listOpsValues) == 1 {
                        intOp, ok := listOpsValues[0].(int_operand_t)
                        if ok { 
                            if intOp.value == 2 {
                                idxDividerA = idx + 1
                            }
                            if intOp.value == 6 {
                                idxDividerB = idx + 1
                            }
                        }
                    }
                }
            }
        }
        }
    }

    return fmt.Sprintf("%v", idxDividerA * idxDividerB)
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
