import java.io.File
import kotlin.math.abs
import kotlin.math.max
import kotlin.math.min

data class Pos(val row: Int, val col: Int) {
    operator fun plus(other: Pos) =
        Pos(row = row + other.row, col = col + other.col)

    fun distance(other: Pos): Int =
        abs(row - other.row) + abs(col - other.col)
}

enum class Direction(val relativePos: Pos) {
    Up(Pos(-1, 0)),
    Down(Pos(1, 0)),
    Left(Pos(0, -1)),
    Right(Pos(0, 1));
}

data class Blizzard(val direction: Direction) {
    fun toChar() =
        when (direction) {
            Direction.Left -> '<'
            Direction.Right -> '>'
            Direction.Up -> '^'
            Direction.Down -> 'v'
        }
}

sealed interface Tile {
    data class BlizzardTile(val blizzards: List<Blizzard>) : Tile {
        constructor(single: Blizzard) : this(listOf(single))

        operator fun plus(other: Blizzard) =
            copy(blizzards = blizzards + other)
    }

    object Wall : Tile
}

data class State(
    val start: Pos = Pos(0, 1),
    val end: Pos,
    val currentPos: Pos = Pos(0, 1),
    val steps: Int = 0,
    val stepsTook: List<Pos> = listOf(Pos(0, 1)),
) {
    override fun equals(other: Any?): Boolean {
        if (this === other) return true
        if (javaClass != other?.javaClass) return false

        other as State

        if (currentPos != other.currentPos) return false
        if (steps != other.steps) return false

        return true
    }

    override fun hashCode(): Int {
        var result = currentPos.hashCode()
        result = 31 * result + steps
        return result
    }
}

const val input = "input"

// Start history with input, step 0
val mapHistory = mutableListOf(parseMap(input))
var mapHeight: Int = 0
var mapWidth: Int = 0

fun main(args: Array<String>) {
    part1()
    part2()
}

fun part1() {
    val initState = State(end = Pos(mapHeight, mapWidth - 1))
    val endState = findExit(initState)
    println("Part 1: ${endState.steps}")
}

fun part2() {
    val startLine = Pos(0, 1)
    val finishLine = Pos(mapHeight, mapWidth - 1)
    val initState = State(end = finishLine)
    val endState0 = findExit(initState) // at finish line
    val endState1 = findExit(endState0.copy(end = startLine)) // at start line
    val endState2 = findExit(endState1.copy(end = finishLine)) // and finish again xD
    println("Part 2: ${endState2.steps}")
}

fun findExit(initState: State): State {
    val todo = mutableSetOf<State>()
    todo.add(initState)
    var maxStepFound: Int = Int.MAX_VALUE
    lateinit var stateAtEnd: State

    while (todo.isNotEmpty()) {
        val state = todo.first()
        todo.remove(state)

        if (state.currentPos == state.end) {
            maxStepFound = min(state.steps, maxStepFound)
            if (state.steps <= maxStepFound) {
                stateAtEnd = state
            }
        }
        if (state.steps + state.currentPos.distance(state.end) >= maxStepFound) {
            continue
        }

        val nextMap = getMapInHistory(state.steps + 1)
        val availablePos = Direction.values()
            .map { state.currentPos + it.relativePos }
            .filter { nextMap[it] == null && it.row >= 0 && it.col >= 0 }


        val newStates = availablePos
            .map { state.copy(currentPos = it, steps = state.steps + 1) }
            .let {
                val canWait = nextMap[state.currentPos] == null
                when (canWait) {
                    true -> it + state.copy(steps = state.steps + 1)
                    false -> it
                }
            }

        todo.addAll(newStates)
    }
    return stateAtEnd
}

fun drawUpToEnd(state: State) {
    (0 until state.steps).forEach { step ->
        println("Step $step (${state.stepsTook[step]}):")
        draw(step, state.stepsTook[step])
    }
}

fun draw(step: Int, pos: Pos) {
    val map = getMapInHistory(step)

    (0..mapHeight).forEach { row ->
        (0..mapWidth).forEach { col ->
            val currPos = Pos(row, col)

            val char =
                when {
                    currPos == pos -> 'E'
                    map[currPos] is Tile.Wall -> '#'
                    map[currPos] is Tile.BlizzardTile -> {
                        val tile = map[currPos] as Tile.BlizzardTile
                        when (val count = tile.blizzards.count()) {
                            in 2..Int.MAX_VALUE -> count.toString()[0]
                            1 -> tile.blizzards.first().toChar()
                            0 -> error("Empty blizzard tile: $tile")
                            else -> error("Unknown count: $count, for $tile")
                        }
                    }

                    else -> '.'
                }
            print(char)
        }
        println()
    }
}

fun getMapInHistory(step: Int): Map<Pos, Tile> {
    if (mapHistory.count() > step) return mapHistory[step]

    val oldMap = getMapInHistory(step - 1)
    val newMap = mutableMapOf<Pos, Tile>()

    fun moveBlizzard(pos: Pos, blizzard: Blizzard, map: MutableMap<Pos, Tile>) {
        var newPos = pos + blizzard.direction.relativePos

        when {
            newPos.row == 0 ->
                newPos = newPos.copy(row = mapHeight - 1)

            newPos.row == mapHeight ->
                newPos = newPos.copy(row = 1)

            newPos.col == 0 ->
                newPos = newPos.copy(col = mapWidth - 1)

            newPos.col == mapWidth ->
                newPos = newPos.copy(col = 1)
        }

        when (val tile = map[newPos]) {
            null -> map[newPos] = Tile.BlizzardTile(blizzard)
            is Tile.BlizzardTile -> map[newPos] = tile + blizzard
            is Tile.Wall -> error("This should have been handled above!")
        }
    }

    oldMap.forEach { (pos, tile) ->
        if (tile is Tile.Wall) newMap[pos] = tile
    }

    oldMap.forEach { (pos, tile) ->
        if (tile is Tile.Wall) return@forEach
        tile as Tile.BlizzardTile
        tile.blizzards.forEach { moveBlizzard(pos, it, newMap) }
    }

    mapHistory.add(newMap)

    return newMap
}

fun parseLine(line: String, row: Int): Map<Pos, Tile> {
    val map = mutableMapOf<Pos, Tile>()
    line.forEachIndexed { index, c ->
        mapWidth = max(index, mapWidth)
        val pos = Pos(row = row, col = index)
        when (c) {
            '.' -> Unit
            '#' -> map[pos] = Tile.Wall
            '^' -> map[pos] = Tile.BlizzardTile(Blizzard(Direction.Up))
            '>' -> map[pos] = Tile.BlizzardTile(Blizzard(Direction.Right))
            'v' -> map[pos] = Tile.BlizzardTile(Blizzard(Direction.Down))
            '<' -> map[pos] = Tile.BlizzardTile(Blizzard(Direction.Left))
            else -> error("Unknown char $c")
        }
    }
    return map
}

fun parseMap(filename: String): Map<Pos, Tile> {
    val map = mutableMapOf<Pos, Tile>()

    File(filename).readLines().forEachIndexed { row, line ->
        mapHeight = row
        map += parseLine(line, row)
    }

    return map
}