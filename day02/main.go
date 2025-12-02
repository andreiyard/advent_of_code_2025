package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

/* Pseudocode:
Main logic:
1. Read file contents
2. Iterate over ranges
	a. expand range
	b. iterate over each number and check if number is invalid
	c. return list of invalid numbers
3. Sum all invalid numbers

Invalid number logic:
n-digit sequence must appear multiple times
for example:
11
1212
443443

- If we have 1-digit repeating, we don't need to check longer squences
- Sequence (n value) might be as long as l/2
- n must be factor of l
*/

func check(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func getFilenameFromArgs() string {
	if len(os.Args) < 2 {
		fmt.Println("Please provide input filename as arg")
		os.Exit(1)
	}
	return os.Args[1]
}

func isInvalidPart1(number int) bool {
	slog.Debug("Checking number", "number", number)
	numberString := strconv.Itoa(number)
	numberLen := len(numberString)
	if numberLen%2 == 1 {
		return false
	}
	n := numberLen / 2

	seq := numberString[:n]
	slog.Debug(seq)

	// Check that first part is equal to second
	if numberString[n:] == seq {
		slog.Debug("Matched!")
		return true
	} else {
		return false
	}
}

func isInvalid(number int) bool {
	slog.Debug("Checking number", "number", number)
	numberString := strconv.Itoa(number)
	numberLen := len(numberString)
	maxN := numberLen / 2

	// Check all sequence lengths
outer:
	for n := 1; n <= maxN; n++ {
		// if n not factor of numberLen, then skip
		if numberLen%n != 0 {
			continue
		}

		seq := numberString[:n]
		slog.Debug(seq)

		// Check that sequence is repeating
		for i := n; i < numberLen; i += n {
			if numberString[i:i+n] == seq {
				continue
			} else {
				slog.Debug("didnt match", "seq", seq)
				continue outer
			}
		}

		slog.Debug("Matched all!")
		return true
	}

	return false
}

func main() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	filename := getFilenameFromArgs()
	bytes, err := os.ReadFile(filename)
	check(err)
	text := string(bytes)
	text = strings.ReplaceAll(text, "\n", "")

	sumPart1 := 0
	sumPart2 := 0

	for r := range strings.SplitSeq(text, ",") {
		boundaries := strings.Split(r, "-")

		start, _ := strconv.Atoi(boundaries[0])
		end, _ := strconv.Atoi(boundaries[1])
		slog.Info("Checking range", "range", r, "start", start, "end", end, "numbers_to_check", end-start+1)

		// Iterate over all numbers in range (Part 1)
		for i := start; i <= end; i++ {
			if isInvalidPart1(i) {
				slog.Info("Invalid number found", "number", i)
				sumPart1 += i
			}
		}

		// Iterate over all numbers in range (Part 2)
		for i := start; i <= end; i++ {
			if isInvalid(i) {
				slog.Info("Invalid number found", "number", i)
				sumPart2 += i
			}
		}
	}

	slog.Warn("part 1 result", "sum", sumPart1)
	slog.Warn("part 2 result", "sum", sumPart2)

	slog.Debug("result", "isInvalid", isInvalid(1010))
	slog.Debug("result", "isInvalid", isInvalid(1235324))
	slog.Debug("result", "isInvalid", isInvalid(77777))
	slog.Debug("result", "isInvalid", isInvalid(43214321))
	slog.Debug("result", "isInvalid", isInvalid(654321654321431232))
}
