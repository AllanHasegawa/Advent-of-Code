import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.async
import kotlinx.coroutines.runBlocking
import java.io.File
import kotlin.math.max

enum class Rock {
    Ore,
    Clay,
    Obsidian,
    Geode;

    companion object {
        val ALL_BUT_GEODE = listOf(Ore, Clay, Obsidian)
    }
}

typealias Cost = Map<Rock, Int>

fun Map<Rock, Int>.getOrZero(rock: Rock) = getOrDefault(rock, 0)
fun mapOfZeroes(): Map<Rock, Int> = Rock.values().associateWith { 0 }

data class Blueprint(
    val id: Int,
    val costs: Map<Rock, Cost>,
)

data class State(
    val blueprint: Blueprint,
    val totalSimulationMinutes: Int,
    val minute: Int = 0,
    val resources: Map<Rock, Int> = mapOfZeroes(),
    val robots: Map<Rock, Int> = Rock.values().associateWith { if (it == Rock.Ore) 1 else 0 },
    val robotsInProd: Map<Rock, Int> = mapOfZeroes(),
    val robotsAllowedProd: Map<Rock, Boolean> = Rock.values().associateWith { true },
) {
    val timeLeft: Int = totalSimulationMinutes - minute

    val canBuild = Rock.values().associateWith { canAffordRobot(it) }
}

val INPUT_REGEX =
    Regex("Blueprint (\\d+): Each ore robot costs (\\d+) ore. Each clay robot costs (\\d+) ore. Each obsidian robot costs (\\d+) ore and (\\d+) clay. Each geode robot costs (\\d+) ore and (\\d+) obsidian.")

fun main(@Suppress("UNUSED_PARAMETER") args: Array<String>) {
    val blueprints = parseInput(lines())

    // Part 1
    val part1 = runBlocking {
        val simulationTime = 24
        val jobs = blueprints.map { blueprint ->
            async(Dispatchers.IO) {
                qualityLevel(blueprint, simulationTime)
            }
        }
        val qualityLevels = jobs.map { it.await() }
        qualityLevels.sum()
    }
    println(part1)

    // Part 2
    val part2 = runBlocking {
        val simulationTime = 32
        val jobs = blueprints.take(3).map { blueprint ->
            async(Dispatchers.IO) {
                highestNumberOfGeode(blueprint, simulationTime)
            }
        }
        val geodes = jobs.map { it.await() }
        geodes.foldRight(1) { value, prod -> value * prod }
    }
    println(part2)

}

fun highestNumberOfGeode(blueprint: Blueprint, simulationTimeMinutes: Int): Int =
    highestNumberOfGeode(State(blueprint, simulationTimeMinutes), 0)

fun qualityLevel(blueprint: Blueprint, simulationTimeMinutes: Int): Int =
    highestNumberOfGeode(blueprint, simulationTimeMinutes) * blueprint.id

fun highestNumberOfGeode(state: State, maxGeodesSoFar: Int): Int {
    if (state.timeLeft == 0) {
        return state.resources[Rock.Geode]!!
    }
    @Suppress("NAME_SHADOWING")
    var state = deliveryPhase(state)
    if (state.estimatedGeodesUntilTheEnd() <= maxGeodesSoFar) return maxGeodesSoFar

    state = state.copy(minute = state.minute + 1)

    var statesPrime = buyPhase(state)
    statesPrime = statesPrime.map(::gatherPhase)

    val maxGeodes = max(maxGeodesSoFar, state.resources[Rock.Geode]!!)

    return statesPrime.maxOfOrNull { highestNumberOfGeode(it, maxGeodes) }
        ?: maxGeodes
}

fun deliveryPhase(state: State): State =
    state.copy(
        robots = Rock.values().associateWith { rock ->
            state.robots[rock]!! + state.robotsInProd[rock]!!
        },
        robotsInProd = mapOfZeroes(),
    )

fun buyPhase(state: State): List<State> {
    val afterGeodeRobotBuying = state.buyRobot(Rock.Geode)
    if (afterGeodeRobotBuying != null && afterGeodeRobotBuying.robotsInProd[Rock.Geode]!! > 0) {
        return listOf(afterGeodeRobotBuying)
    }

    val afterNotAffordingGeode =
        if (!state.canAffordRobot(Rock.Geode)) {
            state.copy(robotsAllowedProd = Rock.values().associateWith { !state.canAffordRobot(it) })
        } else null

    return Rock.ALL_BUT_GEODE.mapNotNull(state::buyRobot) + listOfNotNull(afterNotAffordingGeode)
}

fun gatherPhase(state: State): State =
    state.copy(
        resources = Rock.values().associateWith {
            (state.resources[it]!! + state.robots[it]!!)
        }
    )

fun State.estimatedGeodesUntilTheEnd(): Int {
    val geodesLeftWithCurrentRobots = timeLeft * robots[Rock.Geode]!!
    val geodesLeftWithFutureRobots =
        if (canBuild[Rock.Geode] == true) (timeLeft * (timeLeft - 1)) / 2
        else ((timeLeft - 1) * (timeLeft - 2)) / 2
    return resources[Rock.Geode]!! + geodesLeftWithCurrentRobots + geodesLeftWithFutureRobots
}

fun State.buyRobot(rock: Rock): State? {
    fun State.actuallyBuyIt(rock: Rock): State =
        (this - blueprint.costs[rock]!!).copy(
            robotsInProd = robotsInProd + (rock to 1),
            robotsAllowedProd = Rock.values().associateWith { true },
        )

    return if (canAffordRobot(rock)) {
        if (rock == Rock.Geode) actuallyBuyIt(rock)
        else if (robotsAllowedProd[rock] == true && needMoreForARobot(rock)) actuallyBuyIt(rock)
        else null
    } else null
}

fun State.needMoreForARobot(rock: Rock): Boolean =
    Rock.values().maxOf { blueprint.costs[it]!![rock]!! } > (robots[rock]!! + robotsInProd[rock]!!)

fun State.canAffordRobot(rock: Rock): Boolean {
    val cost = blueprint.costs[rock]!!
    return Rock.values().all { resources[it]!! >= cost[it]!! }
}

operator fun State.minus(cost: Cost): State = copy(
    resources = Rock.values().associateWith { resources[it]!! - cost[it]!! }
)

fun lines(filename: String = "input"): List<String> = File(filename).readLines()

fun parseInput(lines: List<String>): List<Blueprint> = lines.map(::parseLine)

fun parseLine(line: String): Blueprint {
    val match = INPUT_REGEX.matchEntire(line) ?: error("Invalid input: $line")
    var groupIdx = 1
    fun nextMatch(): Int = match.groups[groupIdx++]!!.value.toInt()

    return Blueprint(
        id = nextMatch(),
        costs = mapOf(
            Rock.Ore to mapOf(Rock.Ore to nextMatch()).fillCosts(),
            Rock.Clay to mapOf(Rock.Ore to nextMatch()).fillCosts(),
            Rock.Obsidian to mapOf(Rock.Ore to nextMatch(), Rock.Clay to nextMatch()).fillCosts(),
            Rock.Geode to mapOf(Rock.Ore to nextMatch(), Rock.Obsidian to nextMatch()).fillCosts(),
        )
    )
}

fun Map<Rock, Int>.fillCosts() =
    Rock.values().associateWith { rock -> this[rock] ?: 0 }