package main

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

/* Pseudocode:
1. Split input to 2 sections (ranges, ids)
2. Parse ranges to struct with ints
3. Iterate over ID (or use goroutines)
4. For each id iterate over all ranges and check if it contains id
5. if contains stop iterating over ranges and increment sum; if doesn't -> continue.
6. Return sum (of fresh product IDs)
*/

type Range struct {
	start, end int
}

func (r Range) contains(num int) bool {
	return num >= r.start && num <= r.end
}

func NewRange(s string) Range {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		panic("invalid range")
	}
	start, err := strconv.Atoi(parts[0])
	if err != nil {
		panic("invalid range start")
	}
	end, err := strconv.Atoi(parts[1])
	if err != nil {
		panic("invalid range end")
	}
	return Range{start, end}
}

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

func parse(data string) ([]Range, []int) {
	// Returns ranges and IDs slices
	data = strings.TrimSpace(data)
	parts := strings.Split(data, "\n\n")

	rangesStrings := strings.Split(parts[0], "\n")
	ranges := make([]Range, len(rangesStrings))
	for i, r := range rangesStrings {
		ranges[i] = NewRange(r)
	}

	idsStrings := strings.Split(parts[1], "\n")
	ids := make([]int, len(idsStrings))
	for i, id := range idsStrings {
		parsedId, err := strconv.Atoi(id)
		if err != nil {
			panic("invalid ID!")
		}
		ids[i] = parsedId
	}
	return ranges, ids
}

func part1(ranges []Range, ids []int) int {
	sum := 0

outer:
	for _, id := range ids {
		for _, r := range ranges {
			if r.contains(id) {
				slog.Debug("Found ID!", "id", id, "range", r)
				sum += 1
				continue outer
			}
		}
		slog.Debug("Wasn't found in any ranges", "id", id)
	}
	return sum
}

func part2(ranges []Range, ids []int) int {
	sum := 0
	return sum
}
func main() {
	setupLogging(true)
	filename := getFilenameFromArgs()
	dataBytes, err := os.ReadFile(filename)
	check(err)
	data := string(dataBytes)

	ranges, ids := parse(data)
	slog.Debug("Parse ranges and IDs", "ranges", ranges, "ids", ids)
	result1 := part1(ranges, ids)
	result2 := part2(ranges, ids)

	slog.Warn("part 1 result", "sum", result1)
	slog.Warn("part 2 result", "sum", result2)
}
