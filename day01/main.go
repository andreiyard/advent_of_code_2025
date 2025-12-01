package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

/* Psuedocode:
0. Start pointer on 50
1. Read input file
2. Iterate over lines (if L - minus value, if R - plus value, after that apply positive MOD, if result is 0 increase counter)
3. Return counter
*/

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide input filename as arg")
		os.Exit(1)
	}

	pointer := 50
	fmt.Printf("starting pos: %d\n", pointer)
	counter_part1 := 0
	counter_part2 := 0

	f, err := os.Open(os.Args[1])
	check(err)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		text := scanner.Text()

		is_passed_zero := false
		move, err := strconv.Atoi(text[1:])
		check(err)

		// Check if move more than one rotation
		full_rotations := move / 100
		counter_part2 += full_rotations
		move %= 100

		if text[0] == 'L' {
			if pointer != 0 && pointer < move {
				is_passed_zero = true
			}
			pointer -= move
		} else {
			if pointer != 0 && pointer+move > 100 {
				is_passed_zero = true
			}
			pointer += move
		}

		fmt.Println(text)
		if is_passed_zero {
			fmt.Println("Passed by 0!")
		}

		// Apply modulo
		pointer += 100
		pointer %= 100

		// Check if pointer is 0
		if pointer == 0 {
			counter_part1 += 1
			counter_part2 += 1
		} else if is_passed_zero { // Check if pointer passed 0
			counter_part2 += 1
		}

		fmt.Printf("current pointer: %d\n", pointer)
	}

	check(scanner.Err())
	fmt.Printf("Part 1 count: %d\n", counter_part1)
	fmt.Printf("Part 2 count: %d\n", counter_part2)
}
