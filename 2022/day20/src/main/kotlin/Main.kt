import LinkedList.Node
import java.io.File

const val DECRIPTION_KEY = 811589153

fun main(args: Array<String>) {
    println("input: ${parseInput()}")

    part1()
    part2()
}

fun part1() {
    val list = LinkedList()
    parseInput().forEach(list::add)

    mixWithOgOrder(list)

    val coords = grooveCoordinates(list)
    println(coords)
    val part1 = coords.sum()
    println("Part 1: $part1")
}

fun part2() {
    val list = LinkedList()
    parseInput().map { it * DECRIPTION_KEY }.forEach(list::add)

    for (i in 1..10) {
        mixWithOgOrder(list)
    }

    val coords = grooveCoordinates(list)
    println(coords)
    val part2 = coords.sum()
    println("Part 2: $part2")
}

fun grooveCoordinates(list: LinkedList): List<Long> {
    var nodeWithZero: Node? = null

    var node = list.currentRoot
    while (node != null) {
        if (node.value == 0L) {
            nodeWithZero = node
            break
        }
        node = node.currentNext
    }
    require(nodeWithZero != null)

    fun forward(list: LinkedList, node: Node, amount: Int): Node {
        if (amount == 0) return node
        val next = node.currentNext ?: list.currentRoot
        return forward(list, next!!, amount - 1)
    }

    val first = forward(list, nodeWithZero, 1000)
    val second = forward(list, first, 1000)
    val third = forward(list, second, 1000)

    return listOf(first, second, third).map(Node::value)
}

fun mixWithOgOrder(list: LinkedList) {
    require(list.currentRoot != null)
    require(list.size > 1)

    var ogNode = list.ogRoot
    while (ogNode != null) {
        list.mixin(ogNode)
        ogNode = ogNode.originalNext
    }
}

private fun lines(filename: String = "input"): List<String> =
    File(filename).readLines()

fun parseInput(): List<Long> =
    lines().map { it.toLong() }

class LinkedList {
    data class Node(
        var originalNext: Node?,
        var originalPrevious: Node?,
        var currentNext: Node?,
        var currentPrevious: Node?,
        val value: Long,
    ) {
        override fun toString(): String =
            "$value [${currentPrevious?.value},${currentNext?.value}]"
    }

    var size: Int = 0
    var currentRoot: Node? = null
    var currentTail: Node? = null
    var ogRoot: Node? = null
    var ogTail: Node? = null

    fun add(value: Long): Node {
        size++

        val node = Node(
            originalNext = null,
            originalPrevious = currentTail,
            currentNext = null,
            currentPrevious = currentTail,
            value = value,
        )
        currentTail?.originalNext = node
        currentTail?.currentNext = node
        currentTail = node
        ogTail = node

        if (currentRoot == null) currentRoot = node
        if (ogRoot == null) ogRoot = node

        return node
    }

    fun mixin(node: Node) {
        require(currentRoot != null)
        require(currentTail != null)

        if (node.value == 0L) return

        fun walk(node: Node, move: Int): Node {
            if (move == 0) return node
            var direction =
                if (move > 0) node.currentNext
                else node.currentPrevious
            if (direction == null) {
                direction = if (move > 0) currentRoot else currentTail
            }
            return walk(direction!!, move - sign(move))
        }

        fun remove(node: Node) {
            val previous = node.currentPrevious
            val next = node.currentNext

            if (previous == null) {
                currentRoot = next!!
            }
            if (next == null) {
                currentTail = previous!!
            }

            previous?.currentNext = next
            next?.currentPrevious = previous
        }

        fun insert(node: Node, dest: Node) {
            val nextNode = dest.currentNext
            if (nextNode == null) {
                currentTail = node
            }

            dest.currentNext = node
            nextNode?.currentPrevious = node

            node.currentPrevious = dest
            node.currentNext = nextNode
        }

        var move = (node.value % (size - 1)).toInt()// + (node.value / size)
        if (move == 0) {
            if (node.value > 0) move++
            else move--
        }
        if (node.value < 0) move--
        if (move == 0) return

        remove(node)

        val dest = walk(node, move)
        insert(node, dest)
    }

    override fun toString(): String {
        fun recursiveToString(node: Node?, nextF: (Node) -> Node?, builder: StringBuilder) {
            when (node) {
                null -> builder.append("null")
                else -> {
                    builder.append(node.value)
                    builder.append(", ")
                    recursiveToString(nextF(node), nextF, builder)
                }
            }
        }

        val builder = StringBuilder()
        builder.append("Size: ")
        builder.append(size)
        builder.appendLine()
        builder.append("Og: ")
        recursiveToString(ogRoot, { it.originalNext }, builder)
        builder.appendLine()
        builder.append("Cr: ")
        recursiveToString(currentRoot, { it.currentNext }, builder)
        builder.appendLine()

        return builder.toString()
    }
}

fun sign(v: Int) = if (v > 0) 1 else -1