use std::fs;

fn main() {
    let test = false;
    let filepath = String::from(if test { "input_test" } else { "input" });
    let input = raw_input(filepath);

    part1(input.clone());
    part2(input.clone());
}

fn part1(input: String) {
}

fn part2(input: String) {
}


fn raw_input(filename: String) -> String {
    fs::read_to_string(filename)
        .expect("File input error")
}
