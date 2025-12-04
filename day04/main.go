package main

import (
	"fmt"
	"iter"
	"log/slog"
	"os"
	"strings"
)

/* Pseudocode:
I plan to use naive approach
1. Go over all cells (skip empty)
2. Check the number of adjacent filled cells
3. if < 4 -> cell matched
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

type Point struct {
	x, y int
}

type Grid struct {
	x, y  int     // Dimensions
	cells [][]int // Cells data
}

func (g Grid) All() iter.Seq[Point] {
	return func(yield func(Point) bool) {
		for i := 0; i < g.x; i++ {
			for j := 0; j < g.y; j++ {
				if !yield(Point{i, j}) {
					return
				}
			}
		}
	}
}

func (g Grid) get(p Point) int {
	return g.cells[p.x][p.y]
}

func (g Grid) isInBounds(p Point) bool {
	x, y := p.x, p.y
	switch {
	case x >= g.x:
		return false
	case x < 0:
		return false
	case y >= g.y:
		return false
	case y < 0:
		return false
	}
	return true
}

func (g Grid) isCellAccessed(p Point) bool {
	n := 0
	for _, p := range g.getAdjacentPoints(p) {
		if g.get(p) == 1 {
			n += 1
		}
	}
	return n < 4
}

func (g Grid) getAdjacentPoints(p Point) []Point {
	points := make([]Point, 0, 9)
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			newPoint := Point{p.x + i, p.y + j}
			// Skip central point
			if newPoint == p {
				continue
			}
			// Skip points out of the grid
			if !g.isInBounds(newPoint) {
				continue
			}
			points = append(points, newPoint)
		}
	}
	return points
}

func NewGrid(data string) Grid {
	// detect dimesions, create slice of slices, return Grid
	data = strings.TrimSpace(data)
	lines := strings.Split(data, "\n")
	x := len(lines)
	y := len(lines[0])

	cells := make([][]int, x)
	for i, line := range lines {
		row := make([]int, y)
		for j := 0; j < len(line); j++ {
			if line[j] == '@' {
				row[j] = 1
			}
		}
		cells[i] = row
	}
	return Grid{x, y, cells}
}

func part1(data string) int {
	sum := 0
	grid := NewGrid(data)
	slog.Debug("New grid from data", "grid", grid)
	for p := range grid.All() {
		// Skip empty cells
		if grid.get(p) != 1 {
			continue
		}
		if grid.isCellAccessed(p) {
			slog.Debug("!Accessed!", "point", p)
			sum += 1
		}
	}
	return sum
}

func (g Grid) removePoints(points []Point) int {
	for _, p := range points {
		g.cells[p.x][p.y] = 0
	}
	return len(points)
}

func part2(data string) int {
	sum := 0
	grid := NewGrid(data)
	slog.Debug("New grid from data", "grid", grid)

	for i := 1; ; i++ {
		pointsToRemove := make([]Point, 0)
		for p := range grid.All() {
			// Skip empty cells
			if grid.get(p) != 1 {
				continue
			}
			if grid.isCellAccessed(p) {
				slog.Debug("!Accessed!", "point", p)
				pointsToRemove = append(pointsToRemove, p)
			}
		}
		removed := grid.removePoints(pointsToRemove)
		if removed > 0 {
			slog.Debug("Removed points", "iteration", i, "nPoints", removed)
			sum += removed
		} else {
			slog.Debug("No removable points found. Exiting")
			break
		}
	}
	return sum
}
func main() {
	setupLogging(false)
	filename := getFilenameFromArgs()
	dataBytes, err := os.ReadFile(filename)
	check(err)
	data := string(dataBytes)

	result1 := part1(data)
	result2 := part2(data)

	slog.Warn("part 1 result", "sum", result1)
	slog.Warn("part 2 result", "sum", result2)
}
