import java.io.File

data class Pos(val row: Int, val col: Int) {
    fun reverse(): Pos =
        when (this) {
            RIGHT -> LEFT
            LEFT -> RIGHT
            UP -> DOWN
            DOWN -> UP
            else -> error("Pos $this can't be reversed.")
        }

    operator fun minus(other: Pos): Pos =
        Pos(row = row - other.row, col = col - other.col)

    operator fun plus(other: Pos) =
        Pos(row = row + other.row, col = col + other.col)

    companion object {
        val RIGHT = Pos(0, 1)
        val DOWN = Pos(1, 0)
        val LEFT = Pos(0, -1)
        val UP = Pos(-1, 0)
    }
}

sealed class Tile {
    abstract val pos: Pos

    data class Wall(override val pos: Pos) : Tile()
    data class Open(override val pos: Pos) : Tile()
    data class Portals(override val pos: Pos, val goTos: Map<Pos, Landing>) : Tile()
}

data class Landing(
    val pos: Pos,
    val facingChange: Pos? = null,
)

sealed class Direction {
    object TurnLeft : Direction() {
        override fun toString(): String = "Turn Left"
    }

    object TurnRight : Direction() {
        override fun toString(): String = "Turn Right"
    }

    data class Forward(val steps: Int) : Direction() {
        override fun toString(): String = "Forward ($steps)"
    }
}

data class State(
    val pos: Pos,
    val facing: Pos = Pos(0, 1),
    val map: Map<Pos, Tile>,
    val directions: List<Direction>,
)

const val INPUT = "input"
fun main(args: Array<String>) {
    println("input: ${lines()}")

    part1()
    part2()
}

fun part1() {
    val (map, directions) = parseInput()
    val startPos = findStartPos(map)
    val initialState = State(
        pos = startPos,
        map = map,
        directions = directions,
    )
    val finalState = initialState.moveAll()

    initialState.debug(10)

    val answer = 1000 * (finalState.pos.row + 1) +
            4 * (finalState.pos.col + 1) +
            finalState.facing.toAnswerValue()
    println("part1: $answer")
}

fun part2() {
    var (map, directions) = parseInput()
    map = makeMapACube(map, true)
    val startPos = findStartPos(map)
    val initialState = State(
        pos = startPos,
        map = map,
        directions = directions,
    )
//    drawMap(initialState)
    val finalState = initialState.moveAll()

//    initialState.debug(10)

    val answer = 1000 * (finalState.pos.row + 1) +
            4 * (finalState.pos.col + 1) +
            finalState.facing.toAnswerValue()
    println("part2: $answer")
}

fun makeMapACube(map: Map<Pos, Tile>, fullInput: Boolean): Map<Pos, Tile> {
    require(fullInput)

    val m = map.toMutableMap()
    // Clear all the portals destinations
    m.mapValuesTo(m) { (_, value) ->
        if (value is Tile.Portals) value.copy(goTos = emptyMap())
        else value
    }

    fun posByFace(rowIdx: Int, colIdx: Int) = Pos(rowIdx * 50, colIdx * 50)

    fun addPortal(pos: Pos, facing: Pos, goTo: Pos, landingFacing: Pos) {
        val ogPortal = (m[pos] as? Tile.Portals) ?: error("Tile at $pos is not Portal, it's ${m[pos]}")
        val portal = ogPortal.copy(goTos = ogPortal.goTos + (facing to Landing(goTo, landingFacing)))
        m[pos] = portal
    }

    fun addPortalLine(
        startPos: Pos,
        startPosInc: Pos,
        startFacing: Pos,
        endPos: Pos,
        endPosInc: Pos,
        endFacing: Pos,
    ) {
        var startPosGuide = startPos
        var endPosGuide = endPos
        repeat(50) {
            addPortal(startPosGuide, startFacing, endPosGuide, endFacing)
            addPortal(endPosGuide, endFacing.reverse(), startPosGuide, startFacing.reverse())
            startPosGuide += startPosInc
            endPosGuide += endPosInc
        }
    }

    /**
     * 12
     * 3
     *45
     *6
     */

    // 1u6l
    addPortalLine(posByFace(0, 1) + Pos.UP, Pos.RIGHT, Pos.UP, posByFace(3, 0) + Pos.LEFT, Pos.DOWN, Pos.RIGHT)
    // 1l4l
    addPortalLine(
        posByFace(0, 1) + Pos.LEFT,
        Pos.DOWN,
        Pos.LEFT,
        posByFace(3, 0) + Pos.LEFT + Pos.UP,
        Pos.UP,
        Pos.RIGHT,
    )
    // 2u6d
    addPortalLine(posByFace(0, 2) + Pos.UP, Pos.RIGHT, Pos.UP, posByFace(4, 0), Pos.RIGHT, Pos.UP)
    // 2d3r
    addPortalLine(posByFace(1, 2), Pos.RIGHT, Pos.DOWN, posByFace(1, 2), Pos.DOWN, Pos.LEFT)
    // 3l4u
    addPortalLine(posByFace(1, 1) + Pos.LEFT, Pos.DOWN, Pos.LEFT, posByFace(2, 0) + Pos.UP, Pos.RIGHT, Pos.DOWN)
    // 5d6r
    addPortalLine(posByFace(3, 1), Pos.RIGHT, Pos.DOWN, posByFace(3, 1), Pos.DOWN, Pos.LEFT)
    // 5r2r
    addPortalLine(posByFace(2, 2), Pos.DOWN, Pos.RIGHT, posByFace(1, 3) + Pos.UP, Pos.UP, Pos.LEFT)

    return m
}

fun State.debug(moves: Int) {
    var s = this
    for (c in 1..moves) {
        println("Moving: ${s.directions.first()}")
        s = s.moveOnce()
        drawMap(s)
    }
}

fun State.moveAll(): State {
    var state = this
    while (state.directions.isNotEmpty()) {
        state = state.moveOnce()
    }
    return state
}

fun State.moveOnce(): State {
    require(directions.isNotEmpty())

    return when (val direction = directions.first()) {
        Direction.TurnLeft,
        Direction.TurnRight -> {
            copy(
                facing = facing.turn(direction),
                directions = directions.drop(1),
            )
        }

        is Direction.Forward -> {
            var steps = direction.steps
            var pos = pos
            var facing = facing

            while (steps > 0) {
                val lookAhead = pos + facing
                when (val tile = map[lookAhead]) {
                    is Tile.Open -> pos = lookAhead
                    is Tile.Wall -> break
                    is Tile.Portals -> {
                        val landing = tile.goTos[facing] ?: error("Portal $tile doesn't have goto for $facing")
                        if (map[landing.pos] is Tile.Portals) {
                            val lookAfterPortal = landing.pos + (landing.facingChange ?: facing)

                            when (map[lookAfterPortal]) {
                                is Tile.Open -> {
                                    if (landing.facingChange != null) facing = landing.facingChange
                                    pos = landing.pos + facing
                                }

                                is Tile.Wall -> break
                                is Tile.Portals -> error("Portal after portal not allowed.")
                                null -> error("No tile after portal.")
                            }
                        } else {
                            error("Tile ${map[landing.pos]} is not a portal.}")
                        }
                    }

                    else -> error("Unknown tile at $pos")
                }
                steps--
            }
            copy(
                pos = pos,
                facing = facing,
                directions = directions.drop(1),
            )
        }
    }
}

fun Pos.toAnswerValue(): Int =
    when (this) {
        Pos.RIGHT -> 0
        Pos.DOWN -> 1
        Pos.LEFT -> 2
        Pos.UP -> 3
        else -> error("Not a facing value $this")
    }

fun Pos.turn(direction: Direction): Pos = when (direction) {
    is Direction.TurnRight -> when (this) {
        Pos.RIGHT -> Pos.DOWN
        Pos.DOWN -> Pos.LEFT
        Pos.LEFT -> Pos.UP
        Pos.UP -> Pos.RIGHT
        else -> error("Cant turn ($this) right")
    }

    is Direction.TurnLeft -> turn(Direction.TurnRight).turn(Direction.TurnRight).turn(Direction.TurnRight)
    else -> error("Cant turn a ($direction)")
}

fun findStartPos(map: Map<Pos, Tile>): Pos {
    val topRow = map.filter { it.value !is Tile.Portals }.minBy { it.key.row }.key.row

    val leftMostColumn =
        map.filter { it.value !is Tile.Portals }.filter { it.key.row == topRow }.minBy { it.key.col }.key.col

    return Pos(topRow, leftMostColumn)
}

fun drawMap(state: State, height: Int = 15, width: Int = 100) {
    val map = state.map
    val array = Array(height) { Array(width) { ' ' } }

    map.values.forEach { tile ->
        if (tile.pos.row < height - 1 && tile.pos.col < width - 1) {
            array[tile.pos.row + 1][tile.pos.col + 1] = when (tile) {
                is Tile.Wall -> '#'
                is Tile.Portals -> '@'
                is Tile.Open -> '.'
            }
        }
    }

    val character = when (state.facing) {
        Pos.RIGHT -> '>'
        Pos.DOWN -> 'V'
        Pos.LEFT -> '<'
        Pos.UP -> '^'
        else -> error("Unknown facing ${state.facing}")
    }
    val pos = state.pos
    if (pos.row < height - 1 && pos.col < width - 1) {
        array[pos.row + 1][pos.col + 1] = character
    }

    for (y in 0 until height) {
        for (x in 0 until width) {
            print(array[y][x])
        }
        println()
    }
}

private fun lines(filename: String = INPUT): List<String> = File(filename).readLines()

private fun parseInput(): Pair<Map<Pos, Tile>, List<Direction>> {
    val map = mutableMapOf<Pos, Tile>()
    val directions = mutableListOf<Direction>()

    lines().forEachIndexed { idx, line ->
        if (line.contains(".")) {
            val mapLine = parseMapLine(line, idx)
            map.putAll(mapLine.map { it.pos to it })
        }

        if (line.contains("R")) {
            directions.addAll(parseDirectionsLine(line))
        }
    }
    addTopBottomPortals(map)

    return map to directions
}

fun addTopBottomPortals(map: MutableMap<Pos, Tile>) {
    var column = 0

    while (true) {
        val columnEntries = map.filterKeys { it.col == column }
        val topRow = columnEntries.filterValues { it !is Tile.Portals }.minByOrNull { it.key.row }?.key ?: break
        val bottomRow = columnEntries.filterValues { it !is Tile.Portals }.maxByOrNull { it.key.row }?.key ?: break

        val topPortalPos = topRow.copy(row = topRow.row - 1)
        val bottomPortalPos = bottomRow.copy(row = bottomRow.row + 1)

        val topPortal = (map[topPortalPos] as? Tile.Portals) ?: Tile.Portals(topPortalPos, mapOf())
        map[topPortalPos] = topPortal.copy(
            goTos = topPortal.goTos + (Pos.UP to Landing(bottomPortalPos))
        )
        val bottomPortal = (map[bottomPortalPos] as? Tile.Portals) ?: Tile.Portals(bottomPortalPos, mapOf())
        map[bottomPortalPos] = bottomPortal.copy(
            goTos = bottomPortal.goTos + (Pos.DOWN to Landing(topPortalPos))
        )

        column++
    }
}

private fun parseMapLine(line: String, row: Int): List<Tile> {
    val tiles = mutableListOf<Tile>()
    for (i in line.indices) {
        val char = line[i]
        val pos = Pos(row, i)
        val tileOrNull = when (char) {
            '.' -> Tile.Open(pos)
            '#' -> Tile.Wall(pos)
            ' ' -> null
            else -> error("Unknown tile at $pos, $char")
        }
        if (tileOrNull != null) {
            tiles.add(tileOrNull)
        }
    }

    val startIndex = line.indexOfFirst { it != ' ' }
    val endIndex = line.length - 1
    val startPortalCol = startIndex - 1
    val endPortalCol = endIndex + 1
    val startPortal = Tile.Portals(Pos(row, startPortalCol), mapOf(Pos.LEFT to Landing(Pos(row, endPortalCol))))
    val endPortal = Tile.Portals(Pos(row, endPortalCol), mapOf(Pos.RIGHT to Landing(Pos(row, startPortalCol))))
    tiles.add(startPortal)
    tiles.add(endPortal)
    return tiles
}

fun parseDirectionsLine(line: String): List<Direction> {
    val lines = line.replace("R", "\nR\n").replace("L", "\nL\n").lines()

    return lines.map { l ->
        when (l) {
            "R" -> Direction.TurnRight
            "L" -> Direction.TurnLeft
            else -> Direction.Forward(l.toInt())
        }
    }
}