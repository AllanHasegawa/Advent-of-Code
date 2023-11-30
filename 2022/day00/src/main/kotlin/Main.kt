import java.io.File

fun main(args: Array<String>) {
    println("input: ${lines()}")
}

private fun lines(filename: String = "input"): List<String> =
    File(filename).readLines()