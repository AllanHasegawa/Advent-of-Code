import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.async
import kotlinx.coroutines.awaitAll
import kotlinx.coroutines.runBlocking
import java.io.File
import java.math.BigInteger
import kotlin.math.abs

const val INPUT_REGEX_PATTERN = """([a-z]{4}): ((\d+)|(([a-z]{4}) ([/*+-]) ([a-z]{4})))"""
val INPUT_REGEX = Regex(INPUT_REGEX_PATTERN)
const val HUMN = "humn"
const val ROOT = "root"

enum class OpType(val visual: String) {
    Plus("+"), Minus("-"), Multiply("*"), Div("/");

    override fun toString(): String = visual
}

data class Line(
    val assigned: Expr.Ref,
    val rightHand: Expr,
) {
    override fun toString(): String = "${assigned.value}: $rightHand"
}

sealed interface Expr {
    data class Literal(val value: BigInteger) : Expr {
        override fun toString(): String = value.toString()
    }

    data class Ref(val value: String) : Expr {
        override fun toString(): String = value
    }

    data class Op(val left: Expr, val op: OpType, val right: Expr) : Expr {
        override fun toString(): String = "$left $op $right"
    }

    data class Humn(val humn: (BigInteger) -> BigInteger) : Expr {
        override fun toString(): String = "Humn"
    }
}

const val INPUT_FILE = "input"
fun main(args: Array<String>) {
    println("input: ${parseInput()}")
    val stack = parseInput().associate { it.assigned to it.rightHand }

    part1(stack)
    part2(stack)
}


private fun part1(stack: Map<Expr.Ref, Expr>) {
    val humnDefaultValue = (stack[Expr.Ref(HUMN)] as Expr.Literal).value
    val rootExprSolved = solve(stack, Expr.Ref(ROOT)) as Expr.Humn
    println("Part1: " + rootExprSolved.humn(humnDefaultValue))
}

private fun part2(stack: Map<Expr.Ref, Expr>) {
    val rootExpr = stack[Expr.Ref(ROOT)] as Expr.Op
    val rootLeftExpr = rootExpr.left as Expr.Ref
    val rootRightExpr = rootExpr.right as Expr.Ref

    // Both inputs have "humn" on the left and root is a plus op.
    val leftSolved = solve(stack, rootLeftExpr) as Expr.Humn
    val rightSolved = solve(stack, rootRightExpr) as Expr.Literal
    require(rootExpr.op == OpType.Plus)

    val humn = leftSolved.humn
    val right = rightSolved.value
    val objectiveF = { it: BigInteger -> humn(it) - right }

    /**
     * I plotted the result of the objective function to understand the shape of the function.
     * It's a linear, decrementing, function with it's zero value close to the center on a -Long to Long range.
     */
//    part2ToCsv(right, humn)

    val (solution, _) = binarySearch(objectiveF, -10_000_000_000_000L, +10_000_000_000_000L)

    println("Part2: $solution")
}

fun solve(stack: Map<Expr.Ref, Expr>, ref: Expr.Ref): Expr {
    if (ref.value == HUMN) return Expr.Humn { it }

    return when (val expr = stack[ref]!!) {
        is Expr.Literal -> expr
        is Expr.Humn -> expr
        is Expr.Op -> {
            val left = solve(stack, expr.left as Expr.Ref)
            val right = solve(stack, expr.right as Expr.Ref)
            doOp(left, expr.op, right)
        }

        is Expr.Ref -> error("Right hand side is never just a ref")
    }
}

fun doOp(left: Expr, op: OpType, right: Expr): Expr =
    when (left) {
        is Expr.Literal -> {
            when (right) {
                is Expr.Literal -> {
                    when (op) {
                        OpType.Plus -> left + right
                        OpType.Minus -> left - right
                        OpType.Multiply -> left * right
                        OpType.Div -> left / right
                    }
                }

                is Expr.Humn -> {
                    Expr.Humn {
                        when (op) {
                            OpType.Plus -> left.value + right.humn(it)
                            OpType.Minus -> left.value - right.humn(it)
                            OpType.Multiply -> left.value * right.humn(it)
                            OpType.Div -> left.value / right.humn(it)
                        }
                    }
                }

                else -> error("Right can't do op: $right")
            }
        }

        is Expr.Humn -> {
            when (right) {
                is Expr.Literal -> {
                    Expr.Humn {
                        when (op) {
                            OpType.Plus -> left.humn(it) + right.value
                            OpType.Minus -> left.humn(it) - right.value
                            OpType.Multiply -> left.humn(it) * right.value
                            OpType.Div -> left.humn(it) / right.value
                        }
                    }
                }

                is Expr.Humn -> {
                    Expr.Humn {
                        val leftHumn = left.humn(it)
                        val rightHumn = right.humn(it)
                        when (op) {
                            OpType.Plus -> leftHumn + rightHumn
                            OpType.Minus -> leftHumn - rightHumn
                            OpType.Multiply -> leftHumn * rightHumn
                            OpType.Div -> leftHumn / rightHumn
                        }
                    }
                }

                else -> error("Right can't do op: $right")
            }
        }

        else -> error("Left can't do op: $left")
    }

/**
 * Took to long :)
 */
fun searchInterval(f: (Long) -> Long, start: Long, end: Long): Pair<Long, Long> {
    val threads = 12
    val step = abs(end - start) / threads
    var solution = 0L

    runBlocking(Dispatchers.Default) {
        val jobs = (0 until threads).map { threadIdx ->
            async {

                var i = start + threadIdx.toLong()

                while (true) {
                    val result = f(i)

                    if (result == 0L) {
                        solution = i
                    }

                    i += threads
                    if (solution != 0L) break
                    if (i > end) break
                }
            }
        }
        jobs.awaitAll()
    }

    return solution to f(solution)
}

/**
 * Slightly hacked binary search function. It searches for values in a DECREMENTING function :)
 * Too lazy to adjust for all cases. It also means the test input would not work xD
 */
fun binarySearch(f: (BigInteger) -> BigInteger, start: Long, end: Long): Pair<Long, BigInteger> {
    var solution = 0L
    var low = start
    var high = end

    while (low <= high) {
        val mid = low + ((high - low) / 2)
        val result = f(mid.toBigInteger())
        if (result > BigInteger.ZERO) {
            low = mid + 1;
        } else if (result < BigInteger.ZERO) {
            high = mid - 1
        } else {
            solution = mid
            break
        }
    }
    return solution to f(solution.toBigInteger())
}

fun Expr.Literal.doOp(other: Expr.Literal, op: (BigInteger, BigInteger) -> BigInteger): Expr.Literal =
    Expr.Literal(op(value, other.value))

operator fun Expr.Literal.plus(other: Expr.Literal): Expr.Literal = doOp(other) { a, b -> a + b }
operator fun Expr.Literal.minus(other: Expr.Literal): Expr.Literal = doOp(other) { a, b -> a - b }
operator fun Expr.Literal.times(other: Expr.Literal): Expr.Literal = doOp(other) { a, b -> a * b }
operator fun Expr.Literal.div(other: Expr.Literal): Expr.Literal = doOp(other) { a, b -> a / b }

private fun lines(filename: String = INPUT_FILE): List<String> =
    File(filename).readLines()

private fun parseInput(): List<Line> =
    lines().map(::parseLine)

private fun parseLine(line: String): Line {
    val match = INPUT_REGEX.matchEntire(line)
    requireNotNull(match) { "$line did not match" }

    val source = Expr.Ref(match.groups[1]!!.value)

    val rightHand =
        if (match.groups[7] == null) {
            Expr.Literal(match.groups[3]!!.value.toLong().toBigInteger())
        } else {
            val opRaw = match.groups[6]!!.value
            val op = OpType.values().first { it.visual == opRaw }
            Expr.Op(
                left = Expr.Ref(match.groups[5]!!.value),
                op = op,
                right = Expr.Ref(match.groups[7]!!.value),
            )
        }

    return Line(source, rightHand)
}

private fun part2ToCsv(right: BigInteger, humn: (BigInteger) -> BigInteger) {
    val samples = 240_000L
    val range = Long.MAX_VALUE / 32
    val step = (range * 2) / samples
    val writer = File("out.csv").bufferedWriter()
    writer.appendLine("-1,$right,0")
    for (i in -range..range step step) {
        val bi = i.toBigInteger()
        writer.appendLine("$i,${humn(bi)},${humn(bi) - right}")
    }
    writer.flush()
    writer.close()
}
