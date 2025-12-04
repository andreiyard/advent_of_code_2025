package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

/* Pseudocode:
1. Read file line by line
2. For each line search for the biggest first digit
	a. Start with "9"
	b. Search the first appearence of the "9" in line (exclude last digit for second part)
	c. if found -> repeat for second digit, bit without limit; if not -> repeat the same with "8", etc
	d. when both parts found convert to int and add to sum
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

func setupLogging(debug bool) {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	opts := &slog.HandlerOptions{
		Level: level,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

const digits = "987654321"

// Change to 12 for part 2
const nDigits = 12

func part1(data string) int {
	sum := 0

	for line := range strings.SplitSeq(data, "\n") {
		if len(line) == 0 {
			continue
		}
		slog.Debug("processing line", "line", line)
		var arr [nDigits]rune
		startI := 0
	nextDigit:
		for digitI := 0; digitI < len(arr); digitI++ {
			slog.Debug("searching for Nth result digit", "n", digitI)
			for _, targetDigit := range digits {
				slog.Debug("trying highest possible digit", "digit", string(targetDigit))
				for i := startI; i < len(line)-(nDigits-1-digitI); i++ {
					// If found
					if line[i] == byte(targetDigit) {
						slog.Debug("!matched!", "index", i, "nextStartI", i+1)
						arr[digitI] = targetDigit
						startI = i + 1
						continue nextDigit
					}
				}
			}
		}
		res := string(arr[:])
		resInt, err := strconv.Atoi(res)
		check(err)
		slog.Info("line solve", "line", line, "res", resInt)
		sum += resInt
	}

	return sum
}

func main() {
	setupLogging(true)
	filename := getFilenameFromArgs()
	dataBytes, err := os.ReadFile(filename)
	check(err)
	data := string(dataBytes)

	result1 := part1(data)
	//result2 := part2(data)

	slog.Warn("part 1 result", "sum", result1)
	//slog.Warn("part 2 result", "sum", sumPart2)
}
