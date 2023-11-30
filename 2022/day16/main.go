package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/maps"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type caves_t map[string]cave_t

type cave_t struct {
	label    string
	flowRate int
	edges    []string
	moveCost map[string]int
}

type state_t struct {
	closedValves map[string]bool
	position     string
	currentFlow  int
	minutesLeft  int
	valvesOpened []string
}

func initState(caves caves_t) state_t {
	closedValves := map[string]bool{}

	for _, c := range caves {
		if c.flowRate > 0 {
			closedValves[c.label] = true
		}
	}

	return state_t{
		closedValves: closedValves,
		currentFlow:  0,
		minutesLeft:  30,
		position:     "AA",
		valvesOpened: []string{},
	}
}

type opened_valve_t struct {
	label    string
	minute   int
	agentIdx int
}

type state2_t struct {
	agents          []agent_t
	minutesPassed   int
	openedValves    []opened_valve_t
	closedValves    map[string]bool
	similuationTime int
	flowSoFar       int
}

type agent_t struct {
	position string
	timeBusy int
}

type action_t int

const (
	Idle        action_t = 0
	Wait                 = 1
	MoveAndOpen          = 2
)

type agent_action_t struct {
	option      action_t
	moveToLabel string
}

func initState2(caves caves_t) state2_t {
	closedValves := map[string]bool{}
	for _, c := range caves {
		if c.flowRate > 0 {
			closedValves[c.label] = true
		}
	}

	agents := []agent_t{
		{
			position: "AA",
			timeBusy: 0,
		},
		{
			position: "AA",
			timeBusy: 0,
		},
	}
	return state2_t{
		agents:          agents,
		minutesPassed:   0,
		openedValves:    []opened_valve_t{},
		flowSoFar:       0,
		closedValves:    closedValves,
		similuationTime: -1,
	}
}

func openValve(state state_t, caves caves_t) state_t {
	cave := caves[state.position]
	newClosedValves := maps.Clone(state.closedValves)
	delete(newClosedValves, cave.label)
	newMinutesLeft := state.minutesLeft - 1
	newFlow := state.currentFlow + newMinutesLeft*cave.flowRate
	return state_t{
		closedValves: newClosedValves,
		currentFlow:  newFlow,
		minutesLeft:  newMinutesLeft,
		position:     state.position,
		valvesOpened: append(state.valvesOpened, cave.label),
	}
}

func openValve2(agentIdx int, state state2_t, caves caves_t) state2_t {
	agent := state.agents[agentIdx]
	agent.timeBusy += 1

	cave := caves[agent.position]
	newClosedValves := maps.Clone(state.closedValves)
	delete(newClosedValves, cave.label)
	newOpenedValves := make([]opened_valve_t, len(state.openedValves))
	copy(newOpenedValves, state.openedValves)
	newOpenedValves = append(newOpenedValves, opened_valve_t{
		label:    cave.label,
		minute:   state.minutesPassed + agent.timeBusy,
		agentIdx: agentIdx,
	})
	timeToOpen := state.minutesPassed + agent.timeBusy
	flowSoFar := state.flowSoFar + (state.similuationTime - timeToOpen) * cave.flowRate
	return state2_t{
		agents:          copyAgents(agent, agentIdx, state.agents),
		minutesPassed:   state.minutesPassed,
		openedValves:    newOpenedValves,
		flowSoFar:       flowSoFar,
		closedValves:    newClosedValves,
		similuationTime: state.similuationTime,
	}
}

func move(toCaveLabel string, state state_t, caves caves_t) state_t {
	cave := caves[state.position]
	cost := cave.moveCost[toCaveLabel]
	minutesLeft := state.minutesLeft - cost
	return state_t{
		closedValves: state.closedValves,
		currentFlow:  state.currentFlow,
		minutesLeft:  minutesLeft,
		position:     toCaveLabel,
		valvesOpened: state.valvesOpened,
	}
}

func move2(toCaveLabel string, agentIdx int, state state2_t, caves caves_t) state2_t {
	agent := state.agents[agentIdx]
	cave := caves[agent.position]

	cost := cave.moveCost[toCaveLabel]
	agent.position = toCaveLabel
	agent.timeBusy += cost
	return state2_t{
		agents:          copyAgents(agent, agentIdx, state.agents),
		minutesPassed:   state.minutesPassed,
		openedValves:    state.openedValves,
		flowSoFar:       state.flowSoFar,
		closedValves:    state.closedValves,
		similuationTime: state.similuationTime,
	}
}

func parseLine(line string) cave_t {
	prefix := "Valve "
	mid := "; tunnels lead to valves "
	mid2 := "; tunnel leads to valve "
	flowDiv := " has flow rate="

	midSplits := strings.Split(line, mid)
	if len(midSplits) == 1 {
		midSplits = strings.Split(line, mid2)
	}

	flowSplits := strings.Split(midSplits[0][len(prefix):], flowDiv)

	label := flowSplits[0]
	flowRate, err := strconv.Atoi(flowSplits[1])
	check(err)

	edges := strings.Split(midSplits[1], ", ")

	return cave_t{label: label, flowRate: flowRate, edges: edges}
}

func parseInput(lines []string) map[string]cave_t {
	caves := map[string]cave_t{}

	for _, l := range lines {
		cave := parseLine(l)
		caves[cave.label] = cave
	}

	return caves
}

func findMoveCost(label string, caves caves_t) map[string]int {
	visited := map[string]bool{label: true}
	costs := map[string]int{label: 0}
	toVisit := caves[label].edges
	currentCost := 1

	for {
		if len(toVisit) == 0 {
			break
		}

		newToVisit := []string{}
		for _, v := range toVisit {
			if visited[v] {
				continue
			}
			visited[v] = true
			costs[v] = currentCost
			newToVisit = append(newToVisit, caves[v].edges...)
		}
		currentCost += 1
		toVisit = newToVisit
	}

	return costs
}

func findAllMoveCost(caves caves_t) {
	for _, label := range maps.Keys(caves) {
		costs := findMoveCost(label, caves)
		c := caves[label]
		c.moveCost = costs
		caves[label] = c
	}
}

func step(states []state_t, idleStates []state_t, caves caves_t) ([]state_t, []state_t, bool) {
	didChange := false
	newStates := []state_t{}
	newIdleStates := []state_t{}
	newIdleStates = append(newIdleStates, idleStates...)

	for _, state := range states {
		if state.minutesLeft <= 0 {
			newIdleStates = append(newIdleStates, state)
			continue
		}

		currentCave := caves[state.position]
		didMove := false

		for _, targetCaveLabel := range maps.Keys(state.closedValves) {
			moveCost := currentCave.moveCost[targetCaveLabel]
			if (moveCost + 1) > state.minutesLeft {
				continue
			}
			stateAfterMove := move(targetCaveLabel, state, caves)
			stateAfterOpen := openValve(stateAfterMove, caves)
			didMove = true
			newStates = append(newStates, stateAfterOpen)
			didChange = true
		}
		if !didMove {
			newIdleStates = append(newIdleStates, state)
		}
	}

	return newStates, newIdleStates, didChange
}

func stepUntilEnd(states []state_t, caves caves_t) []state_t {
	ok := false
	idleStates := []state_t{}

	for {
		states, idleStates, ok = step(states, idleStates, caves)
		if !ok {
			break
		}
	}

	return idleStates
}

func solveP1(input string, lines []string) string {
	runP1 := true
	if !runP1 {
		return "Skipping"
	}
	caves := parseInput(lines)
	findAllMoveCost(caves)
	states := []state_t{initState(caves)}
	states = stepUntilEnd(states, caves)
	sort.Slice(states, func(i, j int) bool {
		return states[i].currentFlow > states[j].currentFlow
	})
	return fmt.Sprintf("%v", states[0].currentFlow)
}

func findAgentActions(agentIdx int, state state2_t, caves caves_t) []agent_action_t {
	agent := state.agents[agentIdx]
	actions := []agent_action_t{}

	if agent.timeBusy > 0 {
		actions = append(actions, agent_action_t{option: Wait})
	} else {
		cave := caves[agent.position]
		for _, target := range maps.Keys(state.closedValves) {
			moveAndOpenCost := cave.moveCost[target] + 1
			if moveAndOpenCost > (state.similuationTime - state.minutesPassed) {
				continue
			}
			actions = append(actions, agent_action_t{option: MoveAndOpen, moveToLabel: target})
		}
		if len(actions) == 0 {
			actions = append(actions, agent_action_t{option: Idle})
		}
	}

	return actions
}

func copyAgents(agent agent_t, agentIdx int, agents []agent_t) []agent_t {
	newAgents := make([]agent_t, len(agents))
	copy(newAgents, agents)
	newAgents[agentIdx] = agent
	return newAgents
}

func doAgentAction(action agent_action_t, agentIdx int, state state2_t, caves caves_t) state2_t {
	if action.option == Idle {
		return state
	}
	if action.option == Wait {
		agent := state.agents[agentIdx]
		agent.timeBusy--
		return state2_t{
			agents:          copyAgents(agent, agentIdx, state.agents),
			minutesPassed:   state.minutesPassed,
			openedValves:    state.openedValves,
			flowSoFar:       state.flowSoFar,
			closedValves:    state.closedValves,
			similuationTime: state.similuationTime,
		}
	}
	if action.option == MoveAndOpen {
		targetPosition := action.moveToLabel
		stateWithAgentMoved := move2(targetPosition, agentIdx, state, caves)
		stateWithOpenValve := openValve2(agentIdx, stateWithAgentMoved, caves)
		return stateWithOpenValve
	}
	panic("Unknown option")
}

func step2(states []state2_t, idleStates []state2_t, caves caves_t) ([]state2_t, []state2_t, bool) {
	didChange := false
	newStates := []state2_t{}
	newIdleStates := []state2_t{}
	newIdleStates = append(newIdleStates, idleStates...)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	if len(states) > 100 {
		// Optimization to prune some states
		if states[0].minutesPassed > 0 {
			sort.Slice(states, func(i, j int) bool {
				return computeFlowRate(states[i], caves) > computeFlowRate(states[j], caves)
			})
			states = states[:int(float64(len(states))*0.95)]
		}
	}

	for _, state := range states {

		if state.minutesPassed >= state.similuationTime || len(state.closedValves) == 0 {
			newIdleStates = append(newIdleStates, state)
			continue
		}
		if r.Int() < (math.MaxInt / 100_000) {
			fmt.Println(state, computeFlowRate(state, caves))
		}

		agentActions := [2][]agent_action_t{}

		for agentIdx := range state.agents {
			actions := findAgentActions(agentIdx, state, caves)
			agentActions[agentIdx] = actions
		}

		for _, manAction := range agentActions[0] {
			for _, eleAction := range agentActions[1] {
				if manAction.option == MoveAndOpen &&
					eleAction.option == MoveAndOpen &&
					manAction.moveToLabel == eleAction.moveToLabel {
					// Optimization to avoid agents moving around (this may be one of the optimizations that broke the output)
					// Exercise passed without this, but before the agents had Move/Open actions separated,
					// so man and elephant could walk around aimlessly. It worked, but took a long time.
					// This is supposed to work, but not. Not sure why and will not optimize now.
					stateWithManOpening := doAgentAction(manAction, 0, state, caves)
					stateWithEleMoving := move2(eleAction.moveToLabel, 1, stateWithManOpening, caves)
					stateWithEleMoving.minutesPassed++
					newStates = append(newStates, stateWithEleMoving)

					stateWithEleOpening := doAgentAction(eleAction, 1, state, caves)
					stateWithManMoving := move2(manAction.moveToLabel, 0, stateWithEleOpening, caves)
					stateWithManMoving.minutesPassed++
					newStates = append(newStates, stateWithManMoving)
					didChange = true
					continue
				}
				statePrime := doAgentAction(manAction, 0, state, caves)
				statePrimePrime := doAgentAction(eleAction, 1, statePrime, caves)
				statePrimePrime.minutesPassed++
				newStates = append(newStates, statePrimePrime)
				didChange = true
			}
		}
	}

	return newStates, newIdleStates, didChange
}

func stepUntilEnd2(states []state2_t, caves caves_t) []state2_t {
	ok := false
	idleStates := []state2_t{}

	for {
		states, idleStates, ok = step2(states, idleStates, caves)
		if len(states) == 0 {
			fmt.Println(ok)
			break
		}
		// if !ok {
		// 	break
		// }
	}

	return idleStates
}

func computeFlowRate(state state2_t, caves caves_t) int {
	sum := 0
	for _, openedValve := range state.openedValves {
		cave := caves[openedValve.label]
		timeOpen := state.similuationTime - openedValve.minute
		sum += cave.flowRate * timeOpen
	}
	return sum
	// return state.flowSoFar
}

func solveP2(input string, lines []string) string {
	caves := parseInput(lines)
	findAllMoveCost(caves)
	initState := initState2(caves)
	initState.similuationTime = 26
	states := []state2_t{initState}
	states = stepUntilEnd2(states, caves)
	sort.Slice(states, func(i, j int) bool {
		return computeFlowRate(states[i], caves) > computeFlowRate(states[j], caves)
	})
	for _, s := range states[:10] {
		fmt.Println(s, computeFlowRate(s, caves))
	}
	return fmt.Sprintf("%v", computeFlowRate(states[0], caves))
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
