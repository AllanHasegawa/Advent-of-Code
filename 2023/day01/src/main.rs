use std::fs;

fn main() {
    let test = false;
    let filepath = String::from(if test { "input_test" } else { "input" });
    let input = raw_input(filepath);

    part1(input.clone());
    part2(input.clone());
}

fn part1(input: String) {
    let sum: i32 = input.lines()
        .map(|line| { extract_calibration_value(String::from(line)) })
        .sum();
    println!("Part1: {sum}")
}

fn part2(input: String) {
    let sum: i32 = input.lines()
        .map(|line| { extract_calibration_value(String::from(line)) })
        .sum();
    println!("Part2: {sum}")
}

fn extract_calibration_value(line: String) -> i32 {
    let values: [(&str, i32); 20] = [
        ("0", 0),
        ("1", 1),
        ("2", 2),
        ("3", 3),
        ("4", 4),
        ("5", 5),
        ("6", 6),
        ("7", 7),
        ("8", 8),
        ("9", 9),
        ("zero", 0),
        ("one", 1),
        ("two", 2),
        ("three", 3),
        ("four", 4),
        ("five", 5),
        ("six", 6),
        ("seven", 7),
        ("eight", 8),
        ("nine", 9),
    ];
    let digits: Vec<i32> = line.chars()
        .enumerate()
        .map(|(i, _)| -> Option<i32> {
            let slice = line.chars().skip(i).collect::<String>();

            values.iter()
                .filter(|(str, _)| -> bool { slice.starts_with(str) })
                .take(1)
                .map(|(_, v)| -> i32 { v.clone() })
                .collect::<Vec<_>>()
                .get(0)
                .copied()
        })
        .filter(|o| -> bool { o.is_some() })
        .map(|o| -> i32 { o.expect("oh no!") })
        .collect();
    let first = digits[0];
    let last = digits[digits.len() - 1];
    return (first * 10) + last;
}

fn raw_input(filename: String) -> String {
    fs::read_to_string(filename)
        .expect("File input error")
}