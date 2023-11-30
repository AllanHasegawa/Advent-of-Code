import java.io.File
import kotlin.math.pow

const val DIGITS = "=-012"

const val input = "input"
fun main(args: Array<String>) {
    part1()
}

fun part1() {
    val sum = lines().map(::snafuToDecimal).sum()
    println("Part 1 (sum): $sum")
    print("Part 1: ")
    println(decimalToSnafu(sum))
}

fun decimalToSnafu(decimal: Long): String =
    generateSequence(decimal) { (it + 2) / 5 }
        .takeWhile { it != 0L }
        .map { (it + 2) % 5 }
        .map { DIGITS[it.toInt()] }
        .joinToString("")
        .reversed()

fun snafuDigitToDecimal(digit: Char): Long =
    (-2L) + DIGITS.indexOf(digit)

fun snafuToDecimal(snafu: String): Long =
    snafu.reversed().foldRightIndexed(0L) { idx, snafuDigit, acc ->
        val digitDecValue = snafuDigitToDecimal(snafuDigit)
        val decValue = digitDecValue * (5.0.pow(idx).toLong())
        acc + decValue
    }

private fun lines(filename: String = input): List<String> =
    File(filename).readLines()