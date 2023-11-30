import java.io.File

data class Pos(val row: Int, val col: Int) {
    operator fun plus(other: Pos): Pos =
        Pos(row = row + other.row, col = col + other.col)
}

enum class Direction(val relativePos: Pos) {
    North(Pos(-1, 0)),
    NorthEast(Pos(-1, 1)),
    East(Pos(0, 1)),
    SouthEast(Pos(1, 1)),
    South(Pos(1, 0)),
    SouthWest(Pos(1, -1)),
    West(Pos(0, -1)),
    NorthWest(Pos(-1, -1));

    fun nextDirectionToCheck() =
        when (this) {
            North -> South
            South -> West
            West -> East
            East -> North
            else -> error("$this is not a checkable direction.")
        }

    fun mainDirectionsOrderFromThis() =
        when (this) {
            North -> listOf(North, South, West, East)
            South -> listOf(South, West, East, North)
            West -> listOf(West, East, North, South)
            East -> listOf(East, North, South, West)
            else -> error("$this is not a main direction!")
        }

    fun toCheck(): List<Direction> =
        when (this) {
            North -> listOf(NorthWest, North, NorthEast)
            East -> listOf(NorthEast, East, SouthEast)
            South -> listOf(SouthEast, South, SouthWest)
            West -> listOf(SouthWest, West, NorthWest)
            else -> error("Only main directions need to check, not $this")
        }
}

data class State(
    val round: Int = 0,
    val allPos: Set<Pos>,
    val firstDirectionToCheck: Direction = Direction.North,
    // A map with <Destination, list of Sources>
    val proposals: Map<Pos, List<Pos>> = emptyMap(),
)

const val inputFile = "input"
fun main(args: Array<String>) {
    part1()
    part2()
}

fun part1() {
    val initState = parseInput(inputFile)
    val stateAfter10 = roundNTimes(initState, 10)
//    draw(stateAfter10)
    val emptyGround = countEmptyTiles(stateAfter10, findBB(stateAfter10))
    println("Part1: $emptyGround")
}

fun part2() {
    val initState = parseInput(inputFile)
    val finalState = roundUntilNoMovement(initState)
//    draw(stateAfter10)
    println("Part2: ${finalState.round}")
}

fun roundNTimes(state: State, times: Int): State =
    (1..times).fold(state) { s, _ -> round(s) }

fun roundUntilNoMovement(state: State): State {
    var lastPos = setOf<Pos>()
    var currentState: State = state

    while (true) {
        currentState = round(currentState)
        if (currentState.allPos == lastPos) break
        lastPos = currentState.allPos.toSet()
    }

    return currentState
}

fun draw(state: State, width: Int = 14, height: Int = 12) {
    (0 until height).forEach { row ->
        (0 until width).forEach { col ->
            val hasElf = state.allPos.contains(Pos(row, col))
            val char = when (hasElf) {
                true -> '#'
                false -> '.'
            }
            print(char)
        }
        println()
    }
}

fun round(state: State): State =
    bumpStep(
        moveStep(
            proposeMovementsStep(state)
        )
    )

fun proposeMovementsStep(state: State): State {
    fun canMove(pos: Pos, direction: Direction): Boolean =
        direction.toCheck()
            .map { directionToCheck -> pos + directionToCheck.relativePos }
            .any { posToCheck -> state.allPos.contains(posToCheck) }
            .not()

    val proposals = mutableMapOf<Pos, List<Pos>>()

    state.allPos.forEach { pos ->
        val shouldMove = Direction.values()
            .map { it.relativePos }
            .any { neighbourRelativePos ->
                val posToLook = pos + neighbourRelativePos
                state.allPos.contains(posToLook)
            }
        if (!shouldMove) return@forEach

        val directionToMoveOrNull = state
            .firstDirectionToCheck
            .mainDirectionsOrderFromThis()
            .firstOrNull { canMove(pos, it) }
        if (directionToMoveOrNull != null) {
            val posToMoveTo = pos + directionToMoveOrNull.relativePos
            proposals[posToMoveTo] = proposals.getOrDefault(posToMoveTo, emptyList()) + pos
        }
    }
    return state.copy(proposals = proposals)
}

fun moveStep(state: State): State {
    val validProposals = state.proposals.filterValues { it.count() == 1 }
    val newAllPos = state.allPos.toMutableSet()

    validProposals.mapValues { it.value.first() }
        .entries
        .forEach { (dst, src) ->
            newAllPos.add(dst)
            newAllPos.remove(src)
        }

    return state.copy(allPos = newAllPos)
}

fun bumpStep(state: State): State =
    state.copy(
        round = state.round + 1,
        firstDirectionToCheck = state.firstDirectionToCheck.nextDirectionToCheck(),
        proposals = emptyMap(),
    )

fun findBB(state: State): Pair<Pos, Pos> {
    val pos = state.allPos
    val minRow = pos.minBy { it.row }.row
    val maxRow = pos.maxBy { it.row }.row
    val minCol = pos.minBy { it.col }.col
    val maxCol = pos.maxBy { it.col }.col

    return Pos(minRow, maxRow) to Pos(minCol, maxCol)
}

fun countEmptyTiles(state: State, bb: Pair<Pos, Pos>): Int {
    val bbWidth = bb.first.col - bb.first.row + 1
    val bbHeight = bb.second.col - bb.second.row + 1
    val bbSize = bbHeight * bbWidth

    return bbSize - state.allPos.count()
}

fun parseInputMap(lines: List<String>): Set<Pos> {
    val set = mutableSetOf<Pos>()
    lines.forEachIndexed { row, line ->
        line.forEachIndexed { col, char ->
            when (char) {
                '#' -> {
                    val pos = Pos(row, col)
                    set.add(pos)
                }
            }
        }
    }
    return set
}

fun parseInput(filename: String): State {
    val lines = File(filename).readLines()
    val allPos = parseInputMap(lines)

    return State(allPos = allPos)
}