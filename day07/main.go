package main

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/andreiyard/advent_of_code_2025/utils"
)

type Point struct {
	x, y int
}

type Cells map[Point]string

func (c Cells) getActiveBeams() []Point {
	var points []Point
	for k, v := range c {
		if v == "S" {
			points = append(points, k)
		}
	}
	return points
}

func (c Cells) printGrid() {
	var maxX, maxY int
	for k, _ := range c {
		if k.x > maxX {
			maxX = k.x
		}
		if k.y > maxY {
			maxY = k.y
		}
	}

	for i := range maxX + 1 {
		for j := range maxY + 1 {
			val := c[Point{i, j}]
			switch val {
			case "":
				fmt.Print(".")
			default:
				fmt.Print(val)
			}
		}
		fmt.Print("\n")
	}
}

func printPathsGrid(c map[Point]int) {
	var maxX, maxY int
	for k, _ := range c {
		if k.x > maxX {
			maxX = k.x
		}
		if k.y > maxY {
			maxY = k.y
		}
	}

	for i := range maxX + 1 {
		for j := range maxY + 1 {
			val := c[Point{i, j}]
			switch val {
			case 0:
				fmt.Print(".")
			default:
				fmt.Print(val)
			}
		}
		fmt.Print("\n")
	}
}
func (c Cells) extendBeam(start Point, limitX int) (splitted bool) {
	// Start for point and check next x point until end of the grid or split
	for x := start.x; x < limitX; x++ {
		point := Point{x, start.y}
		val := c[point]
		if val == "^" {
			// add new beam points to the grid because split
			left := Point{x, start.y - 1}
			right := Point{x, start.y + 1}
			c[left] = "S"
			c[right] = "S"
			// mark splitted
			splitted = true
			c[point] = "V"
			break
		} else if val == "V" {
			// if already was used to split beam do not count
			// only create splitted beams
			left := Point{x, start.y - 1}
			right := Point{x, start.y + 1}
			c[left] = "S"
			c[right] = "S"
			break
		} else {
			c[point] = "|"
		}
	}
	// remove starting point from grid
	c[start] = "|"
	return
}

func parse(data string) (Cells, int) {
	data = strings.TrimSpace(data)
	lines := strings.Split(data, "\n")
	cells := make(Cells)

	for i, line := range lines {
		for j, char := range line {
			point := Point{i, j}
			if char == '.' {
				continue
			}
			cells[point] = string(char)
		}
	}

	return cells, len(lines)
}

func part1(cells Cells, limitX int) int {
	sum := 0

	// while grid contains beams starts
	for {
		//cells.printGrid()
		nextBeams := cells.getActiveBeams()
		slog.Debug("Next beams to check", "beamsStarts", nextBeams)
		if len(nextBeams) == 0 {
			break
		}
		for _, beamStart := range nextBeams {
			if cells.extendBeam(beamStart, limitX) {
				sum++
				slog.Debug("Beam splitted", "beam", beamStart, "sum", sum)
			} else {
				slog.Debug("Beam NOT splitted", "beam", beamStart, "sum", sum)
			}
		}
	}
	cells.printGrid()

	return sum
}

func (c Cells) getWithX(x int, value string) Cells {
	result := make(Cells)
	for p, val := range c {
		if p.x == x && val == value {
			result[p] = val
		}
	}
	return result
}

func part2(cells Cells, limitX int) int {
	// go from top to bottom and track possible paths in each moment
	paths := make(map[Point]int)

	for i := range limitX {
		for p, _ := range cells.getWithX(i, "V") {
			// Just set the same paths num as x-1 "|" cell
			prevPoint := Point{i - 1, p.y}
			paths[p] = paths[prevPoint]
		}
		for p, _ := range cells.getWithX(i, "|") {
			// if x == 0, then set paths num to 1
			if i == 0 {
				paths[p] = 1
				continue
			}
			// If x - 1 cell also contains "|" -> set the same paths num
			// If next to "V" add path nums to sum of "V" cells on the right and left
			summarized := 0
			prevPoint := Point{i - 1, p.y}
			if cells[prevPoint] == "|" {
				summarized = paths[prevPoint]
			}
			// Check left and right cells (if V -> add to current paths num)
			checkAndGetPaths := func(point Point) int {
				if cells[point] == "V" {
					return paths[point]
				} else {
					return 0
				}

			}
			leftPoint := Point{i, p.y - 1}
			summarized += checkAndGetPaths(leftPoint)
			rightPoint := Point{i, p.y + 1}
			summarized += checkAndGetPaths(rightPoint)

			paths[p] = summarized

		}
	}
	if os.Getenv("DEBUG") != "" {
		printPathsGrid(paths)
	}
	slog.Debug("Calculated paths grid", "paths", paths)

	sum := 0
	for p, pathsNum := range paths {
		// Take only last line
		if p.x == limitX-1 {
			sum += pathsNum
		}

	}
	return sum
}

func main() {
	utils.SetupLoggingEnv()
	filename := utils.GetFilenameFromArgs()
	dataBytes, err := os.ReadFile(filename)
	utils.Check(err)
	data := string(dataBytes)
	slog.Debug("Got data", "data", data)

	cells, limitX := parse(data)
	slog.Debug("Parsed input", "cells", cells, "limitX", limitX)
	result1 := part1(cells, limitX)

	// Take resulting grid from part1
	result2 := part2(cells, limitX)

	slog.Warn("part 1 result", "sum", result1)
	slog.Warn("part 2 result", "sum", result2)
}
